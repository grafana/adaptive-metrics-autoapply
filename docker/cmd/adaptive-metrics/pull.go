package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/grafana/adaptive-metrics-autoapply/docker/internal"
)

func pull(args []string) {
	defaultWorkingDir := "./"
	if workingDirEnvVar := os.Getenv("INPUT_WORKING-DIR"); workingDirEnvVar != "" {
		defaultWorkingDir = workingDirEnvVar
	}

	flags := flag.NewFlagSet("pull", flag.ExitOnError)
	workingDir := flags.String("working-dir", defaultWorkingDir, "The path to the working directory.")
	userAgent := flags.String("user-agent", "gh-action-autoapply", "The user-agent to use when making requests against the API.")
	writeSegments := flags.Bool("write-segments", false, "Optionally write a segments.json file to disk.")

	err := flags.Parse(args)
	if err != nil {
		log.Fatalf("failed to parse flags: %v", err)
	}

	apiURL := mustGetEnv("GRAFANA_AM_API_URL")
	apiKey := mustGetEnv("GRAFANA_AM_API_KEY")

	c := internal.NewClient(&http.Client{}, *userAgent, apiURL, apiKey)

	// Fetch all segments.
	segments, err := c.FetchSegments()
	if err != nil {
		log.Fatalf("failed to fetch segments: %v", err)
	}

	if *writeSegments {
		log.Printf("writing segments.json with %d segments", len(segments))
		err = writeJSONToFile(filepath.Join(*workingDir, "segments.json"), segments)
		if err != nil {
			log.Fatalf("failed to write segments.json: %v", err)
		}
	}

	// Add the default segment.
	segments = append(segments, internal.DefaultSegment)

	gha, err := newGithubActionWorkflowCommands()
	if err != nil {
		log.Fatalf("failed to create github action workflow commands: %v", err)
	}

	totalSeriesChange := 0
	totalSeries := 0
	output := new(strings.Builder)
	for _, segment := range segments {
		// Fetch recommendations for each segment.
		recs, err := c.FetchRecommendations(segment, true)
		if err != nil {
			log.Fatalf("failed to fetch recommendations for segment %s: %v", segment.Name, err)
		}

		// Sort exact match rules first, then sort by metric name.
		slices.SortStableFunc(recs, func(a, b internal.Recommendation) int {
			// If both are exact matches, sort by metric name.
			if isExactMatch(a) && isExactMatch(b) {
				return strings.Compare(a.Metric, b.Metric)
			}
			// Otherwise sort exact matches first
			if a.MatchType != b.MatchType {
				if isExactMatch(a) {
					return -1
				}
				return 1
			}
			// Otherwise don't change anything, since it may change the semantics of the ruleset.
			return 0
		})

		// Strip the managed_by field from the recommendations. This adds unnecessary noise to the files, and is overwritten when applying the rules anyway.
		for i, r := range recs {
			r.ManagedBy = ""
			recs[i] = r
		}

		// Write the recommendations to a file.
		var filename string
		if segment == internal.DefaultSegment {
			filename = "recommendations.json"
		} else {
			filename = fmt.Sprintf("recommendations-%s.json", segment.Name)
		}
		log.Printf("writing recommendations for segment %s to %s with %d rules", segment.Name, filename, len(recs))
		err = writeJSONToFile(filepath.Join(*workingDir, filename), internal.ConvertVerboseToRules(recs))
		if err != nil {
			log.Fatalf("failed to write recommendations for segment %s: %v", segment.Name, err)
		}

		writeChanges(output, segment, recs)

		segmentChange := seriesChangeForSegment(recs)
		err = gha.writeOutput(fmt.Sprintf("series-change-%s", segment.Name), strconv.Itoa(segmentChange))
		if err != nil {
			log.Fatalf("failed to write series-change output for segment %s: %v", segment.Name, err)
		}

		segmentTotal := totalSeriesForSegment(recs)
		err = gha.writeOutput(fmt.Sprintf("series-total-%s", segment.Name), strconv.Itoa(segmentTotal))
		if err != nil {
			log.Fatalf("failed to write series-total output for segment %s: %v", segment.Name, err)
		}

		totalSeriesChange += segmentChange
		totalSeries += segmentTotal
	}

	err = gha.writeOutput("series-change", strconv.Itoa(totalSeriesChange))
	if err != nil {
		log.Fatalf("failed to write series-change output: %v", err)
	}

	err = gha.writeOutput("series-total", strconv.Itoa(totalSeries))
	if err != nil {
		log.Fatalf("failed to write series-total output: %v", err)
	}

	err = gha.writeStepSummary(output.String())
	if err != nil {
		log.Fatalf("failed to write step summary: %v", err)
	}
}

func totalSeriesForSegment(recs []internal.Recommendation) int {
	var total int
	for _, rec := range recs {
		total += rec.CurrentSeriesCount
	}
	return total
}

func seriesChangeForSegment(recs []internal.Recommendation) int {
	var total int
	for _, rec := range recs {
		total += rec.RecommendedSeriesCount - rec.CurrentSeriesCount
	}
	return total
}

func writeChanges(output io.Writer, segment internal.Segment, recs []internal.Recommendation) {
	type change struct {
		seriesChange int
		action       string
		metric       string
	}

	var changes []change
	for _, rec := range recs {
		if rec.RecommendedAction == "keep" {
			continue
		}
		changes = append(changes, change{
			seriesChange: rec.RecommendedSeriesCount - rec.CurrentSeriesCount,
			action:       rec.RecommendedAction,
			metric:       rec.Metric,
		})
	}

	if len(changes) == 0 {
		return
	}

	sort.Slice(changes, func(i, j int) bool {
		return changes[i].seriesChange > changes[j].seriesChange
	})

	fmt.Fprintf(output, "## Segment %q\n", segment.Name)

	fmt.Fprintf(output, "### Series Change\n")
	fmt.Fprintf(output, "Total series change: %d\n", seriesChangeForSegment(recs))
	fmt.Fprintf(output, "Total series: %d\n", totalSeriesForSegment(recs))
	fmt.Fprintf(output, "Percentage change: %.2f%%\n", float64(seriesChangeForSegment(recs))/float64(totalSeriesForSegment(recs))*100)

	fmt.Fprintln(output, "| Metric | Action | Series Change |")
	fmt.Fprintln(output, "|--------|--------|---------------|")
	for _, c := range changes {
		fmt.Fprintf(output, "| %s | %s | %d |\n", c.metric, c.action, c.seriesChange)
	}
}

func isExactMatch(rule internal.Recommendation) bool {
	return rule.MatchType == "exact" || rule.MatchType == ""
}

func mustGetEnv(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("missing required env var %s", key)
	}

	return val
}

func writeJSONToFile(filePath string, obj any) error {
	out, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return err
	}

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	_, err = f.Write(out)
	return err
}

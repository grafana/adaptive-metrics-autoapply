package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"slices"
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

	for _, segment := range segments {
		// Fetch recommendations for each segment.
		recs, err := c.FetchRecommendations(segment, false)
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

		// Write the recommendations to a file.
		var filename string
		if segment == internal.DefaultSegment {
			filename = "recommendations.json"
		} else {
			filename = fmt.Sprintf("recommendations-%s.json", segment.Name)
		}
		log.Printf("writing recommendations for segment %s to %s with %d rules", segment.Name, filename, len(recs))
		err = writeJSONToFile(filepath.Join(*workingDir, filename), recs)
		if err != nil {
			log.Fatalf("failed to write recommendations for segment %s: %v", segment.Name, err)
		}
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

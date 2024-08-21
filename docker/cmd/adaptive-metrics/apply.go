package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/grafana/adaptive-metrics-autoapply/docker/internal"
)

func apply(args []string) {
	defaultDryRun := false
	if dryRunEnvVar := os.Getenv("INPUT_DRY-RUN"); dryRunEnvVar != "" {
		var err error
		defaultDryRun, err = strconv.ParseBool(dryRunEnvVar)
		if err != nil {
			log.Fatalf("error parsing INPUT_DRY-RUN: %s", err)
		}
	}

	defaultWorkingDir := "./"
	if workingDirEnvVar := os.Getenv("INPUT_WORKING-DIR"); workingDirEnvVar != "" {
		defaultWorkingDir = workingDirEnvVar
	}

	flags := flag.NewFlagSet("apply", flag.ExitOnError)
	workingDir := flags.String("working-dir", defaultWorkingDir, "The path to the working directory.")
	dryRun := flags.Bool("dry-run", defaultDryRun, "dry run; print changes but do not apply them")
	userAgent := flags.String("user-agent", "gh-action-autoapply", "The user-agent to use when making requests against the API.")

	err := flags.Parse(args)
	if err != nil {
		log.Fatalf("failed to parse flags: %v", err)
	}

	err = os.Chdir(*workingDir)
	if err != nil {
		log.Fatalf("failed to change working directory: %v", err)
	}

	apiURL := mustGetEnv("GRAFANA_AM_API_URL")
	apiKey := mustGetEnv("GRAFANA_AM_API_KEY")

	c := internal.NewClient(&http.Client{}, *userAgent, apiURL, apiKey)

	segments, err := c.FetchSegments()
	if err != nil {
		log.Fatalf("failed to read segments: %v", err)
	}

	output := os.Stdout
	if summaryOutputPath := os.Getenv("GITHUB_STEP_SUMMARY"); summaryOutputPath != "" {
		output, err = os.Open(summaryOutputPath)
		if err != nil {
			log.Fatalf("failed to open summary output file: %v", err)
		}
		defer output.Close()
	}

	totalChanges := 0
	unchangedSegments := 0
	segments = append(segments, internal.DefaultSegment)
	for _, segment := range segments {
		changes, err := applySegment(output, c, segment, *dryRun)
		if err != nil {
			log.Fatalf("failed to apply segment %s: %v", segment.Name, err)
		}

		totalChanges += changes
		if changes == 0 {
			unchangedSegments++
		}
	}

	if totalChanges == 0 {
		fmt.Fprintln(output, "No changes detected in aggregation rules.")
	} else {
		fmt.Fprintln(output, "#### Summary")
		fmt.Fprintf(output, "- %d changes detected in aggregation rules\n", totalChanges)
		fmt.Fprintf(output, "- %d modified segments\n", len(segments)-unchangedSegments)
		fmt.Fprintf(output, "- %d unmodified segments\n", unchangedSegments)
	}
}

func applySegment(output io.Writer, client *internal.Client, segment internal.Segment, dryRun bool) (int, error) {
	filename := fmt.Sprintf("recommendations-%s.json", segment.Name)
	if segment == internal.DefaultSegment {
		filename = "recommendations.json"
	}

	rules, err := readJSONFile[[]internal.Recommendation](filename)
	if err != nil {
		if !os.IsNotExist(err) {
			return 0, fmt.Errorf("failed to read rules: %w", err)
		}
		log.Printf("no rules found for segment %q", segment.Name)
		rules = []internal.Recommendation{}
	}

	err = client.ValidateRules(rules)
	if err != nil {
		return 0, fmt.Errorf("failed to validate rules: %w", err)
	}

	currentState, etag, err := client.GetRules(segment)
	if err != nil {
		return 0, fmt.Errorf("failed to get current rules: %w", err)
	}

	changes := writeDiff(output, segment, currentState, rules)

	if !dryRun {
		return changes, client.UpdateRules(segment, etag, rules)
	}
	return changes, nil
}

func readJSONFile[T any](path string) (T, error) {
	var result T
	file, err := os.Open(path)
	if err != nil {
		return result, err
	}

	err = json.NewDecoder(file).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

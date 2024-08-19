package main

import (
	"encoding/json"
	"flag"
	"fmt"
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

	segments, err := readJSONFile[[]internal.Segment]("segments.json")
	if err != nil {
		log.Fatalf("failed to read segments: %v", err)
	}

	segments = append(segments, internal.DefaultSegment)
	for _, segment := range segments {
		err := applySegment(c, segment, *dryRun)
		if err != nil {
			log.Fatalf("failed to apply segment %s: %v", segment.Name, err)
		}
	}
}

func applySegment(client *internal.Client, segment internal.Segment, dryRun bool) error {
	filename := fmt.Sprintf("recommendations-%s.json", segment.Name)
	if segment == internal.DefaultSegment {
		filename = "recommendations.json"
	}

	rules, err := readJSONFile[[]internal.Recommendation](filename)
	if err != nil {
		return fmt.Errorf("failed to read rules: %w", err)
	}

	log.Printf("applying segment %q num-rules=%d dry-run=%t", segment.Name, len(rules), dryRun)

	if dryRun {
		return nil
	}

	return client.UpdateRules(segment, rules)
}

func readJSONFile[T any](path string) (T, error) {
	var result T
	file, err := os.Open(path)
	if err != nil {
		return result, err
	}

	err = json.NewDecoder(file).Decode(&result)
	if err != nil {
		return result, nil
	}

	return result, nil
}

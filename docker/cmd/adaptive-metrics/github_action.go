package main

import (
	"fmt"
	"io"
	"os"
)

type githubActionWorkflowCommands struct {
	outputFile  io.WriteCloser
	summaryFile io.WriteCloser
}

func newGithubActionWorkflowCommands() (*githubActionWorkflowCommands, error) {
	if os.Getenv("GITHUB_ACTIONS") != "true" {
		return nil, nil
	}

	outputFile, err := os.OpenFile(os.Getenv("GITHUB_OUTPUT"), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	summaryFile, err := os.OpenFile(os.Getenv("GITHUB_STEP_SUMMARY"), os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	return &githubActionWorkflowCommands{
		outputFile:  outputFile,
		summaryFile: summaryFile,
	}, nil
}

func (c *githubActionWorkflowCommands) writeStepSummary(summary string) error {
	if c == nil {
		return nil
	}

	_, err := io.WriteString(c.summaryFile, summary)
	return err
}

func (c *githubActionWorkflowCommands) writeOutput(name, value string) error {
	if c == nil {
		return nil
	}

	_, err := fmt.Fprintf(c.outputFile, "%s=%s\n", name, value)
	return err
}

func (c *githubActionWorkflowCommands) close() {
	if c == nil {
		return
	}

	c.outputFile.Close()
	c.summaryFile.Close()
}

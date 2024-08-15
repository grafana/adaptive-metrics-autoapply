package main

import (
	"context"
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/terraform-exec/tfexec"
)

func main() {

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

	workingDir := flag.String("working-dir", defaultWorkingDir, "The path to the working directory.")
	dryRun := flag.Bool("dry-run", defaultDryRun, "dry run; run terraform plan instead of terraform apply")
	flag.Parse()

	installer := &releases.ExactVersion{
		Product: product.Terraform,
		Version: version.Must(version.NewVersion("1.9")),
	}

	execPath, err := installer.Install(context.Background())
	if err != nil {
		log.Fatalf("error installing Terraform: %s", err)
	}

	tf, err := tfexec.NewTerraform(*workingDir, execPath)
	if err != nil {
		log.Fatalf("error running NewTerraform: %s", err)
	}

	tf.SetStderr(os.Stderr)
	tf.SetStdout(os.Stdout)

	err = tf.Init(context.Background(), tfexec.Upgrade(true))
	if err != nil {
		log.Fatalf("error running Init: %s", err)
	}

	if *dryRun {
		_, err = tf.PlanJSON(context.Background(), os.Stdout)
		if err != nil {
			log.Fatalf("error running Plan: %s", err)
		}
	} else {
		err = tf.Apply(context.Background())
		if err != nil {
			log.Fatalf("error running Apply: %s", err)
		}
	}
}

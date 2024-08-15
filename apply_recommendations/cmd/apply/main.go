package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/terraform-exec/tfexec"
)

func main() {
	workingDir := flag.String("working-dir", os.Getenv("INPUT_WORKING-DIR"), "The path to the working directory.")
	dryRun := flag.Bool("dry-run", false, "dry run; run terraform plan instead of terraform apply")
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

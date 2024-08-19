package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("missing command, available commands: pull, apply")
	}

	switch os.Args[1] {
	case "pull":
		pull(os.Args[2:])
	case "apply":
		apply(os.Args[2:])
	default:
		log.Fatalf("unknown command %s, available commands: pull, apply", os.Args[1])
	}
}

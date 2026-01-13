package main

import (
	"fmt"
	"os"
	"simpledeploy/internal/cli"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: simpledeploy <command>")
		fmt.Println("Commands: deploy, status")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "deploy":
		if err := cli.Deploy(); err != nil {
			fmt.Println("Failed to deploy ", err)
			os.Exit(1)
		}

	case "status":
		if err := cli.Status(); err != nil {
			fmt.Println("Status failed:", err)
			os.Exit(1)
		}
	default:
		fmt.Println("Unknown command:", command)
		os.Exit(1)
	}

}

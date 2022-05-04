package main

import "fmt"

var (
	Version   = "0.0.1"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

func main() {
	fmt.Printf("Ohana v%s (built %s, commit %s)\n", Version, BuildTime, GitCommit)
}

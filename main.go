package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/Jason-ZW/autok3s-geo/pkg/cmd"
)

var (
	gitVersion   string
	gitCommit    string
	gitTreeState string
	buildDate    string
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	rootCmd := cmd.Command()
	rootCmd.AddCommand(cmd.ServeCommand(), cmd.VersionCommand(gitVersion, gitCommit, gitTreeState, buildDate))

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

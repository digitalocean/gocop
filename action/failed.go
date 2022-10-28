package action

import (
	"fmt"
	"log"
	"strings"

	"github.com/digitalocean/gocop/gocop"
	"github.com/spf13/cobra"
)

type failedCmdFlags struct {
	src       string
	test2json bool
}

var failedFlags failedCmdFlags

var failedCmd = &cobra.Command{
	Use:   "failed",
	Short: "lists failed test events from test run",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var parser gocop.Parser
		if failedFlags.test2json {
			parser = &gocop.Test2JSONParser{}
		} else {
			parser = &gocop.StandardParser{}
		}

		pkgs, err := gocop.ParseFileFailedPackages(parser, failedFlags.src)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Print(strings.Join(pkgs, "\n"))
	},
}

func init() {
	RootCmd.AddCommand(failedCmd)

	failedCmd.Flags().StringVarP(&failedFlags.src, "src", "s", "", "source test output file")
	failedCmd.Flags().BoolVarP(
		&failedFlags.test2json,
		"test2json", "", false,
		"set to true if the test output format is test2json format",
	)
	err := failedCmd.MarkFlagRequired("src")
	if err != nil {
		log.Fatal(err)
	}
}

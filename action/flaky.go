package action

import (
	"fmt"
	"log"
	"strings"

	"github.com/digitalocean/gocop/gocop"
	"github.com/spf13/cobra"
)

type flakyCmdFlags struct {
	retests   []string
	test2json bool
}

var flakyFlags flakyCmdFlags

var flakyCmd = &cobra.Command{
	Use:   "flaky",
	Short: "lists packages suspected of having flaky tests",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var parser gocop.Parser
		if flakyFlags.test2json {
			parser = &gocop.Test2JSONParser{}
		} else {
			parser = &gocop.StandardParser{}
		}

		pkgs, err := gocop.FlakyFilePackages(parser, flakyFlags.retests...)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Print(strings.Join(pkgs, "\n"))
	},
}

func init() {
	RootCmd.AddCommand(flakyCmd)

	flakyCmd.Flags().StringSliceVarP(&flakyFlags.retests, "retests", "r", []string{}, "source output for retests")
	flakyCmd.Flags().BoolVarP(
		&flakyFlags.test2json,
		"test2json", "", false,
		"set to true if the test output format is test2json format",
	)
	err := flakyCmd.MarkFlagRequired("retests")
	if err != nil {
		log.Fatal(err)
	}
}

package action

import (
	"fmt"
	"log"
	"strings"

	"github.com/digitalocean/gocop/gocop"
	"github.com/spf13/cobra"
)

type flakyCmdFlags struct {
	retests []string
}

var flakyFlags flakyCmdFlags

var flakyCmd = &cobra.Command{
	Use:   "flaky",
	Short: "lists packages suspected of having flaky tests",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		pkgs := gocop.FlakyFilePackages(flakyFlags.retests...)
		fmt.Print(strings.Join(pkgs, "\n"))
	},
}

func init() {
	RootCmd.AddCommand(flakyCmd)

	flakyCmd.Flags().StringSliceVarP(&flakyFlags.retests, "retests", "r", []string{}, "source output for retests")
	err := flakyCmd.MarkFlagRequired("retests")
	if err != nil {
		log.Fatal(err)
	}
}

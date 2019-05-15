package action

import (
	"fmt"
	"strings"

	"github.com/digitalocean/gocop/gocop"
	"github.com/spf13/cobra"
)

var retests []string

var flakyCmd = &cobra.Command{
	Use:   "flaky",
	Short: "lists packages suspected of having flaky tests",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		pkgs := gocop.FlakyFile(retests...)
		fmt.Print(strings.Join(pkgs, "\n"))
	},
}

func init() {
	RootCmd.AddCommand(flakyCmd)

	flakyCmd.Flags().StringSliceVarP(&retests, "retests", "r", []string{}, "source output for retests")
	flakyCmd.MarkFlagRequired("retests")
}

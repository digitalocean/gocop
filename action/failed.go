package action

import (
	"fmt"
	"strings"

	"github.com/digitalocean/gocop/gocop"
	"github.com/spf13/cobra"
)

var src string

var failedCmd = &cobra.Command{
	Use:   "failed",
	Short: "lists failed packages from test run",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		pkgs := gocop.ParseFile(src)
		fmt.Print(strings.Join(pkgs, "\n"))
	},
}

func init() {
	RootCmd.AddCommand(failedCmd)

	failedCmd.Flags().StringVarP(&src, "src", "s", "", "source test output file")
	failedCmd.MarkFlagRequired("src")
}

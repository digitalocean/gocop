package action

import (
	"fmt"
	"log"
	"strings"

	"github.com/digitalocean/gocop/gocop"
	"github.com/spf13/cobra"
)

type failedCmdFlags struct {
	src string
}

var failedFlags failedCmdFlags

var failedCmd = &cobra.Command{
	Use:   "failed",
	Short: "lists failed packages from test run",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		pkgs, err := gocop.ParseFileFailedPackages(failedFlags.src)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(strings.Join(pkgs, "\n"))
	},
}

func init() {
	RootCmd.AddCommand(failedCmd)

	failedCmd.Flags().StringVarP(&failedFlags.src, "src", "s", "", "source test output file")
	err := failedCmd.MarkFlagRequired("src")
	if err != nil {
		log.Fatal(err)
	}
}

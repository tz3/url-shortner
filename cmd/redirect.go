package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tz3/url-shortner/cmd/redirect"
)

var redirectCmd = &cobra.Command{
	Use:   "Redirect Urls",
	Short: "Subcommand related to prediction of the urls",
	RunE:  Noop,
}

func init() {
	redirectCmd.AddCommand(redirect.Serve)
}

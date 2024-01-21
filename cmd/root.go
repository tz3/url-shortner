// Package cmd root.go
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var Root = &cobra.Command{
	Use:   "url-shortner",
	Short: "Link",
	Long: `Take any lengthy URL and make it into something shorter.
A secondary purpose of this application is me (Moutaz Chaara) to sharpen my skills in GOLANG.
If prospective employers come looking, here's some
code!`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("This application still in WIP.")
	},
}

// init to declare configuration
func init() {
	Root.AddCommand(redirectCmd)
}

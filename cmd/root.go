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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main() func. It only needs to happen once to the root.
func Execute(c *cobra.Command) int {
	err := c.Execute()
	if err != nil {
		fmt.Println("Error:", err)
		return ExitGeneralError
	}
	return ExitSuccess
}

// init to declare configuration
func init() {
}

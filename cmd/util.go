package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

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

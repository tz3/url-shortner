// util_test.go
package cmd

import (
	"fmt"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	for _, tc := range []struct {
		name string
		cmd  *cobra.Command
		exit int
	}{
		{
			name: "Everything OK",
			cmd: &cobra.Command{
				RunE: func(cmd *cobra.Command, args []string) error {
					return nil
				},
			},
			exit: ExitSuccess,
		},
		{
			name: "Passed an error",
			cmd: &cobra.Command{
				RunE: func(cmd *cobra.Command, args []string) error {
					return fmt.Errorf("its broke")
				},
			},
			exit: ExitGeneralError,
		},
		{
			name: "Passed an error (without specific exit code)",
			cmd: &cobra.Command{
				RunE: func(cmd *cobra.Command, args []string) error {
					return fmt.Errorf("I am very mysterious")
				},
			},
			exit: ExitGeneralError,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tc := tc
			t.Parallel()

			exitCode := Execute(tc.cmd)

			assert.Equal(t, tc.exit, exitCode)
		})
	}
}

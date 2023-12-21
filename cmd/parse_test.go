package cmd

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseCmd(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		expectedError  string
	}{
		{
			name:           "Valid parse",
			args:           []string{"parse", "--file-path", "tests/empty.zed"},
			expectedOutput: "",
		},
		{
			name:           "Invalid parse",
			args:           []string{"parse", "--file-path", "tests/broken.zed"},
			expectedOutput: "Complied error parse error in `schema`, line 1, column 1: Expected end of statement or definition, found: TokenTypeError", // Add the expected error message
		},
		// Add more test cases as needed
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			var stdout bytes.Buffer
			var rootCmd = &cobra.Command{}
			rootCmd.AddCommand(parseCmd)
			rootCmd.SetOut(&stdout)
			rootCmd.SetArgs(test.args)
			err := rootCmd.Execute()
			if err != nil {
				fmt.Printf("error %s", err)
			}
			assert.Equal(t, test.expectedOutput, stdout.String())
		})
	}
}

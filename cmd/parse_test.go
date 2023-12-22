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
			expectedOutput: "Complied error: parse error in `schema`, line 1, column 1: Expected end of statement or definition, found: TokenTypeError\n", // Add the expected error message
		},
		{
			name:           "associativity parse",
			args:           []string{"parse", "--file-path", "tests/associativity.zed"},
			expectedOutput: "Complied error: parse error in `schema`, line 2, column 5: error in permission union: invalid Relation.UsersetRewrite: embedded message failed validation | caused by: invalid UsersetRewrite.Union: embedded message failed validation | caused by: invalid SetOperation.Child[0]: embedded message failed validation | caused by: invalid SetOperation_Child.ComputedUserset: embedded message failed validation | caused by: invalid ComputedUserset.Relation: value does not match regex pattern \"^[a-z][a-z0-9_]{1,62}[a-z0-9]$\"\n", // Add the expected error message
		},
		{
			name:           "basic parse",
			args:           []string{"parse", "--file-path", "tests/basic.zed"},
			expectedOutput: "Schema validation error: could not lookup definition `sometype` for relation `foo`: object definition `sometype` not found\n", // Add the expected error message
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

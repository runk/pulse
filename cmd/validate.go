package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/runk/pulse/internal/policy"
	"github.com/spf13/cobra"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate a policy file without running it",
	Long: `Validation ensures that the policy file is correctly formatted and can be parsed without errors.
This is useful for catching syntax errors or misconfigurations before attempting to execute the policy.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filename := args[0]

		stdout := cmd.OutOrStdout()

		fmt.Fprintf(stdout, "Validating policy file: %s\n", filename)

		fd, err := os.OpenFile(filename, os.O_RDONLY, 0)
		if err != nil {
			return fmt.Errorf("Cannot open policy file: %w. Check that file exists and is accessible.", err)
		}
		defer fd.Close()

		_, err = policy.ReadPolicy(filename)
		if err != nil {
			errorMsg := color.New(color.FgRed).Sprintf("Policy validation failed: %v", err)
			return errors.New(errorMsg)
		}

		successMsg := color.New(color.FgGreen, color.Bold).Sprint("Success!")
		fmt.Fprintf(stdout, "%s Policy is valid.\n", successMsg)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// validateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// validateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

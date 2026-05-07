package cmd

import (
	"fmt"
	"os"

	"github.com/runk/pulse/internal/policy"
	"github.com/spf13/cobra"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate a policy file without running it",
	Long: `Validation ensures that the policy file is correctly formatted and can be parsed without errors.
This is useful for catching syntax errors or misconfigurations before attempting to execute the policy.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("validate called")

		if len(args) != 1 {
			fmt.Println("Please provide exactly one argument to validate.")
			os.Exit(1)
		}

		filename := args[0]

		_, err := policy.ReadPolicy(filename)
		if err != nil {
			fmt.Println("Error reading policy:", err)
			os.Exit(1)
		}
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

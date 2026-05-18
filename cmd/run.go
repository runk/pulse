package cmd

import (
	"fmt"

	"github.com/runk/pulse/internal/policy"
	"github.com/runk/pulse/internal/runner"
	"github.com/spf13/cobra"
)

var (
	concurrency uint16
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a policy with all its checks",
	Long: `The run command executes the specified policy file, performing all defined checks and actions.
	This command will read the policy file, validate its structure, and then execute the checks in the order they are defined.
You can specify the level of concurrency for running checks using the --concurrency flag.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		stdout := cmd.OutOrStdout()

		if len(args) != 1 {
			return fmt.Errorf("Please provide a single policy file, e.g. 'pulse run policy.json'")
		}

		filename := args[0]
		fmt.Fprintf(stdout, "Using policy file: %s\n", filename)

		policy, err := policy.ReadPolicy(filename)
		if err != nil {
			return fmt.Errorf("Error reading policy: %w", err)
		}

		// fmt.Fprintf(stdout, "Policy loaded: %+v\n", policy)
		err = policy.Validate()
		if err != nil {
			return err
		}

		runner.Execute(policy.Checks, concurrency)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	runCmd.Flags().Uint16VarP(&concurrency, "concurrency", "c", 4, "Level of concurrency")
}

package cmd

import (
	"errors"
	"fmt"
	"io"

	"github.com/fatih/color"
	"github.com/runk/pulse/internal/policy"
	"github.com/runk/pulse/internal/runner"
	"github.com/spf13/cobra"
)

var (
	concurrency uint16
	timeout     uint32
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

		err = policy.Validate()
		if err != nil {
			return err
		}

		errs := make(chan error, 1)
		results := make(chan runner.Result, len(policy.Checks))

		go func() {
			defer close(results)
			errs <- runner.Execute(policy.Checks, results, concurrency, timeout)
		}()

		if ok := formatResults(stdout, results); !ok {
			return errors.New("Some checks failed")
		}

		if err = <-errs; err != nil {
			return err
		}

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
	runCmd.Flags().Uint32VarP(&timeout, "timeout", "t", 3000, "Timeout for each check in milliseconds")
}

func formatResults(stdout io.Writer, results chan runner.Result) bool {
	pass := color.New(color.FgGreen).Sprint("PASS")
	fail := color.New(color.FgRed).Sprint("FAIL")
	ok := true
	for result := range results {
		if !result.Ok {
			ok = false
		}

		typ := color.New(color.FgBlack).Sprint(result.Type)

		if result.Ok {
			fmt.Fprintf(stdout, "%s %s %s\n", pass, typ, result.Subject)
		} else {
			fmt.Fprintf(stdout, "%s %s %s: %s\n", fail, typ, result.Subject, result.Message)
		}
	}

	return ok
}

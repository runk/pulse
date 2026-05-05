package cmd

import (
	"fmt"

	"github.com/runk/pulse/internal/policy"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a policy with all its checks",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("run called")

		if len(args) != 1 {
			fmt.Println("Please provide a single policy file, e.g. 'pulse run policy.json'")
			return
		}

		filename := args[0]
		fmt.Println("Using policy file:", filename)

		policy, err := policy.ReadPolicy(filename)
		if err != nil {
			fmt.Println("Error reading policy:", err)
			return
		}

		fmt.Printf("Policy loaded: %+v\n", policy)
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

	runCmd.Flags().Uint16P("concurrency", "c", 4, "Level of concurrency")
}

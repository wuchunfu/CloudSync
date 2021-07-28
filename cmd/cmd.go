package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wuchunfu/CloudSync/cmd/config"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:          "CloudSync",
	SilenceUsage: true,
	Short:        "Main application",
	Long:         `CloudSync is an cloud sync tool implemented with golang.`,
	Example:      "CloudSync CloudSync",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			tips()
			return fmt.Errorf("requires at least 1 arg(s), only received 0")
		}
		if cmd.Use != args[0] {
			tips()
			return fmt.Errorf("invalid args specified: %s", args[0])
		}
		return nil
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		tips()
	},
}

func init() {
	rootCmd.AddCommand(config.StartCmd)
}

func tips() {
	welcome := `Welcome to CloudSync.`
	help := `You can use -h to view the command.`
	banner := `
_______________               _________________                        
__  ____/___  /______ ____  ________  /__  ___/_____  _________ _______
_  /     __  / _  __ \_  / / /_  __  / _____ \ __  / / /__  __ \_  ___/
/ /___   _  /  / /_/ // /_/ / / /_/ /  ____/ / _  /_/ / _  / / // /__  
\____/   /_/   \____/ \__,_/  \__,_/   /____/  _\__, /  /_/ /_/ \___/  
                                               /____/
    `
	fmt.Printf("%s\n", welcome)
	fmt.Printf("%s\n", banner)
	fmt.Printf("%s\n", help)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

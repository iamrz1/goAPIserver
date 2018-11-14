package cmd

import (
	"fmt"
	"os"

	srvr "github.com/iamrz1/GoApiServer/srvr"
	"github.com/spf13/cobra"
)

var port string
var verbose bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "GoApiServer",
	Short: "API endpints for CRUD operations",
	Long:  `Use flags to change the default port`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		srvr.PostMain(port, verbose)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func init() {

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&port, "port", "8080", "Set the port name for the API server. Default 8080")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Set verbose mode on for the API server. Default is false.")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

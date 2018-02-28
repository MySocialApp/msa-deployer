package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	log "github.com/sirupsen/logrus"
)

var cfgFile string
var clientFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:	"deployer",
	Short:	"MySocialApp deployer for application and databases",
	Long:	`This application is used manage client's infrastructure`,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of msa-deployer",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("v0.2")
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
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(versionCmd)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./.deployer.yaml)")
	rootCmd.PersistentFlags().StringVar(&clientFile, "clientfile", "", "config file (default is ./clients.csv)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in current directory
		viper.AddConfigPath(".")
		viper.SetConfigName(".deployer")
	}
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig() ; err != nil {
		log.Error("Config file not found:")
		log.Fatal(err)
	}
	log.Debug("Using config file: ", viper.ConfigFileUsed())

	configValidator("gitlab_project_id")
	configValidator("gitlab_token")

}

func configValidator(key string) bool {
	if len(viper.GetString("gitlab_project_id")) == 0 {
		log.Error("Can't access mandatory information in your config file, please set '" + key + "' in " + viper.ConfigFileUsed())
		os.Exit(1)
	}
	return true
}
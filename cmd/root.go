package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"path/filepath"
	log "github.com/sirupsen/logrus"
)

var cfgFile string
var clientFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:	"msa-k8s-deployer",
	Short:	"MySocialApp deployer for application and databases",
	Long:	`This application is used manage client's infrastructure`,
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
		// Find current directory.
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in current directory with name ".msa-k8s-deployer" (without extension).
		viper.AddConfigPath(dir)
		viper.SetConfigName(".deployer.yaml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Info("Using config file: ", viper.ConfigFileUsed())
	} else {
		log.Error("Config file not found ")
		os.Exit(1)
	}

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
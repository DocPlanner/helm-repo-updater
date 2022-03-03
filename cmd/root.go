package cmd

import (
	"fmt"
	"os"
	"text/template"

	"github.com/docplanner/helm-repo-updater/internal/app/git"
	"github.com/docplanner/helm-repo-updater/internal/app/updater"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// HelmUpdaterConfig contains global configuration and required runtime data
type HelmUpdaterConfig struct {
	DryRun           bool
	LogLevel         string
	AppName          string
	UpdateApps       []updater.ChangeEntry
	GitCommitMessage *template.Template
	GitCredentials   *git.Credentials
	GitConf          *git.Conf
}

// GitConf contains the configuration for the git repository
type GitConf struct {
	RepoURL string
	Branch  string
	File    string
}

// UpdateApp contains the information about the update app
type UpdateApp struct {
	Name  string
	File  string
	Key   string
	Image string
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "helm-repo-updater",
	Short: "Helm repo updater",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.helm-repo-updater.yaml)")

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
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".helm-repo-updater" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".helm-repo-updater")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

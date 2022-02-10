package cmd

import (
	"fmt"
	"github.com/docplanner/helm-repo-updater/internal/app/logger"
	"os"
	"path"
	"text/template"

	"github.com/docplanner/helm-repo-updater/internal/app/git"
	"github.com/docplanner/helm-repo-updater/internal/app/updater"
	"github.com/spf13/cobra"
)

const (
	GitCommitUser  = "git-commit-user"
	GitCommitEmail = "git-commit-email"
	GitPassword    = "git-password"
	GitBranch      = "git-branch"
	GitRepoUrl     = "git-repo-url"
	GitFile        = "git-file"
	GitDir         = "git-dir"

	AppName          = "app-name"
	SshPrivateKey    = "ssh-private-key"
	DryRun           = "dry-run"
	LogLevel         = "logLevel"
	LogStashProtocol = "log-stash-protocol"
	LogStashAddr     = "log-stash-addr"
	HelmKeyValues    = "helm-key-values"
)

var cfg = updater.HelmUpdaterConfig{}

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs the helm repo updater",
	Run: func(cmd *cobra.Command, args []string) {
		gitUser, _ := cmd.Flags().GetString(GitCommitUser)
		gitEmail, _ := cmd.Flags().GetString(GitCommitEmail)
		gitPass, _ := cmd.Flags().GetString(GitPassword)
		gitBranch, _ := cmd.Flags().GetString(GitBranch)
		gitRepoURL, _ := cmd.Flags().GetString(GitRepoUrl)
		gitFile, _ := cmd.Flags().GetString(GitFile)
		gitDir, _ := cmd.Flags().GetString(GitDir)
		sshKey, _ := cmd.Flags().GetString(SshPrivateKey)
		appName, _ := cmd.Flags().GetString(AppName)
		logLevel, _ := cmd.Flags().GetString(LogLevel)
		logStashProtocol, _ := cmd.Flags().GetString(LogStashProtocol)
		logStashAddress, _ := cmd.Flags().GetString(LogStashAddr)
		dryRun, _ := cmd.Flags().GetBool(DryRun)
		helmKVs, _ := cmd.Flags().GetStringToString(HelmKeyValues)

		if len(helmKVs) == 0 {
			if err := cmd.Help(); err != nil {
				return
			}

			os.Exit(1)
		}

		var updateApps []updater.ChangeEntry
		for k, v := range helmKVs {
			updateApps = append(updateApps, updater.ChangeEntry{
				Key:      k,
				NewValue: v,
			})
		}

		gitCredentials := &git.Credentials{
			Username:   gitUser,
			Email:      gitEmail,
			Password:   gitPass,
			SSHPrivKey: sshKey,
		}

		gitConf := &git.Conf{
			RepoURL: gitRepoURL,
			Branch:  gitBranch,
		}

		logger := loadLogger(logLevel, logStashProtocol, logStashAddress)

		if tpl, err := template.New("commitMessage").Parse(git.DefaultGitCommitMessage); err != nil {
			logger.Fatal("could not parse commit message template", err)

			return
		} else {
			logger.Info("Successfully parsed commit message template")

			gitConf.Message = tpl
		}

		cfg = updater.HelmUpdaterConfig{
			DryRun:         dryRun,
			AppName:        appName,
			UpdateApps:     updateApps,
			File:           path.Join(gitDir, appName, gitFile),
			GitCredentials: gitCredentials,
			GitConf:        gitConf,
			Logger:         logger,
		}

		if err := runImageUpdater(cfg); err != nil {
			logger.Error("Can not run image updater", err)
		}
	},
}

func loadLogger(logLevel, logStashProtocol, logStashAddr string) logger.Logger {
	if logStashProtocol == "" || logStashAddr == "" {
		z := logger.NewConsoleZeroLogger(logLevel)

		return logger.NewZeroLogger(z)
	}

	z := logger.NewConsoleELKZeroLogger(logLevel, logStashProtocol, logStashAddr)

	return logger.NewZeroLogger(z)
}

func runImageUpdater(cfg updater.HelmUpdaterConfig) error {
	syncState := updater.NewSyncIterationState()

	cfg.Logger.DebugWithContext("Updating values", map[string]interface{}{
		"changes": cfg.UpdateApps,
	})

	return func(cfg updater.HelmUpdaterConfig) error {
		cfg.Logger.DebugWithContext("Processing application", map[string]interface{}{
			"application": cfg.AppName,
			"file":        cfg.File,
		})

		_, err := updater.UpdateApplication(cfg, syncState)
		if err != nil {
			cfg.Logger.Error(
				fmt.Sprintf("Can not update application: %s", cfg.AppName),
				err,
			)

			return err
		}

		cfg.Logger.Info(
			fmt.Sprintf("Application %s successfully updated", cfg.AppName),
		)

		return nil
	}(cfg)
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().String(GitCommitUser, "", "Username to use for Git commits")
	runCmd.Flags().String(GitCommitEmail, "", "E-Mail address to use for Git commits")
	runCmd.Flags().String(GitPassword, "", "Password for github user")
	runCmd.Flags().String(GitBranch, "develop", "branch")
	runCmd.Flags().String(GitRepoUrl, "", "git repo url")
	runCmd.Flags().String(GitFile, "", "file eg. values.yaml")
	runCmd.Flags().String(GitDir, "", "file eg. /production/charts/")
	runCmd.Flags().String(AppName, "", "app name")
	runCmd.Flags().String(SshPrivateKey, "", "ssh private key (only using ")
	runCmd.Flags().Bool(DryRun, false, "run in dry-run mode. If set to true, do not perform any changes")
	runCmd.Flags().String(LogLevel, "info", "set the loglevel to one of trace|debug|info|warn|error")
	runCmd.Flags().String(LogStashProtocol, "", "set the protocol like udp/tcp for LogStash (optional)")
	runCmd.Flags().String(LogStashAddr, "", "set the http address to achieve logstash (optional)")
	runCmd.Flags().StringToString(HelmKeyValues, nil, "helm key-values sets")

	_ = runCmd.MarkFlagRequired(GitCommitUser)
	_ = runCmd.MarkFlagRequired(GitCommitEmail)
	_ = runCmd.MarkFlagRequired(GitRepoUrl)
	_ = runCmd.MarkFlagRequired(GitFile)
	_ = runCmd.MarkFlagRequired(GitDir)
	_ = runCmd.MarkFlagRequired(HelmKeyValues)
	_ = runCmd.MarkFlagRequired(AppName)
}

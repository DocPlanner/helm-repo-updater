package cmd

import (
	"fmt"
	"os"
	"path"
	"text/template"

	"github.com/argoproj-labs/argocd-image-updater/pkg/log"
	"github.com/docplanner/helm-repo-updater/internal/app/git"
	"github.com/docplanner/helm-repo-updater/internal/app/updater"
	"github.com/spf13/cobra"
)

var cfg = updater.HelmUpdaterConfig{}

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs the helm repo updater",
	Run: func(cmd *cobra.Command, args []string) {
		gitUser, _ := cmd.Flags().GetString("git-commit-user")
		gitEmail, _ := cmd.Flags().GetString("git-commit-email")
		gitPass, _ := cmd.Flags().GetString("git-password")
		gitBranch, _ := cmd.Flags().GetString("git-branch")
		gitRepoURL, _ := cmd.Flags().GetString("git-repo-url")
		gitFile, _ := cmd.Flags().GetString("git-file")
		gitDir, _ := cmd.Flags().GetString("git-dir")
		sshKey, _ := cmd.Flags().GetString("ssh-private-key")
		appName, _ := cmd.Flags().GetString("app-name")
		logLevel, _ := cmd.Flags().GetString("loglevel")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		helmKVs, _ := cmd.Flags().GetStringToString("helm-key-values")

		if err := log.SetLogLevel(logLevel); err != nil {
			fmt.Println(err)

			os.Exit(1)
		}

		if len(helmKVs) == 0 {
			err := cmd.Help()
			if err != nil {
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

		if tpl, err := template.New("commitMessage").Parse(git.DefaultGitCommitMessage); err != nil {
			log.Fatalf("could not parse commit message template: %v", err)

			return
		} else {
			log.Debugf("Successfully parsed commit message template")

			gitConf.Message = tpl
		}

		cfg = updater.HelmUpdaterConfig{
			DryRun:         dryRun,
			LogLevel:       logLevel,
			AppName:        appName,
			UpdateApps:     updateApps,
			File:           path.Join(gitDir, appName, gitFile),
			GitCredentials: gitCredentials,
			GitConf:        gitConf,
		}

		err := runImageUpdater(cfg)
		if err != nil {
			log.Errorf("Error: %v", err)
		}

	},
}

func runImageUpdater(cfg updater.HelmUpdaterConfig) error {

	syncState := updater.NewSyncIterationState()

	fmt.Println(cfg.UpdateApps)

	err := func(cfg updater.HelmUpdaterConfig) error {
		log.Debugf("Processing application %s in directory %s", cfg.AppName, cfg.File)

		err := updater.UpdateApplication(cfg, syncState)
		if err != nil {
			return err
		}

		return nil
	}(cfg)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().String("git-commit-user", "", "Username to use for Git commits")
	runCmd.Flags().String("git-commit-email", "", "E-Mail address to use for Git commits")
	runCmd.Flags().String("git-password", "", "Password for github user")
	runCmd.Flags().String("git-branch", "develop", "branch")
	runCmd.Flags().String("git-repo-url", "", "git repo url")
	runCmd.Flags().String("git-file", "", "file eg. values.yaml")
	runCmd.Flags().String("git-dir", "", "file eg. /production/charts/")
	runCmd.Flags().String("app-name", "", "app name")
	runCmd.Flags().String("ssh-private-key", "", "ssh private key (only using ")
	runCmd.Flags().Bool("dry-run", false, "run in dry-run mode. If set to true, do not perform any changes")
	runCmd.Flags().String("loglevel", "info", "set the loglevel to one of trace|debug|info|warn|error")
	runCmd.Flags().StringToString("helm-key-values", nil, "helm key-values sets")

	_ = runCmd.MarkFlagRequired("git-commit-user")
	_ = runCmd.MarkFlagRequired("git-commit-email")
	_ = runCmd.MarkFlagRequired("git-repo-url")
	_ = runCmd.MarkFlagRequired("git-file")
	_ = runCmd.MarkFlagRequired("git-dir")
	_ = runCmd.MarkFlagRequired("helm-key-values")
	_ = runCmd.MarkFlagRequired("app-name")

}

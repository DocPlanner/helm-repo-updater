package cmd

import (
	"fmt"
	"os"
	"path"
	"text/template"

	"github.com/docplanner/helm-repo-updater/internal/app/git"
	"github.com/docplanner/helm-repo-updater/internal/app/log"
	"github.com/docplanner/helm-repo-updater/internal/app/updater"
	"github.com/spf13/cobra"
)

const (
	// GitCommitUser is the user login used for commit changes
	GitCommitUser = "git-commit-user"
	// GitCommitEmail is the email used for commit changes
	GitCommitEmail = "git-commit-email"
	// GitPassword is the git password used for auth
	GitPassword = "git-password"
	// GitBranch is the branch of the git repository
	GitBranch = "git-branch"
	// GitRepoURL is the git repository url
	GitRepoURL = "git-repo-url"
	// GitFile is the file that is going to be changed
	GitFile = "git-file"
	// GitDir is the directory where the file to be changed is located
	GitDir = "git-dir"
	// AppName is the name of the helm application
	AppName = "app-name"
	// SSHPrivateKey is the location of the SSH private key used for auth
	SSHPrivateKey = "ssh-private-key"
	// UseSSHPrivateKeyAsInline indicates if the SSHPrivateKey is going to be created based in a string provided
	UseSSHPrivateKeyAsInline = "use-ssh-private-key-as-inline"
	// DryRun is going to indicate if the changes are going to be committed or not
	DryRun = "dry-run"
	// LogLevel will indicate the log level
	LogLevel = "logLevel"
	// HelmKeyValues will be used for indicate the key and values to be changed in helm
	HelmKeyValues = "helm-key-values"
	// AllowErrorNothingToUpdate represents that is allowed the error nothing to update
	AllowErrorNothingToUpdate = "allow-nothing-to-update"
	// AllowErrorNothingToUpdateMessage represents the allowed error that will be the exception for make an os.Exit(1) call when is detected
	AllowErrorNothingToUpdateMessage = "nothing to update, skipping commit"
)

var cfg = updater.HelmUpdaterConfig{}

// runImageUpdater checks and apply the necessary update in the helm application
func runImageUpdater(cfg updater.HelmUpdaterConfig) error {

	syncState := updater.NewSyncIterationState()

	err := func(cfg updater.HelmUpdaterConfig) error {
		log.Debugf("Processing application %s in directory %s", cfg.AppName, cfg.File)

		_, err := updater.UpdateApplication(cfg, syncState)
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

// checkExecutionRunImageUpdater represents the check of the execution of the runImageUpdater command
func checkExecutionRunImageUpdater(cfg updater.HelmUpdaterConfig, logCtx *log.Context, appName string) {
	if err := runImageUpdater(cfg); err != nil {
		if err.Error() != AllowErrorNothingToUpdateMessage || !cfg.AllowErrorNothingToUpdate {
			logCtx.Errorf("Error trying to update the %s application: %v", appName, err)
			os.Exit(1)
		}
		logCtx.Infof("%s", err.Error())
		return
	}
}

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs the helm repo updater",
	Run: func(cmd *cobra.Command, args []string) {
		gitUser, _ := cmd.Flags().GetString(GitCommitUser)
		gitEmail, _ := cmd.Flags().GetString(GitCommitEmail)
		gitPass, _ := cmd.Flags().GetString(GitPassword)
		gitBranch, _ := cmd.Flags().GetString(GitBranch)
		gitRepoURL, _ := cmd.Flags().GetString(GitRepoURL)
		gitFile, _ := cmd.Flags().GetString(GitFile)
		gitDir, _ := cmd.Flags().GetString(GitDir)
		sshKey, _ := cmd.Flags().GetString(SSHPrivateKey)
		appName, _ := cmd.Flags().GetString(AppName)
		logLevel, _ := cmd.Flags().GetString(LogLevel)
		dryRun, _ := cmd.Flags().GetBool(DryRun)
		useSSHPrivateKeyAsInline, _ := cmd.Flags().GetBool(UseSSHPrivateKeyAsInline)
		helmKVs, _ := cmd.Flags().GetStringToString(HelmKeyValues)
		allowErrorNothingToUpdate, _ := cmd.Flags().GetBool(AllowErrorNothingToUpdate)

		if err := log.SetLogLevel(logLevel); err != nil {
			fmt.Println(err)

			os.Exit(1)
		}

		if len(helmKVs) == 0 {
			if err := cmd.Help(); err != nil {
				return
			}

			os.Exit(1)
		}

		var updateApps []updater.ChangeEntry
		var tpl *template.Template
		var err error
		for k, v := range helmKVs {
			updateApps = append(updateApps, updater.ChangeEntry{
				Key:      k,
				NewValue: v,
			})
		}

		gitCredentials := &git.Credentials{
			Username:             gitUser,
			Email:                gitEmail,
			Password:             gitPass,
			SSHPrivKey:           sshKey,
			SSHPrivKeyFileInline: useSSHPrivateKeyAsInline,
		}

		gitConf := &git.Conf{
			RepoURL: gitRepoURL,
			Branch:  gitBranch,
		}

		logCtx := log.WithContext().AddField("application", appName)

		if tpl, err = template.New("commitMessage").Parse(git.DefaultGitCommitMessage); err != nil {
			logCtx.Fatalf("could not parse commit message template: %v", err)

			return
		}
		logCtx.Debugf("Successfully parsed commit message template")
		gitConf.Message = tpl

		cfg = updater.HelmUpdaterConfig{
			DryRun:                    dryRun,
			LogLevel:                  logLevel,
			AppName:                   appName,
			UpdateApps:                updateApps,
			File:                      path.Join(gitDir, appName, gitFile),
			GitCredentials:            gitCredentials,
			GitConf:                   gitConf,
			AllowErrorNothingToUpdate: allowErrorNothingToUpdate,
		}

		checkExecutionRunImageUpdater(cfg, logCtx, appName)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().String(GitCommitUser, "", "Username to use for Git commits")
	runCmd.Flags().String(GitCommitEmail, "", "e-mail address to use for Git commits")
	runCmd.Flags().String(GitPassword, "", "Password for github user")
	runCmd.Flags().String(GitBranch, "develop", "git repo branch")
	runCmd.Flags().String(GitRepoURL, "", "git repo url")
	runCmd.Flags().String(GitFile, "", "file eg. values.yaml")
	runCmd.Flags().String(GitDir, "", "file eg. /production/charts/")
	runCmd.Flags().String(AppName, "", "app name")
	runCmd.Flags().String(SSHPrivateKey, "", "ssh private key")
	runCmd.Flags().Bool(UseSSHPrivateKeyAsInline, false, "ssh private key inline creation, if true it will use ssh-private-key as input for create ssh private key file in temporal directory")
	runCmd.Flags().Bool(DryRun, false, "run in dry-run mode. If set to true, do not perform any changes")
	runCmd.Flags().String(LogLevel, "info", "set the loglevel to one of trace|debug|info|warn|error")
	runCmd.Flags().StringToString(HelmKeyValues, nil, "helm key-values sets")
	runCmd.Flags().Bool(AllowErrorNothingToUpdate, true, "allow the error message 'nothing to update, skipping commit' and finish without exit 1 the execution")

	_ = runCmd.MarkFlagRequired(GitCommitUser)
	_ = runCmd.MarkFlagRequired(GitCommitEmail)
	_ = runCmd.MarkFlagRequired(GitRepoURL)
	_ = runCmd.MarkFlagRequired(GitFile)
	_ = runCmd.MarkFlagRequired(HelmKeyValues)
	_ = runCmd.MarkFlagRequired(AppName)

}

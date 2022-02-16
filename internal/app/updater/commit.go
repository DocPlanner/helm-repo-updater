package updater

import (
	"fmt"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/argoproj-labs/argocd-image-updater/ext/git"
	git_hru "github.com/docplanner/helm-repo-updater/internal/app/git"
	"github.com/docplanner/helm-repo-updater/internal/app/log"
)

const (
	// default origin branch
	defaultOriginBranch = "origin"
)

// UpdateApplication update all values of a single application.
func UpdateApplication(cfg HelmUpdaterConfig, state *SyncIterationState) (*[]ChangeEntry, error) {
	appsChanges, err := commitChangesLocked(cfg, state)
	if err != nil {
		log.Errorf("Could not update application spec: %v", err)

		return nil, err
	}

	log.Infof("Successfully updated the live application spec")

	return appsChanges, nil

}

// commitChangesLocked commits the changes to the git repository.
func commitChangesLocked(cfg HelmUpdaterConfig, state *SyncIterationState) (*[]ChangeEntry, error) {
	lock := state.GetRepositoryLock(cfg.GitConf.RepoURL)
	lock.Lock()
	defer lock.Unlock()

	return commitChangesGit(cfg, writeOverrides)
}

// configureGitClientCredentials set username and e-mail address in gitClient to identify the committer
func configureGitClientCredentials(gitC git.Client, gitCredentials git_hru.Credentials) (git.Client, error) {
	if gitCredentials.Username != "" && gitCredentials.Email != "" {
		err := gitC.Config(gitCredentials.Username, gitCredentials.Email)
		if err != nil {
			return nil, err
		}
	}
	return gitC, nil
}

// initAndFetchGitRepository initializes a local git repository and sets the remote origin
// and fetches latest updates from origin
func initAndFetchGitRepository(repoUrl string, tempRoot string, creds git.Creds, gitCredentials git_hru.Credentials) (git.Client, error) {
	gitC, err := git.NewClientExt(repoUrl, tempRoot, creds, true, false, "")
	if err != nil {
		return nil, err
	}

	err = gitC.Init()
	if err != nil {
		return nil, err
	}

	err = gitC.Fetch("")
	if err != nil {
		return nil, err
	}

	// Set username and e-mail address used to identify the commiter
	if gitCredentials.Username != "" && gitCredentials.Email != "" {
		err = gitC.Config(gitCredentials.Username, gitCredentials.Email)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("it's necessary provide an username and an email for configure git client")
	}

	gitC, err = configureGitClientCredentials(gitC, gitCredentials)
	if err != nil {
		return nil, err
	}

	return gitC, nil
}

// commitAndPushGitChanges perfoms a git commit for the given pathSpec to the currently checked
// out branch and after pushes local changes to the remote branch
func commitAndPushGitChanges(gitC git.Client, opts *git.CommitOptions, checkOutBranch string) error {
	err := gitC.Commit("", opts)
	if err != nil {
		return err
	}
	err = gitC.Push(defaultOriginBranch, checkOutBranch, false)
	if err != nil {
		return err
	}
	return nil
}

// configureCommitOptions creates a git.CommitOptions based in the appName the apps to
// change and the helm repo updater config message template generated
func configureCommitOptions(appName string, apps []ChangeEntry, helmUpdaterConfigMessage *template.Template) (*git.CommitOptions, error) {
	var gitCommitMessage string

	logCtx := log.WithContext().AddField("application", appName)

	if len(apps) > 0 && helmUpdaterConfigMessage != nil {
		gitCommitMessage = TemplateCommitMessage(helmUpdaterConfigMessage, appName, apps)
	}

	commitOpts := &git.CommitOptions{}
	if gitCommitMessage != "" {
		cm, err := ioutil.TempFile("", "image-updater-commit-msg")
		if err != nil {
			return nil, fmt.Errorf("cold not create temp file: %v", err)
		}
		logCtx.Debugf("Writing commit message to %s", cm.Name())
		err = ioutil.WriteFile(cm.Name(), []byte(gitCommitMessage), 0600)
		if err != nil {
			_ = cm.Close()
			return nil, fmt.Errorf("could not write commit message to %s: %v", cm.Name(), err)
		}
		commitOpts.CommitMessagePath = cm.Name()
		_ = cm.Close()
		defer os.Remove(cm.Name())
	}
	return commitOpts, nil
}

// createTempFileInDirectory creates a temporal directory where a copy of
// the git repository is going to be stored.
func createTempFileInDirectory(dirName string, applicationName string) (*string, error) {
	logCtx := log.WithContext().AddField("application", applicationName)
	tempRoot, err := ioutil.TempDir(os.TempDir(), dirName)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := os.RemoveAll(tempRoot)
		if err != nil {
			logCtx.Errorf("could not remove temp dir: %v", err)
		}
	}()
	return &tempRoot, nil
}

func getCheckoutBranchAndSystemRef(gitConfBranch string, applicationName string, gitC git.Client) (*string, error) {
	logCtx := log.WithContext().AddField("application", applicationName)
	checkOutBranch := gitConfBranch

	logCtx.Tracef("targetRevision for update is '%s'", checkOutBranch)

	if checkOutBranch == "" || checkOutBranch == "HEAD" {
		checkOutBranch, err := gitC.SymRefToBranch(checkOutBranch)
		logCtx.Infof("resolved remote default branch to '%s' and using that for operations", checkOutBranch)
		if err != nil {
			return nil, err
		}
	}
	return &checkOutBranch, nil
}

// commitChangesGit commits any changes required for updating one or more values
// after the UpdateApplication cycle has finished.
func commitChangesGit(cfg HelmUpdaterConfig, write changeWriter) (*[]ChangeEntry, error) {
	var apps []ChangeEntry

	logCtx := log.WithContext().AddField("application", cfg.AppName)

	creds, err := cfg.GitCredentials.NewGitCreds(cfg.GitConf.RepoURL)
	if err != nil {
		return nil, fmt.Errorf("could not get creds for repo '%s': %v", cfg.AppName, err)
	}
	var gitC git.Client

	tempRoot, err := createTempFileInDirectory(fmt.Sprintf("git-%s", cfg.AppName), cfg.AppName)
	if err != nil {
		return nil, err
	}

	gitC, err = initAndFetchGitRepository(cfg.GitConf.RepoURL, *tempRoot, creds, *cfg.GitCredentials)
	if err != nil {
		return nil, err
	}

	checkOutBranch, err := getCheckoutBranchAndSystemRef(cfg.GitConf.Branch, cfg.AppName, gitC)
	if err != nil {
		return nil, err
	}

	err = gitC.Checkout(*checkOutBranch)
	if err != nil {
		return nil, err
	}

	// write changes to files
	if apps, err = write(cfg, gitC); err != nil {
		return nil, err
	}

	commitOpts, err := configureCommitOptions(cfg.AppName, apps, cfg.GitConf.Message)
	if err != nil {
		return nil, err
	}

	if cfg.DryRun {
		logCtx.Infof("dry run, not committing changes")
		return &apps, nil
	}

	err = commitAndPushGitChanges(gitC, commitOpts, *checkOutBranch)
	if err != nil {
		return nil, err
	}

	return &apps, nil
}

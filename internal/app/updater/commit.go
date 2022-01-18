package updater

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/argoproj-labs/argocd-image-updater/ext/git"
	"github.com/argoproj-labs/argocd-image-updater/pkg/log"
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

// commitChangesGit commits any changes required for updating one or more values
// after the UpdateApplication cycle has finished.
func commitChangesGit(cfg HelmUpdaterConfig, write changeWriter) (*[]ChangeEntry, error) {
	var apps []ChangeEntry
	var gitCommitMessage string

	logCtx := log.WithContext().AddField("application", cfg.AppName)

	creds, err := cfg.GitCredentials.NewGitCreds(cfg.GitConf.RepoURL)
	if err != nil {
		return nil, fmt.Errorf("could not get creds for repo '%s': %v", cfg.AppName, err)
	}
	var gitC git.Client
	tempRoot, err := ioutil.TempDir(os.TempDir(), fmt.Sprintf("git-%s", cfg.AppName))
	if err != nil {
		return nil, err
	}
	defer func() {
		err := os.RemoveAll(tempRoot)
		if err != nil {
			logCtx.Errorf("could not remove temp dir: %v", err)
		}
	}()

	gitC, err = git.NewClientExt(cfg.GitConf.RepoURL, tempRoot, creds, false, false, "")
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
	if cfg.GitCredentials.Username != "" && cfg.GitCredentials.Email != "" {
		err = gitC.Config(cfg.GitCredentials.Username, cfg.GitCredentials.Email)
		if err != nil {
			return nil, err
		}
	}

	checkOutBranch := cfg.GitConf.Branch

	logCtx.Tracef("targetRevision for update is '%s'", checkOutBranch)
	if checkOutBranch == "" || checkOutBranch == "HEAD" {
		checkOutBranch, err = gitC.SymRefToBranch(checkOutBranch)
		logCtx.Infof("resolved remote default branch to '%s' and using that for operations", checkOutBranch)
		if err != nil {
			return nil, err
		}
	}

	err = gitC.Checkout(checkOutBranch)
	if err != nil {
		return nil, err
	}

	// write changes to files
	if apps, err = write(cfg, gitC); err != nil {
		return nil, err
	}

	commitOpts := &git.CommitOptions{}
	if len(apps) > 0 && cfg.GitConf.Message != nil {
		gitCommitMessage = TemplateCommitMessage(cfg.GitConf.Message, cfg.AppName, apps)
	}

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

	if cfg.DryRun {
		logCtx.Infof("dry run, not committing changes")

		return &apps, nil
	}

	err = gitC.Commit("", commitOpts)
	if err != nil {
		return nil, err
	}
	err = gitC.Push("origin", checkOutBranch, false)
	if err != nil {
		return nil, err
	}

	return &apps, nil
}

package updater

import (
	"fmt"
	"io/ioutil"
	"os"
	"text/template"
	"time"

	"github.com/docplanner/helm-repo-updater/internal/app/log"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
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

// cloneRepository clones the git repository in a temporal directory.
func cloneRepository(appName string, repoUrl string, authCreds transport.AuthMethod, tempRoot string) (*git.Repository, error) {
	logCtx := log.WithContext().AddField("application", appName)
	logCtx.Infof("git clone %s ", repoUrl)
	r, err := git.PlainClone(tempRoot, false, &git.CloneOptions{
		Auth:     authCreds,
		URL:      repoUrl,
		Progress: os.Stdout,
	})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// commitAndPushGitChanges perfoms a git commit for the given pathSpec to the currently checked
// out branch and after pushes local changes to the remote branch
func commitAndPushGitChanges(appName string, commitMessage string, gitW git.Worktree, gitUserName string, gitEmail string, tempRoot string, gitAuth transport.AuthMethod) error {
	logCtx := log.WithContext().AddField("application", appName)
	logCtx.Infof("git commit -m %s ", commitMessage)
	commit, err := gitW.Commit(commitMessage, &git.CommitOptions{
		Author: &object.Signature{
			Name:  gitUserName,
			Email: gitEmail,
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}
	r, err := git.PlainOpen(tempRoot)
	if err != nil {
		return err
	}
	// Prints the current HEAD to verify that all worked well.
	logCtx.Debugf("git show -s")
	obj, err := r.CommitObject(commit)
	if err != nil {
		return err
	}
	logCtx.Infof("obj: %s", obj)

	logCtx.Infof("git push")
	// push using default options
	err = r.Push(&git.PushOptions{
		Auth: gitAuth,
	})
	if err != nil {
		return err
	}
	return nil
}

// configureCommitOptions creates a git.CommitOptions based in the appName the apps to
// change and the helm repo updater config message template generated
func configureCommitMessage(appName string, apps []ChangeEntry, helmUpdaterConfigMessage *template.Template) (*string, error) {
	var gitCommitMessage string

	logCtx := log.WithContext().AddField("application", appName)

	if len(apps) > 0 && helmUpdaterConfigMessage != nil {
		gitCommitMessage = TemplateCommitMessage(helmUpdaterConfigMessage, appName, apps)
	}

	if gitCommitMessage != "" {
		cm, err := ioutil.TempFile("", appName)
		if err != nil {
			return nil, fmt.Errorf("cold not create temp file: %v", err)
		}
		logCtx.Debugf("Writing commit message to %s", cm.Name())
		err = ioutil.WriteFile(cm.Name(), []byte(gitCommitMessage), 0600)
		if err != nil {
			_ = cm.Close()
			return nil, fmt.Errorf("could not write commit message to %s: %v", cm.Name(), err)
		}
		gitCommitMessage = cm.Name()
		_ = cm.Close()
		defer os.Remove(cm.Name())
	}
	return &gitCommitMessage, nil
}

// createTempFileInDirectory creates a temporal directory where a copy of
// the git repository is going to be stored.
func createTempFileInDirectory(dirName string, applicationName string) (*string, error) {
	logCtx := log.WithContext().AddField("application", applicationName)
	tempRoot, err := ioutil.TempDir(os.TempDir(), dirName)
	if err != nil {
		return nil, err
	}
	logCtx.Debugf("Created temporal directory to clone repository %s", tempRoot)
	defer func() {
		err := os.RemoveAll(tempRoot)
		if err != nil {
			logCtx.Errorf("could not remove temp dir: %v", err)
		}
	}()
	return &tempRoot, nil
}

// checkBranchExists check if a specific branch in a repository was already created in the origin
func checkBranchExists(gitR git.Repository, checkOutBranch plumbing.ReferenceName) error {
	_, err := gitR.ResolveRevision(plumbing.Revision(checkOutBranch))
	if err != nil {
		return err
	}
	return nil
}

func getCheckoutName(gitConfBranch string, applicationName string, gitR git.Repository) (*plumbing.ReferenceName, error) {
	var checkOutBranch plumbing.ReferenceName
	logCtx := log.WithContext().AddField("application", applicationName)

	logCtx.Tracef("targetRevision for update is '%s'", checkOutBranch)

	if gitConfBranch == "" || gitConfBranch == "HEAD" {
		// retrieving the branch being pointed by head
		ref, err := gitR.Head()
		if err != nil {
			return nil, err
		}
		checkOutBranch = ref.Name()
		return &checkOutBranch, nil
	}
	checkOutBranch = plumbing.NewBranchReferenceName(gitConfBranch)
	return &checkOutBranch, nil
}

func getGitRepositoryWithBranchUpdated(gitW git.Worktree, gitR git.Repository,
	creds transport.AuthMethod, gitConfBranch string, appName string) (*git.Worktree, error) {
	checkOutBranch, err := getCheckoutName(gitConfBranch, appName, gitR)
	if err != nil {
		return nil, err
	}

	err = gitW.Checkout(&git.CheckoutOptions{
		Branch: *checkOutBranch,
	})
	if err != nil {
		return nil, err
	}

	err = checkBranchExists(gitR, *checkOutBranch)
	if err != nil {
		return nil, err
	}

	// Pull the latest changes from the origin remote and merge into the current branch
	log.Infof("git pull origin")
	err = gitW.Pull(&git.PullOptions{
		RemoteName: "origin",
		Auth:       creds,
	})
	if err != nil {
		if err.Error() != "already up-to-date" {
			return nil, err
		}
	}

	// Print the latest commit that was just pulled
	ref, err := gitR.Head()
	if err != nil {
		return nil, err
	}
	commit, err := gitR.CommitObject(ref.Hash())
	if err != nil {
		return nil, err
	}
	log.Debugf("The latest commit to the branch %s is %s", checkOutBranch, commit)
	return &gitW, nil
}

// commitChangesGit commits any changes required for updating one or more values
// after the UpdateApplication cycle has finished.
func commitChangesGit(cfg HelmUpdaterConfig, write changeWriter) (*[]ChangeEntry, error) {
	var apps []ChangeEntry

	logCtx := log.WithContext().AddField("application", cfg.AppName)

	creds, err := cfg.GitCredentials.NewGitCreds(cfg.GitConf.RepoURL, cfg.GitCredentials.Password)
	if err != nil {
		return nil, fmt.Errorf("could not get creds for repo '%s': %v", cfg.AppName, err)
	}

	tempRoot, err := createTempFileInDirectory(fmt.Sprintf("git-%s", cfg.AppName), cfg.AppName)
	if err != nil {
		return nil, err
	}

	gitR, err := cloneRepository(cfg.AppName, cfg.GitConf.RepoURL, creds, *tempRoot)
	if err != nil {
		return nil, err
	}

	gitW, err := gitR.Worktree()
	if err != nil {
		return nil, err
	}

	gitW, err = getGitRepositoryWithBranchUpdated(*gitW, *gitR, creds, cfg.GitConf.Branch, cfg.AppName)
	if err != nil {
		return nil, err
	}

	// write changes to files
	if apps, err = write(cfg, *tempRoot, *gitW); err != nil {
		return nil, err
	}

	commitMessage, err := configureCommitMessage(cfg.AppName, apps, cfg.GitConf.Message)
	if err != nil {
		return nil, err
	}

	if cfg.DryRun {
		logCtx.Infof("dry run, not committing changes")
		return &apps, nil
	}

	err = commitAndPushGitChanges(cfg.AppName, *commitMessage, *gitW, cfg.GitCredentials.Username, cfg.GitCredentials.Email, *tempRoot, creds)
	if err != nil {
		return nil, err
	}

	return &apps, nil
}

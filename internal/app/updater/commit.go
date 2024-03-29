package updater

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"text/template"
	"time"

	git_internal "github.com/docplanner/helm-repo-updater/internal/app/git"
	"github.com/docplanner/helm-repo-updater/internal/app/log"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
)

// UpdateApplication update all values of a single application.
func UpdateApplication(cfg HelmUpdaterConfig, state *SyncIterationState) (*[]ChangeEntry, error) {
	logCtx := log.WithContext().AddField("application", cfg.AppName)
	appsChanges, err := commitChangesLocked(cfg, state)
	if err != nil {
		logCtx.Errorf("Could not update application spec: %v", err)

		return nil, err
	}

	logCtx.Infof("Successfully updated the live application spec")

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
func cloneRepository(appName string, repoURL string, authCreds transport.AuthMethod, tempRoot string) (*git.Repository, error) {
	logCtx := log.WithContext().AddField("application", appName)
	logCtx.Infof("Cloning git repository %s in temporal folder located in %s", repoURL, tempRoot)
	r, err := git.PlainClone(tempRoot, false, &git.CloneOptions{
		Auth:     authCreds,
		URL:      repoURL,
		Progress: os.Stdout,
	})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// commitGitChanges commit the changes in the necessary file/s to the working copy
// of git repository
func commitGitChanges(appName string, gitW git.Worktree, commitMessage string, gitUsername string, gitEmail string) (*plumbing.Hash, error) {
	logCtx := log.WithContext().AddField("application", appName)
	// We can verify the current status of the worktree using the method Status.
	logCtx.Debugf("Obtaining current status after changes")
	status, err := gitW.Status()
	if err != nil {
		return nil, err
	}
	logCtx.Debugf("Obtained git status status is: %s", status)

	logCtx.Infof("It's going to commit changes with message: %s", commitMessage)
	commit, err := gitW.Commit(commitMessage, &git.CommitOptions{
		Author: &object.Signature{
			Name:  gitUsername,
			Email: gitEmail,
			When:  time.Now(),
		},
	})
	if err != nil {
		return nil, err
	}
	return &commit, nil
}

// pushGitChanges push the changes to the remote repository
func pushGitChanges(appName string, objC object.Commit, gitR *git.Repository, gitAuth transport.AuthMethod) error {
	logCtx := log.WithContext().AddField("application", appName)
	logCtx.Infof("It's going to push commit with hash %s and message %s", objC.Hash, objC.Message)

	logCtx.Infof("Pushing changes")
	err := gitR.Push(&git.PushOptions{
		Auth: gitAuth,
	})
	if err != nil {
		return err
	}

	logCtx.Infof("Successfully pushed changes")

	return nil
}

// commitAndPushGitChanges perfoms a git commit for the given pathSpec to the currently checked
// out branch and after pushes local changes to the remote branch
func commitAndPushGitChanges(cfg HelmUpdaterConfig, commitMessage string, gitW git.Worktree, tempRoot string, gitAuth transport.AuthMethod) error {
	logCtx := log.WithContext().AddField("application", cfg.AppName)

	targetFile := path.Join(cfg.GitConf.File, cfg.File)
	logCtx.Infof("Adding file %s to git for commit changes", targetFile)
	_, err := gitW.Add(targetFile)
	if err != nil {
		return err
	}

	commit, err := commitGitChanges(cfg.AppName, gitW, commitMessage, cfg.GitCredentials.Username, cfg.GitCredentials.Email)
	if err != nil {
		return err
	}

	gitR, err := git.PlainOpen(tempRoot)
	if err != nil {
		return err
	}

	logCtx.Debugf("Obtaining current HEAD to verify added changes")
	obj, err := gitR.CommitObject(*commit)
	if err != nil {
		return err
	}
	err = pushGitChanges(cfg.AppName, *obj, gitR, gitAuth)
	if err != nil {
		return err
	}

	return nil
}

// configureCommitMessage configure the git commit message
func configureCommitMessage(appName string, apps []ChangeEntry, helmUpdaterConfigMessage *template.Template) (*string, error) {
	var gitCommitMessage string

	logCtx := log.WithContext().AddField("application", appName)

	if len(apps) > 0 && helmUpdaterConfigMessage != nil {
		gitCommitMessage = TemplateCommitMessage(helmUpdaterConfigMessage, appName, apps)
		logCtx.Debugf("templated commit message successfully with value: %s", gitCommitMessage)
	}

	if gitCommitMessage == "" {
		tpl, err := template.New("commitMessage").Parse(git_internal.DefaultGitCommitMessage)
		if err != nil {
			return nil, fmt.Errorf("could not parse commit message template: %v", err)
		}
		gitCommitMessage = TemplateCommitMessage(tpl, appName, apps)
		logCtx.Debugf("templated commit message successfully with value: %s", gitCommitMessage)
	}
	return &gitCommitMessage, nil
}

// createTempFileInDirectory creates a temporal directory where a copy of
// the git repository is going to be stored.
func createTempFileInDirectory(dirName string, applicationName string, repoURL string) (*string, error) {
	logCtx := log.WithContext().AddField("application", applicationName)
	tempRoot, err := ioutil.TempDir(os.TempDir(), dirName)
	if err != nil {
		return nil, err
	}
	logCtx.Debugf("Created temporal directory %s to clone repository %s", tempRoot, repoURL)
	defer func() {
		err := os.RemoveAll(tempRoot)
		if err != nil {
			logCtx.Errorf("could not remove temp dir: %v", err)
		}
	}()
	return &tempRoot, nil
}

// getCheckoutBranchName obtain the name of the branch to be used
func getCheckoutBranchName(gitConfBranch string, applicationName string, gitR git.Repository) (*plumbing.ReferenceName, error) {
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

// checkBranchExists check if a specific branch in a repository was already created in the origin
func checkBranchExists(gitW git.Worktree, gitR git.Repository, checkOutBranchName plumbing.ReferenceName) (*git.Worktree, error) {
	err := gitW.Checkout(&git.CheckoutOptions{
		Branch: checkOutBranchName,
	})
	if err != nil {
		return nil, err
	}
	_, err = gitR.ResolveRevision(plumbing.Revision(checkOutBranchName))
	if err != nil {
		return nil, err
	}
	return &gitW, nil
}

// getRepositoryWorktreeWithBranchUpdated obtain working tree of git repositoy and checks if an specific
// branch exists already and pull latest changes
func getRepositoryWorktreeWithBranchUpdated(gitConfBranch string, appName string, gitR git.Repository, creds transport.AuthMethod) (*git.Worktree, error) {
	logCtx := log.WithContext().AddField("application", appName)
	gitW, err := gitR.Worktree()
	if err != nil {
		return nil, err
	}
	checkOutBranchName, err := getCheckoutBranchName(gitConfBranch, appName, gitR)
	if err != nil {
		return nil, err
	}

	gitWUpdated, err := checkBranchExists(*gitW, gitR, *checkOutBranchName)
	if err != nil {
		return nil, err
	}
	// Pull the latest changes from the origin remote and merge into the current branch
	logCtx.Infof("Pulling latest changes of branch %s", checkOutBranchName.Short())
	err = gitW.Pull(&git.PullOptions{
		Auth:          creds,
		ReferenceName: *checkOutBranchName,
	})

	if err != nil {
		if err.Error() != "already up-to-date" {
			return nil, err
		}
	}
	return gitWUpdated, nil
}

// fetchLatestChangesGitRepository fetch the latest changes in a git repository
func fetchLatestChangesGitRepository(appName string, gitR git.Repository, creds transport.AuthMethod) (*git.Repository, error) {
	logCtx := log.WithContext().AddField("application", appName)
	logCtx.Debugf("Fetching latest changes of repository")

	err := gitR.Fetch(&git.FetchOptions{
		RefSpecs: []config.RefSpec{"refs/*:refs/*", "HEAD:refs/heads/HEAD"},
		Auth:     creds,
		Force:    true,
	})
	if err != nil {
		return nil, err
	}

	return &gitR, nil
}

// cloneGitRepositoryInBranch clone git repository with a specific branch checking if that branch exists already
func cloneGitRepositoryInBranch(appName string, repoURL string, creds transport.AuthMethod, tempRoot string, gitConfBranch string) (*git.Worktree, error) {
	gitR, err := cloneRepository(appName, repoURL, creds, tempRoot)
	if err != nil {
		return nil, err
	}

	gitR, err = fetchLatestChangesGitRepository(appName, *gitR, creds)
	if err != nil {
		return nil, err
	}

	gitW, err := getRepositoryWorktreeWithBranchUpdated(gitConfBranch, appName, *gitR, creds)
	if err != nil {
		return nil, err
	}

	return gitW, nil
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

	tempRoot, err := createTempFileInDirectory(fmt.Sprintf("git-%s", cfg.AppName), cfg.AppName, cfg.GitConf.RepoURL)
	if err != nil {
		return nil, err
	}

	gitW, err := cloneGitRepositoryInBranch(cfg.AppName, cfg.GitConf.RepoURL, creds, *tempRoot, cfg.GitConf.Branch)
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

	err = commitAndPushGitChanges(cfg, *commitMessage, *gitW, *tempRoot, creds)
	if err != nil {
		return nil, err
	}

	return &apps, nil
}

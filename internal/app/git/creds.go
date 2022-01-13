package git

import (
	"fmt"

	"github.com/argoproj-labs/argocd-image-updater/ext/git"
)

// Credentials is a git credential config
type Credentials struct {
	Username   string
	Password   string
	Email      string
	SSHPrivKey string
}

// NewGitCreds returns the credentials for the given repo url.
func (c Credentials) NewGitCreds(repoURL string) (git.Creds, error) {
	if isSshUrl(repoURL) {
		return c.fromSsh(repoURL)
	}

	if isHttpsUrl(repoURL) {
		return c.fromHttps(repoURL)
	}

	return nil, unknownRepositoryType(repoURL)
}

func isSshUrl(repoUrl string) bool {
	ok, _ := git.IsSSHURL(repoUrl)

	return ok
}

func isHttpsUrl(repoUrl string) bool {
	return git.IsHTTPSURL(repoUrl)
}

func (c Credentials) fromSsh(repoUrl string) (git.Creds, error) {
	if c.allowsSshAuth() {
		return git.NewSSHCreds(c.SSHPrivKey, "", true), nil
	}

	return nil, sshPrivateKeyNotProvided(repoUrl)
}

func (c Credentials) fromHttps(repoURL string) (git.Creds, error) {
	if c.allowsHttpsAuth() {
		return git.NewHTTPSCreds(c.Username, c.Password, "", "", true, ""), nil
	}

	return nil, httpsUserAndPasswordNotProvided(repoURL)
}

func (c Credentials) allowsSshAuth() bool {
	return c.SSHPrivKey != ""
}

func (c Credentials) allowsHttpsAuth() bool {
	return c.Username != "" && c.Password != ""
}

func sshPrivateKeyNotProvided(repoUrl string) error {
	return fmt.Errorf("sshPrivKey not provided for authenticatication to repository %s", repoUrl)
}

func httpsUserAndPasswordNotProvided(repoUrl string) error {
	return fmt.Errorf("no value provided for username and password for authentication to repository %s", repoUrl)
}

func unknownRepositoryType(repoUrl string) error {
	return fmt.Errorf("unknown repository type for git repository URL %s", repoUrl)
}

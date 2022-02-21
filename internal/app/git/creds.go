package git

import (
	"fmt"
	"regexp"

	"github.com/docplanner/helm-repo-updater/internal/app/log"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

var (
	sshURLRegex   = regexp.MustCompile("^(ssh://)?([^/:]*?)@[^@]+$")
	httpsURLRegex = regexp.MustCompile("^(https://).*")
)

// Credentials is a git credential config
type Credentials struct {
	Username   string
	Password   string
	Email      string
	SSHPrivKey string
}

// NewGitCloneOpts returns the options that are going to be used to clone the repository.
func (c Credentials) NewGitCreds(repoURL string, password string) (transport.AuthMethod, error) {
	if isSSHURL(repoURL) {
		gitSSHCredentials, err := c.fromSsh(repoURL, password)
		if err != nil {
			return nil, err
		}
		return gitSSHCredentials, nil
	}

	if isHTTPSURL(repoURL) {
		gitCreds, err := c.fromHttps(repoURL)
		if err != nil {
			return nil, err
		}
		return gitCreds, nil
	}

	return nil, unknownRepositoryType(repoURL)
}

// IsSSHURL returns true if supplied URL is SSH URL
func isSSHURL(url string) bool {
	matches := sshURLRegex.FindStringSubmatch(url)
	return len(matches) > 2
}

// IsHTTPSURL returns true if supplied URL is HTTPS URL
func isHTTPSURL(url string) bool {
	return httpsURLRegex.MatchString(url)
}

func generateAuthForSSH(repoUrl string, userName string, privateKeyFile string, password string) (ssh.AuthMethod, error) {
	publicKeys, err := ssh.NewPublicKeysFromFile("git", privateKeyFile, password)
	if err != nil {
		log.Warnf("generate publickeys failed: %s\n", err.Error())
		return nil, err
	}
	return publicKeys, err
}

func generatAuthForHttps(username string, password string) *http.BasicAuth {
	return &http.BasicAuth{
		Username: username,
		Password: password,
	}
}

func (c Credentials) fromSsh(repoUrl string, password string) (ssh.AuthMethod, error) {
	if c.allowsSshAuth() {
		sshPublicKeys, err := generateAuthForSSH(repoUrl, c.Username, c.SSHPrivKey, password)
		if err != nil {
			return nil, err
		}
		return sshPublicKeys, nil
	}

	return nil, sshPrivateKeyNotProvided(repoUrl)
}

func (c Credentials) fromHttps(repoURL string) (*http.BasicAuth, error) {
	if c.allowsHttpsAuth() {
		return generatAuthForHttps(c.Username, c.Password), nil
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

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

// NewGitCreds returns credentials for use with go-git library
func (c Credentials) NewGitCreds(repoURL string, password string) (transport.AuthMethod, error) {
	if isSSHURL(repoURL) {
		gitSSHCredentials, err := c.fromSSH(repoURL, password)
		if err != nil {
			return nil, err
		}
		return gitSSHCredentials, nil
	}

	if isHTTPSURL(repoURL) {
		gitCreds, err := c.from(repoURL)
		if err != nil {
			return nil, err
		}
		return gitCreds, nil
	}

	return nil, unknownRepositoryType(repoURL)
}

// isSSHURL returns true if supplied URL is SSH URL
func isSSHURL(url string) bool {
	matches := sshURLRegex.FindStringSubmatch(url)
	return len(matches) > 2
}

// isHTTPSURL returns true if supplied URL is a valid HTTPS URL
func isHTTPSURL(url string) bool {
	return httpsURLRegex.MatchString(url)
}

// generateAuthForSSH generate the necessary public keys as auth for git repository using
// the provided privateKeyFile containing a valid SSH private key
func generateAuthForSSH(repoUrl string, userName string, privateKeyFile string, password string) (ssh.AuthMethod, error) {
	publicKeys, err := ssh.NewPublicKeysFromFile("git", privateKeyFile, password)
	if err != nil {
		log.Warnf("generate publickeys failed: %s\n", err.Error())
		return nil, err
	}
	return publicKeys, err
}

// fromSSH generate a valid credentials using ssh key
func (c Credentials) fromSSH(repoUrl string, password string) (ssh.AuthMethod, error) {
	if c.allowsSshAuth() {
		sshPublicKeys, err := generateAuthForSSH(repoUrl, c.Username, c.SSHPrivKey, password)
		if err != nil {
			return nil, err
		}
		return sshPublicKeys, nil
	}

	return nil, sshPrivateKeyNotProvided(repoUrl)
}

// generatAuthFor generate a valid credentials for go-git library using
// username and password
func generatAuthFor(username string, password string) *http.BasicAuth {
	return &http.BasicAuth{
		Username: username,
		Password: password,
	}
}

// from generate a valid credentials for go-git library using
// username and passowrd
func (c Credentials) from(repoURL string) (*http.BasicAuth, error) {
	if c.allowsAuth() {
		return generatAuthFor(c.Username, c.Password), nil
	}

	return nil, UserAndPasswordNotProvided(repoURL)
}

// allowSshAuth check if necessary attributes for generate an SSH
// credentials are provided
func (c Credentials) allowsSshAuth() bool {
	return c.SSHPrivKey != ""
}

// allowsAuth check if necessary attributes for generate and
// credentials are provided
func (c Credentials) allowsAuth() bool {
	return c.Username != "" && c.Password != ""
}

// sshPrivateKeyNotProvided return an error used when sshPrivKey
// is not provided for generate and SSH credentials
func sshPrivateKeyNotProvided(repoUrl string) error {
	return fmt.Errorf("sshPrivKey not provided for authenticatication to repository %s", repoUrl)
}

// UserAndPasswordNotProvided return an error used when
// username or password are not provided for generate and  credentials
func UserAndPasswordNotProvided(repoUrl string) error {
	return fmt.Errorf("no value provided for username and password for authentication to repository %s", repoUrl)
}

// unknownRepositoryType return an error used when
// the repository provided is not  or SSH type
func unknownRepositoryType(repoUrl string) error {
	return fmt.Errorf("unknown repository type for git repository URL %s", repoUrl)
}

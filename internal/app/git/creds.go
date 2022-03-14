package git

import (
	"fmt"
	"os"
	"regexp"

	"github.com/docplanner/helm-repo-updater/internal/app/log"
	app_utils "github.com/docplanner/helm-repo-updater/internal/app/utils"
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
	Username             string
	Password             string
	Email                string
	SSHPrivKey           string
	SSHPrivKeyFileInline bool
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
func generateAuthForSSH(repoURL string, userName string, privateKeyFile string, SSHPrivKeyFileInline bool, password string) (ssh.AuthMethod, error) {
	sshPrivKeyFileName := privateKeyFile
	if SSHPrivKeyFileInline {
		sshPrivKeyFile, err := app_utils.CreateAndWriteContentInTempFile("sshPrivKey", privateKeyFile)
		if err != nil {
			return nil, err
		}
		sshPrivKeyFileName = sshPrivKeyFile.Name()
		// close and remove the temporary file at the end of the program
		defer sshPrivKeyFile.Close()
		defer os.Remove(sshPrivKeyFileName)
		log.Infof("Generated file in %s location with content of SSH private key provided as input", sshPrivKeyFileName)
	}
	publicKeys, err := ssh.NewPublicKeysFromFile("git", sshPrivKeyFileName, password)
	if err != nil {
		log.Warnf("generate publickeys failed: %s\n", err.Error())
		return nil, err
	}
	return publicKeys, err
}

// fromSSH generate a valid credentials using ssh key
func (c Credentials) fromSSH(repoURL string, password string) (ssh.AuthMethod, error) {
	if c.allowsSSHAuth() {
		sshPublicKeys, err := generateAuthForSSH(repoURL, c.Username, c.SSHPrivKey, c.SSHPrivKeyFileInline, password)
		if err != nil {
			return nil, err
		}
		return sshPublicKeys, nil
	}

	return nil, sshPrivateKeyNotProvided(repoURL)
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
func (c Credentials) allowsSSHAuth() bool {
	return c.SSHPrivKey != ""
}

// allowsAuth check if necessary attributes for generate and
// credentials are provided
func (c Credentials) allowsAuth() bool {
	return c.Username != "" && c.Password != ""
}

// sshPrivateKeyNotProvided return an error used when sshPrivKey
// is not provided for generate and SSH credentials
func sshPrivateKeyNotProvided(repoURL string) error {
	return fmt.Errorf("sshPrivKey not provided for authenticatication to repository %s", repoURL)
}

// UserAndPasswordNotProvided return an error used when
// username or password are not provided for generate and  credentials
func UserAndPasswordNotProvided(repoURL string) error {
	return fmt.Errorf("no value provided for username and password for authentication to repository %s", repoURL)
}

// unknownRepositoryType return an error used when
// the repository provided is not  or SSH type
func unknownRepositoryType(repoURL string) error {
	return fmt.Errorf("unknown repository type for git repository URL %s", repoURL)
}

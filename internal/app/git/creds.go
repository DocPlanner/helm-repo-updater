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

// NewCreds returns the credentials for the given repo url.
func (g Credentials) NewCreds(repoURL string) (git.Creds, error) {
	if ok, _ := git.IsSSHURL(repoURL); ok {
		if g.SSHPrivKey != "" {
			return git.NewSSHCreds(g.SSHPrivKey, "", true), nil
		} else {
			return nil, fmt.Errorf(
				"sshPrivKey not provided for authenticatication to repository %s",
				repoURL,
			)
		}
	} else if git.IsHTTPSURL(repoURL) {
		if g.Username != "" && g.Password != "" {
			return git.NewHTTPSCreds(g.Username, g.Password, "", "", true, ""), nil
		} else {
			return nil, fmt.Errorf(
				"no value provided for username and password for authentication to repository %s",
				repoURL,
			)
		}
	}

	return nil, fmt.Errorf("unknown repository type for git repository URL %s", repoURL)
}

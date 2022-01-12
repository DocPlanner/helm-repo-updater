package git

import (
	"fmt"
	"log"
	"testing"

	"github.com/argoproj-labs/argocd-image-updater/ext/git"
	app_utils "github.com/docplanner/helm-repo-updater/internal/app/utils"
	"github.com/google/go-cmp/cmp"
	"gotest.tools/v3/assert"
)

const validGitCredentialsEmail = "test-user@docplanner.com"
const validGitCredentialsUsername = "test-user"
const validGitCredentialsPassword = "test-password"
const validSSHPrivKeyRelativeRoute = "/test-git-server/private_keys/helm-repo-updater-test"
const validGitRepoSSHURL = "git@github.com:kubernetes/kubernetes.git"
const validGitRepoHTTPSURL = "https://github.com/kubernetes/kubernetes.git"
const invalidGitRepoURL = "github.com/kubernetes/kubernetes.git"

func TestNewCredsSSHURLSShPrivKey(t *testing.T) {

	sshPrivKeyRoute, err := app_utils.GetRouteRelativePath(2, validSSHPrivKeyRelativeRoute)
	if err != nil {
		log.Fatal(err)
	}

	g := Credentials{
		Email:      validGitCredentialsEmail,
		SSHPrivKey: *sshPrivKeyRoute,
	}

	repoURL := validGitRepoSSHURL

	creds, err := g.NewGitCreds(repoURL)

	if err != nil {
		log.Fatal(err)
	}

	expectedCreds := git.NewSSHCreds(*sshPrivKeyRoute, "", true)

	assert.DeepEqual(t, creds, expectedCreds, cmp.AllowUnexported(git.SSHCreds{}))
}

func TestNewCredsHTPPSURLUsernamePassword(t *testing.T) {

	g := Credentials{
		Email:    validGitCredentialsEmail,
		Username: validGitCredentialsUsername,
		Password: validGitCredentialsPassword,
	}

	creds, err := g.NewGitCreds(validGitRepoHTTPSURL)

	if err != nil {
		log.Fatal(err)
	}

	expectedCreds := git.NewHTTPSCreds(g.Username, g.Password, "", "", true, "")

	assert.DeepEqual(t, creds, expectedCreds, cmp.AllowUnexported(git.HTTPSCreds{}))
}

func TestNewCredsSSHURLWithoutSShPrivKey(t *testing.T) {

	g := Credentials{
		Email:      validGitCredentialsEmail,
		SSHPrivKey: "",
	}

	repoURL := validGitRepoSSHURL

	_, err := g.NewGitCreds(repoURL)

	expectedErrorMessage := fmt.Sprintf(
		"sshPrivKey not provided for authenticatication to repository %s",
		repoURL,
	)

	assert.Error(t, err, expectedErrorMessage)
}

func TestNewCredsHTPPSURLWithoutUsernameWithPassword(t *testing.T) {

	g := Credentials{
		Email:    validGitCredentialsEmail,
		Username: "",
		Password: validGitCredentialsPassword,
	}

	repoURL := validGitRepoHTTPSURL

	_, err := g.NewGitCreds(repoURL)

	expectedErrorMessage := fmt.Sprintf(
		"no value provided for username and password for authentication to repository %s",
		repoURL,
	)

	assert.Error(t, err, expectedErrorMessage)
}

func TestNewCredsHTPPSURLWitUsernameWithoutPassword(t *testing.T) {

	g := Credentials{
		Email:    validGitCredentialsEmail,
		Username: validGitCredentialsUsername,
		Password: "",
	}

	repoURL := validGitRepoHTTPSURL

	_, err := g.NewGitCreds(repoURL)

	expectedErrorMessage := fmt.Sprintf(
		"no value provided for username and password for authentication to repository %s",
		repoURL,
	)

	assert.Error(t, err, expectedErrorMessage)
}

func TestNewCredsInvalidURL(t *testing.T) {

	g := Credentials{
		Email:    validGitCredentialsEmail,
		Username: validGitCredentialsUsername,
		Password: validGitCredentialsPassword,
	}

	repoURL := invalidGitRepoURL

	_, err := g.NewGitCreds(repoURL)

	expectedErrorMessage := fmt.Sprintf(
		"unknown repository type for git repository URL %s",
		repoURL,
	)

	assert.Error(t, err, expectedErrorMessage)
}

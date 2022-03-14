package git

import (
	"fmt"
	"log"
	"testing"

	app_utils "github.com/docplanner/helm-repo-updater/internal/app/utils"
	"gotest.tools/assert"
)

const (
	validGitCredentialsEmail     = "test-user@docplanner.com"
	validGitCredentialsUsername  = "test-user"
	validGitCredentialsPassword  = "test-password"
	validSSHPrivKeyRelativeRoute = "/test-git-server/private_keys/helm-repo-updater-test"
	validSSHPrivKeyString        = `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
QyNTUxOQAAACCNAv45QMrXuGuWk7uadYNOlL1B2q/g2pw1g+xP5oD/EQAAAJjmelMZ5npT
GQAAAAtzc2gtZWQyNTUxOQAAACCNAv45QMrXuGuWk7uadYNOlL1B2q/g2pw1g+xP5oD/EQ
AAAEA8ubLuVW5jc+Q9a2divLLVfm0Up+eus/9f7/HvACUD2o0C/jlAyte4a5aTu5p1g06U
vUHar+DanDWD7E/mgP8RAAAAFWRldm9wc0Bkb2NwbGFubmVyLmNvbQ==
-----END OPENSSH PRIVATE KEY-----`
	validGitRepoSSHURL   = "git@github.com:kubernetes/kubernetes.git"
	validGitRepoHTTPSURL = "https://github.com/kubernetes/kubernetes.git"
	invalidGitRepoURL    = "github.com/kubernetes/kubernetes.git"
	invalidPrivKeyRoute  = "/tmp/key-dont-exists"
)

func TestNewCredsSSHURLSSHPrivKey(t *testing.T) {

	sshPrivKeyRoute, err := app_utils.GetRouteRelativePath(2, validSSHPrivKeyRelativeRoute)
	if err != nil {
		log.Fatal(err)
	}

	g := Credentials{
		Username:   validGitCredentialsUsername,
		Email:      validGitCredentialsEmail,
		SSHPrivKey: *sshPrivKeyRoute,
	}

	repoURL := validGitRepoSSHURL

	creds, err := g.NewGitCreds(repoURL, g.Password)
	if err != nil {
		log.Fatal(err)
	}

	expectedCredsString := "user: git, name: ssh-public-keys"
	assert.DeepEqual(t, creds.String(), expectedCredsString)
}

func TestNewCredsSSHURLSSHPrivKeyFromString(t *testing.T) {
	g := Credentials{
		Username:             validGitCredentialsUsername,
		Email:                validGitCredentialsEmail,
		SSHPrivKey:           validSSHPrivKeyString,
		SSHPrivKeyFileInline: true,
	}

	repoURL := validGitRepoSSHURL

	creds, err := g.NewGitCreds(repoURL, g.Password)
	if err != nil {
		log.Fatal(err)
	}

	expectedCredsString := "user: git, name: ssh-public-keys"
	assert.DeepEqual(t, creds.String(), expectedCredsString)
}

func TestNewCredsHTPPSURLUsernamePassword(t *testing.T) {

	g := Credentials{
		Email:    validGitCredentialsEmail,
		Username: validGitCredentialsUsername,
		Password: validGitCredentialsPassword,
	}

	creds, err := g.NewGitCreds(validGitRepoHTTPSURL, g.Password)

	if err != nil {
		log.Fatal(err)
	}

	expectedCredsString := fmt.Sprintf("http-basic-auth - %s:*******", g.Username)
	assert.DeepEqual(t, creds.String(), expectedCredsString)
}

func TestNewCredsSSHURLSSHErroCreatingPublicKeys(t *testing.T) {

	g := Credentials{
		Username:   validGitCredentialsUsername,
		Email:      validGitCredentialsEmail,
		SSHPrivKey: invalidPrivKeyRoute,
	}

	repoURL := validGitRepoSSHURL

	_, err := g.NewGitCreds(repoURL, g.Password)

	expectedErrorMessage := fmt.Sprintf(
		"open %s: no such file or directory",
		invalidPrivKeyRoute,
	)

	assert.Error(t, err, expectedErrorMessage)
}

func TestNewCredsSSHURLWithoutSShPrivKey(t *testing.T) {

	g := Credentials{
		Email:      validGitCredentialsEmail,
		SSHPrivKey: "",
	}

	repoURL := validGitRepoSSHURL

	_, err := g.NewGitCreds(repoURL, g.Password)

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

	_, err := g.NewGitCreds(repoURL, g.Password)

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

	_, err := g.NewGitCreds(repoURL, g.Password)

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

	_, err := g.NewGitCreds(repoURL, g.Password)

	expectedErrorMessage := fmt.Sprintf(
		"unknown repository type for git repository URL %s",
		repoURL,
	)

	assert.Error(t, err, expectedErrorMessage)
}

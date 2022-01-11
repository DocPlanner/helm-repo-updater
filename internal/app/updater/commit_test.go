package updater

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/template"

	"github.com/docplanner/helm-repo-updater/internal/app/git"
	"gotest.tools/assert"
)

const validGitCredentialsEmail = "test-user@docplanner.com"
const validGitCredentialsUsername = "test-user"
const validGitRepoHost = "ssh://git@localhost:2222"
const validGitRepoRoute = "/git-server/repos/test-repo.git"
const validGitRepoURL = validGitRepoHost + validGitRepoRoute
const invalidGitRepoRoute = "/git-server/repos/test-r"
const invalidGitRepoURL = validGitRepoHost + invalidGitRepoRoute
const validSSHPrivKeyRelativeRoute = "/test-git-server/private_keys/helm-repo-updater-test"
const validGitRepoBranch = "develop"
const invalidGitRepoBranch = "developp"
const validHelmAppName = "example-app"
const validHelmAppFileToChange = validHelmAppName + "/values.yaml"

func getRouteRelativePath(numRelativePath int, relativePath string) (*string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	parent := filepath.Dir(wd)
	s := strings.Split(parent, "/")
	s = s[:len(s)-numRelativePath]
	finalPath := strings.Join(s, "/")
	finalPath = finalPath + relativePath
	return &finalPath, nil
}

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestUpdateApplicationDryRunNoChanges(t *testing.T) {

	sshPrivKeyRoute, err := getRouteRelativePath(2, validSSHPrivKeyRelativeRoute)
	if err != nil {
		log.Fatal(err)
	}

	gCred := git.Credentials{
		Email:      validGitCredentialsEmail,
		Username:   validGitCredentialsUsername,
		SSHPrivKey: *sshPrivKeyRoute,
	}

	gConf := git.Conf{
		RepoURL: validGitRepoURL,
		Branch:  validGitRepoBranch,
		File:    "",
	}

	changeEntry := ChangeEntry{
		OldValue: "1.0.0",
		NewValue: "1.0.0",
		File:     "",
		Key:      ".image.tag",
	}
	changeEntries := []ChangeEntry{
		changeEntry,
	}

	cfg := HelmUpdaterConfig{
		DryRun:         true,
		LogLevel:       "info",
		AppName:        validHelmAppName,
		UpdateApps:     changeEntries,
		File:           validHelmAppFileToChange,
		GitCredentials: &gCred,
		GitConf:        &gConf,
	}

	syncState := NewSyncIterationState()
	apps, err := UpdateApplication(cfg, syncState)

	if err != nil {
		log.Fatal(err)
	}

	assert.DeepEqual(t, *apps, changeEntries)
}

func TestUpdateApplicationDryRun(t *testing.T) {

	sshPrivKeyRoute, err := getRouteRelativePath(2, validSSHPrivKeyRelativeRoute)
	if err != nil {
		log.Fatal(err)
	}

	gCred := git.Credentials{
		Email:      validGitCredentialsEmail,
		Username:   validGitCredentialsUsername,
		SSHPrivKey: *sshPrivKeyRoute,
	}

	gConf := git.Conf{
		RepoURL: validGitRepoURL,
		Branch:  validGitRepoBranch,
		File:    "",
	}

	changeEntry := ChangeEntry{
		OldValue: "1.0.0",
		NewValue: "1.1.0",
		File:     "",
		Key:      ".image.tag",
	}
	changeEntries := []ChangeEntry{
		changeEntry,
	}

	cfg := HelmUpdaterConfig{
		DryRun:         true,
		LogLevel:       "info",
		AppName:        validHelmAppName,
		UpdateApps:     changeEntries,
		File:           validHelmAppFileToChange,
		GitCredentials: &gCred,
		GitConf:        &gConf,
	}

	syncState := NewSyncIterationState()
	apps, err := UpdateApplication(cfg, syncState)

	if err != nil {
		log.Fatal(err)
	}

	assert.DeepEqual(t, *apps, changeEntries)
}

func TestUpdateApplicationDryRunNoRepoURL(t *testing.T) {

	sshPrivKeyRoute, err := getRouteRelativePath(2, validSSHPrivKeyRelativeRoute)
	if err != nil {
		log.Fatal(err)
	}

	gCred := git.Credentials{
		Email:      validGitCredentialsEmail,
		Username:   validGitCredentialsUsername,
		SSHPrivKey: *sshPrivKeyRoute,
	}

	gConf := git.Conf{
		RepoURL: "",
		Branch:  validGitRepoBranch,
		File:    "",
	}

	changeEntry := ChangeEntry{
		OldValue: "1.0.0",
		NewValue: "1.1.0",
		File:     "",
		Key:      ".image.tag",
	}
	changeEntries := []ChangeEntry{
		changeEntry,
	}

	cfg := HelmUpdaterConfig{
		DryRun:         true,
		LogLevel:       "info",
		AppName:        validHelmAppName,
		UpdateApps:     changeEntries,
		File:           validHelmAppFileToChange,
		GitCredentials: &gCred,
		GitConf:        &gConf,
	}

	syncState := NewSyncIterationState()
	_, err = UpdateApplication(cfg, syncState)
	errorMessage := fmt.Sprintf("could not get creds for repo '%s': unknown repository type for git repository URL", cfg.AppName)

	assert.ErrorContains(t, err, errorMessage)
}

func TestUpdateApplicationDryRunInvalidGitRepo(t *testing.T) {

	sshPrivKeyRoute, err := getRouteRelativePath(2, validSSHPrivKeyRelativeRoute)
	if err != nil {
		log.Fatal(err)
	}

	gCred := git.Credentials{
		Email:      validGitCredentialsEmail,
		Username:   validGitCredentialsUsername,
		SSHPrivKey: *sshPrivKeyRoute,
	}

	gConf := git.Conf{
		RepoURL: invalidGitRepoURL,
		Branch:  validGitRepoBranch,
		File:    "",
	}

	changeEntry := ChangeEntry{
		OldValue: "1.0.0",
		NewValue: "1.1.0",
		File:     "",
		Key:      ".image.tag",
	}
	changeEntries := []ChangeEntry{
		changeEntry,
	}

	cfg := HelmUpdaterConfig{
		DryRun:         true,
		LogLevel:       "info",
		AppName:        validHelmAppName,
		UpdateApps:     changeEntries,
		File:           validHelmAppFileToChange,
		GitCredentials: &gCred,
		GitConf:        &gConf,
	}

	syncState := NewSyncIterationState()
	_, err = UpdateApplication(cfg, syncState)
	assert.ErrorContains(t, err, fmt.Sprintf("fatal: '%s' does not appear to be a git repository", invalidGitRepoRoute))
	assert.ErrorContains(t, err, "fatal: Could not read from remote repository.")
}

func TestUpdateApplicationDryRunInvalidGitBranch(t *testing.T) {

	sshPrivKeyRoute, err := getRouteRelativePath(2, validSSHPrivKeyRelativeRoute)
	if err != nil {
		log.Fatal(err)
	}

	gCred := git.Credentials{
		Email:      validGitCredentialsEmail,
		Username:   validGitCredentialsUsername,
		SSHPrivKey: *sshPrivKeyRoute,
	}

	gConf := git.Conf{
		RepoURL: validGitRepoURL,
		Branch:  invalidGitRepoBranch,
		File:    "",
	}

	changeEntry := ChangeEntry{
		OldValue: "1.0.0",
		NewValue: "1.1.0",
		File:     "",
		Key:      ".image.tag",
	}
	changeEntries := []ChangeEntry{
		changeEntry,
	}

	cfg := HelmUpdaterConfig{
		DryRun:         true,
		LogLevel:       "info",
		AppName:        validHelmAppName,
		UpdateApps:     changeEntries,
		File:           validHelmAppFileToChange,
		GitCredentials: &gCred,
		GitConf:        &gConf,
	}

	syncState := NewSyncIterationState()
	_, err = UpdateApplication(cfg, syncState)
	expectedErrorMessage := fmt.Sprintf("`git checkout --force %s` failed exit status 1: error: pathspec '%s' did not match any file(s) known to git", invalidGitRepoBranch, invalidGitRepoBranch)

	assert.Error(t, err, expectedErrorMessage)
}

func TestUpdateApplicationDryRuNoBranch(t *testing.T) {

	sshPrivKeyRoute, err := getRouteRelativePath(2, validSSHPrivKeyRelativeRoute)
	if err != nil {
		log.Fatal(err)
	}

	gCred := git.Credentials{
		Email:      validGitCredentialsEmail,
		Username:   validGitCredentialsUsername,
		SSHPrivKey: *sshPrivKeyRoute,
	}

	gConf := git.Conf{
		RepoURL: validGitRepoURL,
		Branch:  "",
		File:    "",
	}

	changeEntry := ChangeEntry{
		OldValue: "1.0.0",
		NewValue: "1.1.0",
		File:     "",
		Key:      ".image.tag",
	}
	changeEntries := []ChangeEntry{
		changeEntry,
	}

	cfg := HelmUpdaterConfig{
		DryRun:         true,
		LogLevel:       "info",
		AppName:        validHelmAppName,
		UpdateApps:     changeEntries,
		File:           validHelmAppFileToChange,
		GitCredentials: &gCred,
		GitConf:        &gConf,
	}

	syncState := NewSyncIterationState()
	_, err = UpdateApplication(cfg, syncState)
	expectedErrorMessage := "could not resolve symbolic ref '': `git symbolic-ref ` failed exit status 128: fatal: No such ref:"

	assert.Error(t, err, expectedErrorMessage)
}

func TestUpdateApplicationDryRunWithGitMessage(t *testing.T) {
	temp, err := template.New("commit-message").Parse("Simple change in app")

	if err != nil {
		log.Fatal(err)
	}

	sshPrivKeyRoute, err := getRouteRelativePath(2, validSSHPrivKeyRelativeRoute)
	if err != nil {
		log.Fatal(err)
	}

	gCred := git.Credentials{
		Email:      validGitCredentialsEmail,
		Username:   validGitCredentialsUsername,
		SSHPrivKey: *sshPrivKeyRoute,
	}

	gConf := git.Conf{
		RepoURL: validGitRepoURL,
		Branch:  validGitRepoBranch,
		File:    "",
		Message: temp,
	}

	changeEntry := ChangeEntry{
		OldValue: "1.0.0",
		NewValue: "1.1.0",
		File:     "",
		Key:      ".image.tag",
	}
	changeEntries := []ChangeEntry{
		changeEntry,
	}

	cfg := HelmUpdaterConfig{
		DryRun:         true,
		LogLevel:       "info",
		AppName:        validHelmAppName,
		UpdateApps:     changeEntries,
		File:           validHelmAppFileToChange,
		GitCredentials: &gCred,
		GitConf:        &gConf,
	}

	syncState := NewSyncIterationState()
	apps, err := UpdateApplication(cfg, syncState)

	if err != nil {
		log.Fatal(err)
	}

	assert.DeepEqual(t, *apps, changeEntries)
}

func TestUpdateApplicationDryRunInvalidKey(t *testing.T) {

	sshPrivKeyRoute, err := getRouteRelativePath(2, validSSHPrivKeyRelativeRoute)
	if err != nil {
		log.Fatal(err)
	}

	gCred := git.Credentials{
		Email:      validGitCredentialsEmail,
		Username:   validGitCredentialsUsername,
		SSHPrivKey: *sshPrivKeyRoute,
	}

	gConf := git.Conf{
		RepoURL: validGitRepoURL,
		Branch:  validGitRepoBranch,
		File:    "",
	}

	changeEntry1 := ChangeEntry{
		OldValue: "1.0.0",
		NewValue: "1.1.0",
		File:     "",
		Key:      ".image.tag",
	}
	changeEntry2 := ChangeEntry{
		OldValue: "1.0.0",
		NewValue: "1.1.0",
		File:     "",
		Key:      "image.tag",
	}
	changeEntries := []ChangeEntry{
		changeEntry1,
		changeEntry2,
	}
	expectedChangedEntries := []ChangeEntry{
		changeEntry1,
	}

	cfg := HelmUpdaterConfig{
		DryRun:         true,
		LogLevel:       "info",
		AppName:        validHelmAppName,
		UpdateApps:     changeEntries,
		File:           validHelmAppFileToChange,
		GitCredentials: &gCred,
		GitConf:        &gConf,
	}

	syncState := NewSyncIterationState()
	apps, err := UpdateApplication(cfg, syncState)

	if err != nil {
		log.Fatal(err)
	}

	assert.DeepEqual(t, *apps, expectedChangedEntries)
}

//TODO: Check check with git-repo-server created with docker
func TestUpdateApplication(t *testing.T) {

	sshPrivKeyRoute, err := getRouteRelativePath(2, validSSHPrivKeyRelativeRoute)
	if err != nil {
		log.Fatal(err)
	}

	gCred := git.Credentials{
		Email:      validGitCredentialsEmail,
		Username:   validGitCredentialsUsername,
		SSHPrivKey: *sshPrivKeyRoute,
	}

	gConf := git.Conf{
		RepoURL: validGitRepoURL,
		Branch:  validGitRepoBranch,
		File:    "",
	}

	changeEntry := ChangeEntry{
		OldValue: "1.0.0",
		NewValue: "1.1.0",
		File:     "",
		Key:      ".image.tag",
	}
	changeEntries := []ChangeEntry{
		changeEntry,
	}

	cfg := HelmUpdaterConfig{
		DryRun:         false,
		LogLevel:       "info",
		AppName:        validHelmAppName,
		UpdateApps:     changeEntries,
		File:           validHelmAppFileToChange,
		GitCredentials: &gCred,
		GitConf:        &gConf,
	}

	syncState := NewSyncIterationState()
	apps, err := UpdateApplication(cfg, syncState)
	if err != nil {
		log.Fatal(err)
	}
	assert.DeepEqual(t, *apps, changeEntries)
}

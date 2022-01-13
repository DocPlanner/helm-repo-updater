package updater

import (
	"fmt"
	"log"
	"os"
	"testing"
	"text/template"

	"github.com/docplanner/helm-repo-updater/internal/app/git"
	app_utils "github.com/docplanner/helm-repo-updater/internal/app/utils"
	"gotest.tools/v3/assert"
)

const validGitCredentialsEmail = "test-user@docplanner.com"
const validGitCredentialsUsername = "test-user"

// TODO: Remove after CI tests docker layer name
const validGitRepoHost = "ssh://git@git-server:22"
const validGitRepoRoute = "/git-server/repos/test-repo.git"
const validGitRepoURL = validGitRepoHost + validGitRepoRoute
const invalidGitRepoRoute = "/git-server/repos/test-r"
const invalidGitRepoURL = validGitRepoHost + invalidGitRepoRoute
const validSSHPrivKeyRelativeRoute = "/test-git-server/private_keys/helm-repo-updater-test"
const validGitRepoBranch = "develop"
const invalidGitRepoBranch = "developp"
const validHelmAppName = "example-app"
const validHelmAppFileToChange = validHelmAppName + "/values.yaml"

func loadSSHKeyPath(sshPrivKeyPath string) (string, error) {
	dat, err := os.ReadFile(sshPrivKeyPath)
	if err != nil {
		return "", err
	}

	return string(dat), nil
}

func TestUpdateApplicationDryRunNoChanges(t *testing.T) {

	sshPrivKeyRoute, err := app_utils.GetRouteRelativePath(2, validSSHPrivKeyRelativeRoute)
	if err != nil {
		log.Fatal(err)
	}

	sshPrivKeyData, err := loadSSHKeyPath(*sshPrivKeyRoute)
	if err != nil {
		log.Fatal(err)
	}

	gCred := git.Credentials{
		Email:      validGitCredentialsEmail,
		Username:   validGitCredentialsUsername,
		SSHPrivKey: sshPrivKeyData,
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

	sshPrivKeyRoute, err := app_utils.GetRouteRelativePath(2, validSSHPrivKeyRelativeRoute)
	if err != nil {
		log.Fatal(err)
	}

	sshPrivKeyData, err := loadSSHKeyPath(*sshPrivKeyRoute)
	if err != nil {
		log.Fatal(err)
	}

	gCred := git.Credentials{
		Email:      validGitCredentialsEmail,
		Username:   validGitCredentialsUsername,
		SSHPrivKey: sshPrivKeyData,
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

	sshPrivKeyRoute, err := app_utils.GetRouteRelativePath(2, validSSHPrivKeyRelativeRoute)
	if err != nil {
		log.Fatal(err)
	}

	sshPrivKeyData, err := loadSSHKeyPath(*sshPrivKeyRoute)
	if err != nil {
		log.Fatal(err)
	}

	gCred := git.Credentials{
		Email:      validGitCredentialsEmail,
		Username:   validGitCredentialsUsername,
		SSHPrivKey: sshPrivKeyData,
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

	sshPrivKeyRoute, err := app_utils.GetRouteRelativePath(2, validSSHPrivKeyRelativeRoute)
	if err != nil {
		log.Fatal(err)
	}

	sshPrivKeyData, err := loadSSHKeyPath(*sshPrivKeyRoute)
	if err != nil {
		log.Fatal(err)
	}

	gCred := git.Credentials{
		Email:      validGitCredentialsEmail,
		Username:   validGitCredentialsUsername,
		SSHPrivKey: sshPrivKeyData,
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

	sshPrivKeyRoute, err := app_utils.GetRouteRelativePath(2, validSSHPrivKeyRelativeRoute)
	if err != nil {
		log.Fatal(err)
	}

	sshPrivKeyData, err := loadSSHKeyPath(*sshPrivKeyRoute)
	if err != nil {
		log.Fatal(err)
	}

	gCred := git.Credentials{
		Email:      validGitCredentialsEmail,
		Username:   validGitCredentialsUsername,
		SSHPrivKey: sshPrivKeyData,
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

	sshPrivKeyRoute, err := app_utils.GetRouteRelativePath(2, validSSHPrivKeyRelativeRoute)
	if err != nil {
		log.Fatal(err)
	}

	sshPrivKeyData, err := loadSSHKeyPath(*sshPrivKeyRoute)
	if err != nil {
		log.Fatal(err)
	}

	gCred := git.Credentials{
		Email:      validGitCredentialsEmail,
		Username:   validGitCredentialsUsername,
		SSHPrivKey: sshPrivKeyData,
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

	sshPrivKeyRoute, err := app_utils.GetRouteRelativePath(2, validSSHPrivKeyRelativeRoute)
	if err != nil {
		log.Fatal(err)
	}

	sshPrivKeyData, err := loadSSHKeyPath(*sshPrivKeyRoute)
	if err != nil {
		log.Fatal(err)
	}

	gCred := git.Credentials{
		Email:      validGitCredentialsEmail,
		Username:   validGitCredentialsUsername,
		SSHPrivKey: sshPrivKeyData,
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

	sshPrivKeyRoute, err := app_utils.GetRouteRelativePath(2, validSSHPrivKeyRelativeRoute)
	if err != nil {
		log.Fatal(err)
	}

	sshPrivKeyData, err := loadSSHKeyPath(*sshPrivKeyRoute)
	if err != nil {
		log.Fatal(err)
	}

	gCred := git.Credentials{
		Email:      validGitCredentialsEmail,
		Username:   validGitCredentialsUsername,
		SSHPrivKey: sshPrivKeyData,
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

func TestUpdateApplication(t *testing.T) {

	sshPrivKeyRoute, err := app_utils.GetRouteRelativePath(2, validSSHPrivKeyRelativeRoute)
	if err != nil {
		log.Fatal(err)
	}

	sshPrivKeyData, err := loadSSHKeyPath(*sshPrivKeyRoute)
	if err != nil {
		log.Fatal(err)
	}

	gCred := git.Credentials{
		Email:      validGitCredentialsEmail,
		Username:   validGitCredentialsUsername,
		SSHPrivKey: sshPrivKeyData,
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

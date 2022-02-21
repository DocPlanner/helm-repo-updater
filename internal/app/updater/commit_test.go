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

const (
	validGitCredentialsEmail      = "test-user@docplanner.com"
	validGitCredentialsUsername   = "test-user"
	SSHRepoPrefix                 = "ssh://git@"
	SSHRepoLocalHostname          = SSHRepoPrefix + "localhost:2222"
	SSHRepoCIHostname             = SSHRepoPrefix + "git-server"
	validGitRepoRoute             = "/git-server/repos/test-repo.git"
	invalidGitRepoRoute           = "/git-server/repos/test-r"
	SSHPrivKeyRelativePath        = 2
	validSSHPrivKeyRelativeRoute  = "/test-git-server/private_keys/helm-repo-updater-test"
	validGitRepoBranch            = "develop"
	invalidGitRepoBranch          = "developp"
	validHelmAppName              = "example-app"
	validHelmAppFileToChange      = validHelmAppName + "/values.yaml"
	ciDiscoveryEnvironmentName    = "isCI"
	isDevContainerEnvironmentName = "isDevContainer"
)

func TestUpdateApplicationDryRunNoChangesEntries(t *testing.T) {

	sshPrivKeyRoute, err := app_utils.GetRouteRelativePath(2, validSSHPrivKeyRelativeRoute)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	gCred := git.Credentials{
		Email:      validGitCredentialsEmail,
		Username:   validGitCredentialsUsername,
		SSHPrivKey: *sshPrivKeyRoute,
	}

	validGitRepoURL := getSSHRepoHostnameAndPort() + validGitRepoRoute

	gConf := git.Conf{
		RepoURL: validGitRepoURL,
		Branch:  validGitRepoBranch,
		File:    "",
	}

	changeEntries := []ChangeEntry{}

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

	expectedErrorMessage := "nothing to update, skipping commit"

	assert.Error(t, err, expectedErrorMessage)
}

func TestUpdateApplicationDryRunNoChanges(t *testing.T) {

	sshPrivKeyRoute, err := app_utils.GetRouteRelativePath(2, validSSHPrivKeyRelativeRoute)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	gCred := git.Credentials{
		Email:      validGitCredentialsEmail,
		Username:   validGitCredentialsUsername,
		SSHPrivKey: *sshPrivKeyRoute,
	}

	validGitRepoURL := getSSHRepoHostnameAndPort() + validGitRepoRoute

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

	if err != nil {
		log.Fatal(err)
	}

	gCred := git.Credentials{
		Email:      validGitCredentialsEmail,
		Username:   validGitCredentialsUsername,
		SSHPrivKey: *sshPrivKeyRoute,
	}

	validGitRepoURL := getSSHRepoHostnameAndPort() + validGitRepoRoute

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

	sshPrivKeyRoute, err := app_utils.GetRouteRelativePath(2, validSSHPrivKeyRelativeRoute)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	gCred := git.Credentials{
		Email:      validGitCredentialsEmail,
		Username:   validGitCredentialsUsername,
		SSHPrivKey: *sshPrivKeyRoute,
	}

	invalidGitRepoURL := getSSHRepoHostnameAndPort() + invalidGitRepoRoute

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
	expectErrorMessage := "repository not found"
	assert.Error(t, err, expectErrorMessage)
}

func TestUpdateApplicationDryRunInvalidGitBranch(t *testing.T) {

	sshPrivKeyRoute, err := app_utils.GetRouteRelativePath(2, validSSHPrivKeyRelativeRoute)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	gCred := git.Credentials{
		Email:      validGitCredentialsEmail,
		Username:   validGitCredentialsUsername,
		SSHPrivKey: *sshPrivKeyRoute,
	}

	validGitRepoURL := getSSHRepoHostnameAndPort() + validGitRepoRoute

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

	expectedErrorMessage := "reference not found"
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

	if err != nil {
		log.Fatal(err)
	}

	gCred := git.Credentials{
		Email:      validGitCredentialsEmail,
		Username:   validGitCredentialsUsername,
		SSHPrivKey: *sshPrivKeyRoute,
	}

	validGitRepoURL := getSSHRepoHostnameAndPort() + validGitRepoRoute

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

	if err != nil {
		log.Fatal(err)
	}

	gCred := git.Credentials{
		Email:      validGitCredentialsEmail,
		Username:   validGitCredentialsUsername,
		SSHPrivKey: *sshPrivKeyRoute,
	}

	validGitRepoURL := getSSHRepoHostnameAndPort() + validGitRepoRoute

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

	if err != nil {
		log.Fatal(err)
	}

	gCred := git.Credentials{
		Email:      validGitCredentialsEmail,
		Username:   validGitCredentialsUsername,
		SSHPrivKey: *sshPrivKeyRoute,
	}

	validGitRepoURL := getSSHRepoHostnameAndPort() + validGitRepoRoute

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

func getSSHRepoHostnameAndPort() string {
	_, isCI := os.LookupEnv(ciDiscoveryEnvironmentName)
	_, isDevContainerEnvironment := os.LookupEnv(isDevContainerEnvironmentName)
	if !isCI && !isDevContainerEnvironment {
		return SSHRepoLocalHostname
	}
	return SSHRepoCIHostname
}

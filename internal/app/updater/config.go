package updater

import (
	"github.com/docplanner/helm-repo-updater/internal/app/git"
)

// HelmUpdaterConfig contains global configuration and required runtime data
type HelmUpdaterConfig struct {
	DryRun         bool
	LogLevel       string
	AppName        string
	UpdateApps     []ChangeEntry
	File           string
	GitCredentials *git.Credentials
	GitConf        *git.Conf
}

// ChangeEntry represents values that has been changed by Helm Updater
type ChangeEntry struct {
	OldValue string
	NewValue string
	File     string
	Key      string
}

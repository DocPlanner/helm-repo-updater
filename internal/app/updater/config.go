package updater

import (
	"github.com/docplanner/helm-repo-updater/internal/app/git"
	"github.com/docplanner/helm-repo-updater/internal/app/logger"
)

// HelmUpdaterConfig contains global configuration and required runtime data
type HelmUpdaterConfig struct {
	DryRun         bool
	AppName        string
	UpdateApps     []ChangeEntry
	File           string
	GitCredentials *git.Credentials
	GitConf        *git.Conf
	Logger         logger.Logger
}

// ChangeEntry represents values that have been changed by Helm Updater
type ChangeEntry struct {
	OldValue string
	NewValue string
	File     string
	Key      string
}

package updater

import (
	"fmt"
	"github.com/docplanner/helm-repo-updater/internal/app/logger"
	"os"
	"path"

	"github.com/argoproj-labs/argocd-image-updater/ext/git"
	"github.com/docplanner/helm-repo-updater/internal/app/yq"
)

var _ changeWriter = writeOverrides

type changeWriter func(cfg HelmUpdaterConfig, gitC git.Client) (apps []ChangeEntry, err error)

//writeOverrides writes the overrides to the git files
func writeOverrides(cfg HelmUpdaterConfig, gitC git.Client) (apps []ChangeEntry, err error) {
	targetFile := path.Join(gitC.Root(), cfg.GitConf.File, cfg.File)

	apps = make([]ChangeEntry, 0)

	_, err = os.Stat(targetFile)
	if err != nil {
		cfg.Logger.ErrorWithContext("target file doesn't exist", err, logger.LogContext{
			"application": cfg.AppName,
			"file":        cfg.File,
		})

		return apps, err
	}

	apps = overrideValues(apps, cfg, targetFile)

	if len(apps) == 0 {
		return apps, fmt.Errorf("nothing to update, skipping commit")
	}

	return apps, gitC.Add(targetFile)
}

// overrideValues overrides values in the given file
func overrideValues(apps []ChangeEntry, cfg HelmUpdaterConfig, targetFile string) []ChangeEntry {
	var err error

	for _, app := range cfg.UpdateApps {
		// define new entry
		var newEntry ChangeEntry
		var oldValue, newValue *string

		// replace helm parameters
		oldValue, err = yq.ReadKey(app.Key, targetFile)
		if err != nil {
			cfg.Logger.WarningWithContext("can not read the presented key due to error, skipping change", logger.LogContext{
				"application": cfg.AppName,
				"key":         app.Key,
				"error":       err.Error(),
			})

			continue
		}

		newEntry.Key = app.Key
		newEntry.OldValue = *oldValue

		// replace helm parameters
		cfg.Logger.DebugWithContext("settings new value", logger.LogContext{
			"application": cfg.AppName,
			"key":         app.Key,
			"value":       app.NewValue,
		})

		err = yq.InplaceApply(app.Key, app.NewValue, targetFile)
		if err != nil {
			cfg.Logger.WarningWithContext("failed to update key", logger.LogContext{
				"application": cfg.AppName,
				"key":         app.Key,
				"error":       err.Error(),
			})

			newEntry.NewValue = *oldValue

			continue
		}

		// check patched app
		newValue, err = yq.ReadKey(app.Key, targetFile)
		if err != nil {
			cfg.Logger.WarningWithContext("failed to read the patched key, skipping change", logger.LogContext{
				"application": cfg.AppName,
				"key":         app.Key,
				"error":       err.Error(),
			})
			newEntry.NewValue = *oldValue

			continue
		}
		newEntry.NewValue = *newValue

		// check if there is any change
		if oldValue == newValue {
			cfg.Logger.WarningWithContext("target for key is the same, skipping change", logger.LogContext{
				"application": cfg.AppName,
				"key":         app.Key,
			})

			continue
		}

		apps = append(apps, newEntry)
	}

	return apps
}

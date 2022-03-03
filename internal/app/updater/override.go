package updater

import (
	"fmt"
	"os"
	"path"

	"github.com/docplanner/helm-repo-updater/internal/app/log"
	"github.com/docplanner/helm-repo-updater/internal/app/yq"
	git "github.com/go-git/go-git/v5"
)

var _ changeWriter = writeOverrides

type changeWriter func(cfg HelmUpdaterConfig, tempRoot string, gitW git.Worktree) (apps []ChangeEntry, err error)

// writeOverrides writes the overrides to the git files
func writeOverrides(cfg HelmUpdaterConfig, tempRoot string, gitW git.Worktree) (apps []ChangeEntry, err error) {
	targetFile := path.Join(tempRoot, cfg.GitConf.File, cfg.File)

	apps = make([]ChangeEntry, 0)

	_, err = os.Stat(targetFile)
	if err != nil {
		log.WithContext().
			AddField("application", cfg.AppName).
			Errorf("target file %s doesn't exist.", cfg.File)

		return apps, err
	}

	apps = overrideValues(apps, cfg, targetFile)

	if len(apps) == 0 {
		return apps, fmt.Errorf("nothing to update, skipping commit")
	}

	return apps, nil
}

// overrideValues overrides values in the given file
func overrideValues(apps []ChangeEntry, cfg HelmUpdaterConfig, targetFile string) []ChangeEntry {
	var err error

	logCtx := log.WithContext().AddField("application", cfg.AppName)
	for _, app := range cfg.UpdateApps {
		// define new entry
		var newEntry ChangeEntry
		var oldValue, newValue *string

		// replace helm parameters
		oldValue, err = yq.ReadKey(app.Key, targetFile)
		if err != nil {
			logCtx.Infof("failed to read the presented key %s due to error %s, skipping change", app.Key, err.Error())

			continue
		}

		newEntry.Key = app.Key
		newEntry.OldValue = *oldValue

		// replace helm parameters
		logCtx.Infof("Actual value for key %s: %s", app.Key, newEntry.OldValue)
		logCtx.Infof("Setting new value for key %s: %s", app.Key, app.NewValue)
		err = yq.InplaceApply(app.Key, app.NewValue, targetFile)
		if err != nil {
			logCtx.Infof("failed to update key %s: %v", app.Key, err)

			newEntry.NewValue = *oldValue

			continue
		}

		// check patched app
		newValue, err = yq.ReadKey(app.Key, targetFile)
		if err != nil {
			logCtx.Infof("failed to read the patched key %s due to error %s, skipping change", app.Key, err.Error())
			newEntry.NewValue = *oldValue

			continue
		}
		newEntry.NewValue = *newValue

		// check if there is any change
		if oldValue == newValue {
			logCtx.Infof("target for key %s is the same, skipping", app.Key)

			continue
		}

		apps = append(apps, newEntry)
	}

	return apps
}

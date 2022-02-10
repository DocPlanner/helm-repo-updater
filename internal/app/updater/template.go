package updater

import (
	"bytes"
	"github.com/docplanner/helm-repo-updater/internal/app/logger"
	"text/template"
)

type commitMessageChange struct {
	File     string
	Key      string
	OldValue string
	NewValue string
}

type commitMessageTemplate struct {
	AppName    string
	KeyChanges []commitMessageChange
}

// TemplateCommitMessage renders a commit message template and returns it
// as a string. If the template could not be rendered, returns a default message.
func TemplateCommitMessage(logger logger.Logger, tpl *template.Template, appName string, changeList []ChangeEntry) string {
	var cmBuf bytes.Buffer
	changes := make([]commitMessageChange, 0)
	for _, c := range changeList {
		changes = append(changes, commitMessageChange{c.File, c.Key, c.OldValue, c.NewValue})
	}

	tplData := commitMessageTemplate{
		AppName:    appName,
		KeyChanges: changes,
	}

	err := tpl.Execute(&cmBuf, tplData)
	if err != nil {
		logger.ErrorWithContext("could not execute template for git commit message", err, map[string]interface{}{
			"application": appName,
		})

		return "build: update of application " + appName
	}

	return cmBuf.String()
}

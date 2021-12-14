package updater

import (
	"bytes"
	"text/template"

	"github.com/argoproj-labs/argocd-image-updater/pkg/log"
)

// TemplateCommitMessage renders a commit message template and returns it
// as a string. If the template could not be rendered, returns a default message.
func TemplateCommitMessage(tpl *template.Template, appName string, changeList []ChangeEntry) string {
	var cmBuf bytes.Buffer

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
		log.Errorf("could not execute template for Git commit message: %v", err)

		return "build: update of application " + appName
	}

	return cmBuf.String()
}

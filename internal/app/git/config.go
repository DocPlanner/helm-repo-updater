package git

import "text/template"

const DefaultGitCommitMessage = `🚀 automatic update of {{ .AppName }}

{{ range .KeyChanges -}}
updates key {{ .Key }} value from '{{ .OldValue }}' to '{{ .NewValue }}'
{{ end -}}
`

// Conf is the configuration for the git client
type Conf struct {
	RepoURL string
	Branch  string
	File    string
	Message *template.Template
}

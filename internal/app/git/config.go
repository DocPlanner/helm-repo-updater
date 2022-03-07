package git

import "text/template"

// DefaultGitCommitMessage is the default commit message build with the changes detected in the app
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

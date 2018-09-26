package cobra

var HELP_TEMPLATE = `{{ $cmd := . }}Usage: {{if .Runnable}}
  {{.UseLine}}{{if .HasFlags}}{{end}}{{end}}{{if .HasSubCommands}}
  {{ .CommandPath}} [command]{{end}}

Description:
  {{.Long | trim}}{{if gt .Aliases 0}}

Aliases:
  {{.NameAndAliases}}{{end}}
{{ if .HasSubCommands}}
Available Commands: {{range .Commands}}{{if and (.Runnable) (ne .Use "help [command]")}}
  {{rpad .Name .UsagePadding }} {{.Short}}{{end}}{{end}}
{{end}}{{ if .HasFlags}} 
Available Flags:
{{.Flags.FlagUsages}}{{end}}{{if .HasParent}}{{if and (gt .Commands 0) (gt .Parent.Commands 1) }}
Additional help topics: {{if gt .Commands 0 }}{{range .Commands}}{{if not .Runnable}} {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if gt .Parent.Commands 1 }}{{range .Parent.Commands}}{{if .Runnable}}{{if not (eq .Name $cmd.Name) }}{{end}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{end}}
{{end}}{{ if .HasSubCommands }}
Use "{{.Root.Name}} help [command]" for more information about that command.
{{end}}`

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

const (
	banner = `.____           __    __          
|    |    _____/  |__/  |_  ____  
|    |   /  _ \   __\   __\/  _ \ 
|    |__(  <_> )  |  |  | (  <_> )
|_______ \____/|__|  |__|  \____/ 
        \/    `

	appSummary = "Create beautiful drawings using a simple scripting language."
)

// rootCmd represents the base command when called without any subcommands
var (
	rootCmd = &cobra.Command{
		DisableSuggestions:    true,
		DisableFlagsInUseLine: true,
		Use:                   fmt.Sprintf("%s <COMMAND>", appName()),
		Short:                 appSummary,
		Long:                  banner,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	rootCmd.Version = version
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	rootCmd.SetUsageTemplate(`Utilizzo:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Esempi:
  {{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Commandi disponibili:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Parametri:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Parametri Globali:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Digita "{{.CommandPath}} [command] --help" per maggiori informazioni su un comando.{{end}}
`)

	rootCmd.SetVersionTemplate(`{{with .Name}}{{printf "%s " .}}{{end}}{{printf "%s" .Version}} - Luca Sepe <luca.sepe@gmail.com>
`)
}

func appName() string {
	return filepath.Base(os.Args[0])
}

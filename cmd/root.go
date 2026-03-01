package cmd

import (
	"fmt"
	"os"

	"github.com/maolei1024/gopate/pkg/i18n"
	"github.com/spf13/cobra"
)

var (
	verbose bool
	quiet   bool
	lang    string
)

var rootCmd = &cobra.Command{
	Use:   "gopate",
	Short: i18n.T("root.short"),
	Long:  i18n.T("root.long"),
}

// Execute executes the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, i18n.T("flag.verbose"))
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, i18n.T("flag.quiet"))
	rootCmd.PersistentFlags().StringVar(&lang, "lang", "", i18n.T("flag.lang"))

	cobra.OnInitialize(initLang)
}

// initLang handles the --lang flag by switching the language
// and updating all command descriptions before help is displayed.
func initLang() {
	if lang != "" {
		i18n.SetLanguage(lang)
		updateCommandTexts()
	}
}

// updateCommandTexts refreshes all cobra command text with current i18n language.
func updateCommandTexts() {
	// Root
	rootCmd.Short = i18n.T("root.short")
	rootCmd.Long = i18n.T("root.long")

	// Disguise
	disguiseCmd.Use = i18n.T("disguise.use")
	disguiseCmd.Short = i18n.T("disguise.short")
	disguiseCmd.Long = i18n.T("disguise.long")

	// Reveal
	revealCmd.Use = i18n.T("reveal.use")
	revealCmd.Short = i18n.T("reveal.short")
	revealCmd.Long = i18n.T("reveal.long")

	// Inspect
	inspectCmd.Use = i18n.T("inspect.use")
	inspectCmd.Short = i18n.T("inspect.short")
	inspectCmd.Long = i18n.T("inspect.long")

	// Version
	versionCmd.Short = i18n.T("version.short")
}

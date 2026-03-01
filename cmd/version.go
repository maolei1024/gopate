package cmd

import (
	"fmt"

	"github.com/maolei1024/gopate/pkg/i18n"
	"github.com/spf13/cobra"
)

// Version is the version string, can be injected via ldflags at build time
var Version = "1.1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: i18n.T("version.short"),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Gopate v%s\n", Version)
		fmt.Println(i18n.T("msg.compatible"))
		fmt.Println(i18n.T("msg.homepage"))
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

package cmd

import (
	"fmt"
	"os"

	"github.com/maolei1024/gopate/pkg/apate"
	"github.com/maolei1024/gopate/pkg/i18n"
	"github.com/spf13/cobra"
)

var inspectCmd = &cobra.Command{
	Use:   i18n.T("inspect.use"),
	Short: i18n.T("inspect.short"),
	Long:  i18n.T("inspect.long"),
	Args:  cobra.MinimumNArgs(1),
	RunE:  runInspect,
}

func init() {
	rootCmd.AddCommand(inspectCmd)
}

func runInspect(cmd *cobra.Command, args []string) error {
	files, err := collectFiles(args, false)
	if err != nil {
		return err
	}

	for _, filePath := range files {
		result, err := apate.Inspect(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, i18n.Tf("msg.inspect_failed", filePath, err)+"\n")
			continue
		}

		if result.IsDisguised {
			fmt.Printf(i18n.Tf("msg.disguised_yes", filePath) + "\n")
			fmt.Printf(i18n.Tf("msg.file_size", formatSize(result.FileSize)) + "\n")
			fmt.Printf(i18n.Tf("msg.mask_length", result.MaskLength) + "\n")
			fmt.Printf(i18n.Tf("msg.disguise_type", result.DetectedType) + "\n")
			if verbose && len(result.OriginalHeader) > 0 {
				fmt.Printf(i18n.T("msg.original_header"))
				displayLen := len(result.OriginalHeader)
				if displayLen > 32 {
					displayLen = 32
				}
				for j := 0; j < displayLen; j++ {
					fmt.Printf("%02X ", result.OriginalHeader[j])
				}
				if len(result.OriginalHeader) > 32 {
					fmt.Printf("...")
				}
				fmt.Println()
			}
		} else {
			fmt.Printf(i18n.Tf("msg.disguised_no", filePath) + "\n")
			fmt.Printf(i18n.Tf("msg.file_size", formatSize(result.FileSize)) + "\n")
		}
	}

	return nil
}

func formatSize(size int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)
	switch {
	case size >= GB:
		return fmt.Sprintf("%.2f GB", float64(size)/float64(GB))
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/float64(MB))
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/float64(KB))
	default:
		return fmt.Sprintf("%d B", size)
	}
}

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/maolei1024/gopate/pkg/apate"
	"github.com/maolei1024/gopate/pkg/i18n"
	"github.com/spf13/cobra"
)

var (
	revealOutputDir string
	revealInPlace   bool
	revealRecursive bool
	revealDryRun    bool
	revealForce     bool
)

var revealCmd = &cobra.Command{
	Use:   i18n.T("reveal.use"),
	Short: i18n.T("reveal.short"),
	Long:  i18n.T("reveal.long"),
	Args:  cobra.MinimumNArgs(1),
	RunE:  runReveal,
}

func init() {
	revealCmd.Flags().StringVarP(&revealOutputDir, "output-dir", "o", "", i18n.T("flag.output_dir"))
	revealCmd.Flags().BoolVar(&revealInPlace, "in-place", false, i18n.T("flag.in_place"))
	revealCmd.Flags().BoolVarP(&revealRecursive, "recursive", "r", false, i18n.T("flag.recursive"))
	revealCmd.Flags().BoolVar(&revealDryRun, "dry-run", false, i18n.T("flag.dry_run"))
	revealCmd.Flags().BoolVarP(&revealForce, "force", "f", false, i18n.T("flag.force"))
	rootCmd.AddCommand(revealCmd)
}

func runReveal(cmd *cobra.Command, args []string) error {
	// Collect all files
	files, err := collectFiles(args, revealRecursive)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return fmt.Errorf(i18n.T("msg.no_files"))
	}

	if revealDryRun {
		fmt.Printf(i18n.Tf("msg.dry_run_reveal", len(files)) + "\n")
		for _, f := range files {
			fmt.Printf("  %s\n", f)
		}
		return nil
	}

	// Safety warning
	if !revealForce && !quiet {
		fmt.Println(i18n.T("msg.reveal_warning"))
		fmt.Printf(i18n.Tf("msg.confirm_proceed", len(files)))
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "y" && confirm != "Y" {
			fmt.Println(i18n.T("msg.cancelled"))
			return nil
		}
	}

	successCount := 0
	failCount := 0

	for _, filePath := range files {
		if verbose {
			fmt.Printf(i18n.Tf("msg.revealing", filePath) + "\n")
		}

		if revealInPlace {
			// In-place mode
			if err := apate.Reveal(filePath); err != nil {
				failCount++
				if !quiet {
					fmt.Fprintf(os.Stderr, i18n.Tf("msg.failed", filePath, err)+"\n")
				}
				continue
			}
			newPath, err := apate.RenameRevealed(filePath)
			if err != nil {
				failCount++
				if !quiet {
					fmt.Fprintf(os.Stderr, i18n.Tf("msg.rename_failed", filePath, err)+"\n")
				}
				continue
			}
			successCount++
			if verbose {
				fmt.Printf("  → %s\n", newPath)
			}
		} else {
			// New file mode
			outputDir := filepath.Dir(filePath)
			if revealOutputDir != "" {
				outputDir = revealOutputDir
				if err := os.MkdirAll(outputDir, 0755); err != nil {
					return fmt.Errorf(i18n.T("msg.create_outdir_failed"), err)
				}
			}
			// Restored file name: remove last extension
			baseName := filepath.Base(filePath)
			ext := filepath.Ext(baseName)
			restoredName := baseName
			if ext != "" {
				restoredName = baseName[:len(baseName)-len(ext)]
			}
			dstPath := filepath.Join(outputDir, restoredName)
			if err := apate.RevealToFile(filePath, dstPath); err != nil {
				failCount++
				if !quiet {
					fmt.Fprintf(os.Stderr, i18n.Tf("msg.failed", filePath, err)+"\n")
				}
				continue
			}
			successCount++
			if verbose {
				fmt.Printf("  → %s\n", dstPath)
			}
		}
	}

	if !quiet {
		fmt.Printf(i18n.Tf("msg.done", successCount, failCount) + "\n")
	}

	if failCount > 0 {
		return fmt.Errorf(i18n.Tf("msg.some_failed", failCount))
	}
	return nil
}

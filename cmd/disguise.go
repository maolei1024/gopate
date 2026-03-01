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
	disguiseMode      string
	disguiseMaskFile  string
	disguiseOutputDir string
	disguiseInPlace   bool
	disguiseRecursive bool
	disguiseDryRun    bool
)

var disguiseCmd = &cobra.Command{
	Use:   i18n.T("disguise.use"),
	Short: i18n.T("disguise.short"),
	Long:  i18n.T("disguise.long"),
	Args:  cobra.MinimumNArgs(1),
	RunE:  runDisguise,
}

func init() {
	disguiseCmd.Flags().StringVarP(&disguiseMode, "mode", "m", "onekey", i18n.T("flag.mode"))
	disguiseCmd.Flags().StringVar(&disguiseMaskFile, "mask-file", "", i18n.T("flag.mask_file"))
	disguiseCmd.Flags().StringVarP(&disguiseOutputDir, "output-dir", "o", "", i18n.T("flag.output_dir"))
	disguiseCmd.Flags().BoolVar(&disguiseInPlace, "in-place", false, i18n.T("flag.in_place"))
	disguiseCmd.Flags().BoolVarP(&disguiseRecursive, "recursive", "r", false, i18n.T("flag.recursive"))
	disguiseCmd.Flags().BoolVar(&disguiseDryRun, "dry-run", false, i18n.T("flag.dry_run"))
	rootCmd.AddCommand(disguiseCmd)
}

func runDisguise(cmd *cobra.Command, args []string) error {
	// Get mask bytes and extension
	maskHead, maskExt, err := getMaskData(disguiseMode, disguiseMaskFile)
	if err != nil {
		return err
	}

	// Collect all files
	files, err := collectFiles(args, disguiseRecursive)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return fmt.Errorf(i18n.T("msg.no_files"))
	}

	if disguiseDryRun {
		fmt.Printf(i18n.Tf("msg.dry_run_disguise", len(files), disguiseMode) + "\n")
		for _, f := range files {
			fmt.Printf("  %s → %s%s\n", f, filepath.Base(f), maskExt)
		}
		return nil
	}

	successCount := 0
	failCount := 0

	for _, filePath := range files {
		if verbose {
			fmt.Printf(i18n.Tf("msg.disguising", filePath) + "\n")
		}

		if disguiseInPlace {
			// In-place mode
			if err := apate.Disguise(filePath, maskHead); err != nil {
				failCount++
				if !quiet {
					fmt.Fprintf(os.Stderr, i18n.Tf("msg.failed", filePath, err)+"\n")
				}
				continue
			}
			newPath, err := apate.RenameDisguised(filePath, maskExt)
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
			if disguiseOutputDir != "" {
				outputDir = disguiseOutputDir
				if err := os.MkdirAll(outputDir, 0755); err != nil {
					return fmt.Errorf(i18n.T("msg.create_outdir_failed"), err)
				}
			}
			dstPath := filepath.Join(outputDir, filepath.Base(filePath)+maskExt)
			if err := apate.DisguiseToFile(filePath, dstPath, maskHead); err != nil {
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

func getMaskData(mode, maskFile string) ([]byte, string, error) {
	switch mode {
	case "onekey":
		mask, err := apate.GetOnekeyMask()
		if err != nil {
			return nil, "", fmt.Errorf(i18n.T("msg.load_mask_failed"), err)
		}
		return mask, ".mp4", nil
	case "mask":
		if maskFile == "" {
			return nil, "", fmt.Errorf(i18n.T("msg.mask_file_required"))
		}
		mask, err := apate.FileToBytes(maskFile)
		if err != nil {
			return nil, "", fmt.Errorf(i18n.T("msg.read_mask_failed"), err)
		}
		ext := filepath.Ext(maskFile)
		if ext == "" {
			ext = ".bin"
		}
		return mask, ext, nil
	case "exe", "jpg", "mp4", "mov":
		head, ok := apate.ModeHead[mode]
		if !ok {
			return nil, "", fmt.Errorf(i18n.Tf("msg.unknown_mode", mode))
		}
		ext, ok := apate.ModeExtension[mode]
		if !ok {
			return nil, "", fmt.Errorf(i18n.Tf("msg.unknown_mode", mode))
		}
		return head, ext, nil
	default:
		return nil, "", fmt.Errorf(i18n.Tf("msg.unknown_mode_opts", mode))
	}
}

func collectFiles(paths []string, recursive bool) ([]string, error) {
	var allFiles []string
	for _, p := range paths {
		// Support glob patterns
		matches, err := filepath.Glob(p)
		if err != nil {
			return nil, fmt.Errorf(i18n.T("msg.glob_failed"), err)
		}
		if len(matches) == 0 {
			// Not a glob pattern, use raw path
			matches = []string{p}
		}
		for _, m := range matches {
			fi, err := os.Stat(m)
			if err != nil {
				return nil, fmt.Errorf(i18n.Tf("msg.access_failed", m, err))
			}
			if fi.IsDir() {
				if recursive {
					dirFiles, err := apate.GetAllFilesRecursively(m)
					if err != nil {
						return nil, fmt.Errorf(i18n.Tf("msg.walk_dir_failed", m, err))
					}
					allFiles = append(allFiles, dirFiles...)
				} else {
					fmt.Fprintf(os.Stderr, i18n.Tf("msg.skip_dir", m)+"\n")
				}
			} else {
				allFiles = append(allFiles, m)
			}
		}
	}
	return allFiles, nil
}

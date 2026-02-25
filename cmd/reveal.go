package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/maolei1024/gopate/pkg/apate"
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
	Use:   "reveal [文件或目录...]",
	Short: "还原伪装文件",
	Long: `将经过伪装的文件还原为原始格式。

注意: 如果对未经过伪装的文件执行还原操作，可能会导致文件损坏！
请务必做好数据备份。

示例:
  gopate reveal file.zip.mp4                # 还原单个文件
  gopate reveal *.mp4 --in-place            # 批量原地还原
  gopate reveal ./mydir -r --in-place       # 递归还原目录下所有文件
  gopate reveal file.mp4 -o ./restored/     # 还原到指定目录`,
	Args: cobra.MinimumNArgs(1),
	RunE: runReveal,
}

func init() {
	revealCmd.Flags().StringVarP(&revealOutputDir, "output-dir", "o", "", "输出目录（默认输出到源文件同目录）")
	revealCmd.Flags().BoolVar(&revealInPlace, "in-place", false, "原地修改文件（默认生成新文件）")
	revealCmd.Flags().BoolVarP(&revealRecursive, "recursive", "r", false, "递归处理目录")
	revealCmd.Flags().BoolVar(&revealDryRun, "dry-run", false, "预览模式，不实际修改文件")
	revealCmd.Flags().BoolVarP(&revealForce, "force", "f", false, "跳过确认提示")
	rootCmd.AddCommand(revealCmd)
}

func runReveal(cmd *cobra.Command, args []string) error {
	// 收集所有文件
	files, err := collectFiles(args, revealRecursive)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return fmt.Errorf("未找到任何文件")
	}

	if revealDryRun {
		fmt.Printf("预览模式 - 将还原 %d 个文件\n", len(files))
		for _, f := range files {
			fmt.Printf("  %s\n", f)
		}
		return nil
	}

	// 安全提示
	if !revealForce && !quiet {
		fmt.Println("⚠ 警告: 对未经伪装的文件执行还原可能导致文件损坏！请确保已备份。")
		fmt.Printf("即将处理 %d 个文件，是否继续? (y/N): ", len(files))
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "y" && confirm != "Y" {
			fmt.Println("已取消")
			return nil
		}
	}

	successCount := 0
	failCount := 0

	for _, filePath := range files {
		if verbose {
			fmt.Printf("还原: %s\n", filePath)
		}

		if revealInPlace {
			// 原地修改模式
			if err := apate.Reveal(filePath); err != nil {
				failCount++
				if !quiet {
					fmt.Fprintf(os.Stderr, "失败: %s - %v\n", filePath, err)
				}
				continue
			}
			newPath, err := apate.RenameRevealed(filePath)
			if err != nil {
				failCount++
				if !quiet {
					fmt.Fprintf(os.Stderr, "重命名失败: %s - %v\n", filePath, err)
				}
				continue
			}
			successCount++
			if verbose {
				fmt.Printf("  → %s\n", newPath)
			}
		} else {
			// 生成新文件模式
			outputDir := filepath.Dir(filePath)
			if revealOutputDir != "" {
				outputDir = revealOutputDir
				if err := os.MkdirAll(outputDir, 0755); err != nil {
					return fmt.Errorf("创建输出目录失败: %w", err)
				}
			}
			// 还原后文件名：去掉最后一个扩展名
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
					fmt.Fprintf(os.Stderr, "失败: %s - %v\n", filePath, err)
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
		fmt.Printf("完成！成功 %d 个，失败 %d 个\n", successCount, failCount)
	}

	if failCount > 0 {
		return fmt.Errorf("有 %d 个文件处理失败", failCount)
	}
	return nil
}

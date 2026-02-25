package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/maolei1024/gopate/pkg/apate"
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
	Use:   "disguise [文件或目录...]",
	Short: "伪装文件",
	Long: `对指定文件进行格式伪装。

伪装模式:
  onekey  使用内嵌 MP4 面具文件伪装（默认，适用大部分场景）
  mask    使用自定义面具文件伪装（需配合 --mask-file 参数）
  exe     使用 EXE 文件头伪装
  jpg     使用 JPG 文件头伪装
  mp4     使用 MP4 文件头伪装
  mov     使用 MOV 文件头伪装

示例:
  gopate disguise file.zip                         # 一键伪装为 MP4
  gopate disguise file.zip --mode exe              # 伪装为 EXE
  gopate disguise file.zip --mode mask --mask-file cover.png  # 使用自定义面具
  gopate disguise *.zip --mode mp4 --in-place      # 批量原地伪装
  gopate disguise ./mydir -r --mode onekey         # 递归伪装目录下所有文件`,
	Args: cobra.MinimumNArgs(1),
	RunE: runDisguise,
}

func init() {
	disguiseCmd.Flags().StringVarP(&disguiseMode, "mode", "m", "onekey", "伪装模式: onekey|mask|exe|jpg|mp4|mov")
	disguiseCmd.Flags().StringVar(&disguiseMaskFile, "mask-file", "", "自定义面具文件路径（仅 mask 模式需要）")
	disguiseCmd.Flags().StringVarP(&disguiseOutputDir, "output-dir", "o", "", "输出目录（默认输出到源文件同目录）")
	disguiseCmd.Flags().BoolVar(&disguiseInPlace, "in-place", false, "原地修改文件（默认生成新文件）")
	disguiseCmd.Flags().BoolVarP(&disguiseRecursive, "recursive", "r", false, "递归处理目录")
	disguiseCmd.Flags().BoolVar(&disguiseDryRun, "dry-run", false, "预览模式，不实际修改文件")
	rootCmd.AddCommand(disguiseCmd)
}

func runDisguise(cmd *cobra.Command, args []string) error {
	// 获取面具字节和扩展名
	maskHead, maskExt, err := getMaskData(disguiseMode, disguiseMaskFile)
	if err != nil {
		return err
	}

	// 收集所有文件
	files, err := collectFiles(args, disguiseRecursive)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return fmt.Errorf("未找到任何文件")
	}

	if disguiseDryRun {
		fmt.Printf("预览模式 - 将伪装 %d 个文件 (模式: %s)\n", len(files), disguiseMode)
		for _, f := range files {
			fmt.Printf("  %s → %s%s\n", f, filepath.Base(f), maskExt)
		}
		return nil
	}

	successCount := 0
	failCount := 0

	for _, filePath := range files {
		if verbose {
			fmt.Printf("伪装: %s\n", filePath)
		}

		if disguiseInPlace {
			// 原地修改模式
			if err := apate.Disguise(filePath, maskHead); err != nil {
				failCount++
				if !quiet {
					fmt.Fprintf(os.Stderr, "失败: %s - %v\n", filePath, err)
				}
				continue
			}
			newPath, err := apate.RenameDisguised(filePath, maskExt)
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
			if disguiseOutputDir != "" {
				outputDir = disguiseOutputDir
				if err := os.MkdirAll(outputDir, 0755); err != nil {
					return fmt.Errorf("创建输出目录失败: %w", err)
				}
			}
			dstPath := filepath.Join(outputDir, filepath.Base(filePath)+maskExt)
			if err := apate.DisguiseToFile(filePath, dstPath, maskHead); err != nil {
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

func getMaskData(mode, maskFile string) ([]byte, string, error) {
	switch mode {
	case "onekey":
		mask, err := apate.GetOnekeyMask()
		if err != nil {
			return nil, "", fmt.Errorf("加载内嵌面具失败: %w", err)
		}
		return mask, ".mp4", nil
	case "mask":
		if maskFile == "" {
			return nil, "", fmt.Errorf("mask 模式需要指定 --mask-file 参数")
		}
		mask, err := apate.FileToBytes(maskFile)
		if err != nil {
			return nil, "", fmt.Errorf("读取面具文件失败: %w", err)
		}
		ext := filepath.Ext(maskFile)
		if ext == "" {
			ext = ".bin"
		}
		return mask, ext, nil
	case "exe", "jpg", "mp4", "mov":
		head, ok := apate.ModeHead[mode]
		if !ok {
			return nil, "", fmt.Errorf("未知的伪装模式: %s", mode)
		}
		ext, ok := apate.ModeExtension[mode]
		if !ok {
			return nil, "", fmt.Errorf("未知的伪装模式: %s", mode)
		}
		return head, ext, nil
	default:
		return nil, "", fmt.Errorf("未知的伪装模式: %s\n可选: onekey, mask, exe, jpg, mp4, mov", mode)
	}
}

func collectFiles(paths []string, recursive bool) ([]string, error) {
	var allFiles []string
	for _, p := range paths {
		// 支持 glob 通配符
		matches, err := filepath.Glob(p)
		if err != nil {
			return nil, fmt.Errorf("路径匹配失败: %w", err)
		}
		if len(matches) == 0 {
			// 非 glob 模式，直接使用原路径
			matches = []string{p}
		}
		for _, m := range matches {
			fi, err := os.Stat(m)
			if err != nil {
				return nil, fmt.Errorf("无法访问: %s - %w", m, err)
			}
			if fi.IsDir() {
				if recursive {
					dirFiles, err := apate.GetAllFilesRecursively(m)
					if err != nil {
						return nil, fmt.Errorf("遍历目录失败: %s - %w", m, err)
					}
					allFiles = append(allFiles, dirFiles...)
				} else {
					fmt.Fprintf(os.Stderr, "跳过目录: %s (使用 -r 递归处理)\n", m)
				}
			} else {
				allFiles = append(allFiles, m)
			}
		}
	}
	return allFiles, nil
}

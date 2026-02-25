package cmd

import (
	"fmt"
	"os"

	"github.com/maolei1024/gopate/pkg/apate"
	"github.com/spf13/cobra"
)

var inspectCmd = &cobra.Command{
	Use:   "inspect [文件...]",
	Short: "检测文件是否经过 apate/gopate 伪装",
	Long: `分析文件是否经过 apate/gopate 伪装，并显示相关信息。

示例:
  gopate inspect file.mp4              # 检测单个文件
  gopate inspect *.mp4                 # 批量检测
  gopate inspect file.mp4 -v           # 显示详细信息（含原始文件头）`,
	Args: cobra.MinimumNArgs(1),
	RunE: runInspect,
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
			fmt.Fprintf(os.Stderr, "检测失败: %s - %v\n", filePath, err)
			continue
		}

		if result.IsDisguised {
			fmt.Printf("✅ %s - 已伪装\n", filePath)
			fmt.Printf("   文件大小: %s\n", formatSize(result.FileSize))
			fmt.Printf("   面具长度: %d 字节\n", result.MaskLength)
			fmt.Printf("   伪装类型: %s\n", result.DetectedType)
			if verbose && len(result.OriginalHeader) > 0 {
				fmt.Printf("   原始文件头: ")
				displayLen := len(result.OriginalHeader)
				if displayLen > 32 {
					displayLen = 32
				}
				for i := 0; i < displayLen; i++ {
					fmt.Printf("%02X ", result.OriginalHeader[i])
				}
				if len(result.OriginalHeader) > 32 {
					fmt.Printf("...")
				}
				fmt.Println()
			}
		} else {
			fmt.Printf("❌ %s - 未检测到伪装\n", filePath)
			fmt.Printf("   文件大小: %s\n", formatSize(result.FileSize))
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

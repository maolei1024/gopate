package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	verbose bool
	quiet   bool
)

var rootCmd = &cobra.Command{
	Use:   "gopate",
	Short: "Gopate - 文件格式伪装工具 (Go CLI)",
	Long: `Gopate 是一款基于 Go 语言的纯命令行文件格式伪装工具。
完全兼容 apate (C# 版本) 的文件格式，支持跨平台使用。

功能特性:
  • 支持超大文件，瞬间伪装/还原
  • 支持批量处理、递归目录
  • 原始文件头经过加密处理，不易被检测
  • 完全兼容 apate 伪装的文件，可混合使用`,
}

// Execute 执行根命令
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "显示详细输出")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "静默模式，仅显示错误")
}

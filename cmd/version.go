package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version 版本号，构建时可通过 ldflags 注入
var Version = "1.0.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示版本信息",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Gopate v%s\n", Version)
		fmt.Println("兼容 apate v1.4.2 文件格式")
		fmt.Println("项目主页: https://github.com/maolei1024/gopate")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

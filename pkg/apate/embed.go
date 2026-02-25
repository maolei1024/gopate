package apate

import (
	"embed"
)

//go:embed mask.mp4
var embedFS embed.FS

// GetOnekeyMask 获取内嵌的一键伪装面具文件（MP4）
func GetOnekeyMask() ([]byte, error) {
	return embedFS.ReadFile("mask.mp4")
}

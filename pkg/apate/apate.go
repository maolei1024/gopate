// Package apate 提供文件格式伪装/还原的核心引擎。
// 二进制格式与 C# 版本 apate (https://github.com/rippod/apate) 完全兼容。
package apate

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrFileTooSmall = errors.New("file too small to reveal")
	ErrInvalidMask  = errors.New("invalid mask data (length is 0)")
	ErrMaskTooLarge = errors.New("mask file too large")
	ErrNotDisguised = errors.New("file is not disguised by apate, or is corrupted")
	ErrFileNotExist = errors.New("file does not exist")
	ErrIsDirectory  = errors.New("target is a directory, not a file")
)

// InspectResult 包含对伪装文件的分析结果
type InspectResult struct {
	FilePath       string // 文件路径
	FileSize       int64  // 文件总大小
	IsDisguised    bool   // 是否经过 apate 伪装
	MaskLength     int    // 面具字节长度
	OriginalHeader []byte // 还原后的原始文件头（前几个字节）
	DetectedType   string // 根据面具头检测到的伪装类型
}

// Disguise 对文件进行伪装。
// 算法与 C# 版本完全一致：
//  1. 读取原始文件头 (长度 = min(len(maskHead), fileSize))
//  2. 将文件头替换为 maskHead
//  3. 在文件末尾追加：原始文件头的逆序 + 面具长度(4字节小端序)
func Disguise(filePath string, maskHead []byte) error {
	if len(maskHead) == 0 {
		return ErrInvalidMask
	}

	fi, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrFileNotExist
		}
		return fmt.Errorf("failed to get file info: %w", err)
	}
	if fi.IsDir() {
		return ErrIsDirectory
	}

	f, err := os.OpenFile(filePath, os.O_RDWR, 0)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	fileSize := fi.Size()

	// 读取原始文件头
	headLen := len(maskHead)
	if int64(headLen) > fileSize {
		headLen = int(fileSize)
	}
	originalHead := make([]byte, headLen)
	_, err = f.ReadAt(originalHead, 0)
	if err != nil {
		return fmt.Errorf("failed to read original file header: %w", err)
	}

	// 在文件头写入面具字节
	_, err = f.WriteAt(maskHead, 0)
	if err != nil {
		return fmt.Errorf("failed to write mask: %w", err)
	}

	// 准备追加数据：逆序原始文件头 + 面具长度标记
	reversedHead := reverseBytes(originalHead)
	maskLenBytes := make([]byte, MaskLengthIndicatorLength)
	binary.LittleEndian.PutUint32(maskLenBytes, uint32(len(maskHead)))

	// 追加到文件末尾
	appendData := append(reversedHead, maskLenBytes...)

	// 计算写入位置：如果面具比原文件大，文件已被扩展到 len(maskHead)
	writePos := fileSize
	if int64(len(maskHead)) > fileSize {
		writePos = int64(len(maskHead))
	}
	_, err = f.WriteAt(appendData, writePos)
	if err != nil {
		return fmt.Errorf("failed to append data to file: %w", err)
	}

	return nil
}

// DisguiseToFile 伪装文件并输出到新路径（不修改原文件）
func DisguiseToFile(srcPath, dstPath string, maskHead []byte) error {
	// 复制源文件到目标路径
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("failed to read source file: %w", err)
	}
	if err := os.WriteFile(dstPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write destination file: %w", err)
	}
	// 对目标文件进行伪装
	return Disguise(dstPath, maskHead)
}

// Reveal 还原伪装文件。
// 算法与 C# 版本完全一致：
//  1. 从文件最后4字节读取面具长度 (int32 小端序)
//  2. 从文件末尾倒数读取逆序的原始文件头
//  3. 截断文件末尾多余部分
//  4. 将逆序还原后的原始文件头写回文件开头
func Reveal(filePath string) error {
	fi, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrFileNotExist
		}
		return fmt.Errorf("failed to get file info: %w", err)
	}
	if fi.IsDir() {
		return ErrIsDirectory
	}

	fileSize := fi.Size()
	if fileSize < int64(MaskLengthIndicatorLength) {
		return ErrFileTooSmall
	}

	f, err := os.OpenFile(filePath, os.O_RDWR, 0)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	// 从文件末尾读取面具长度标记
	maskLenBytes := make([]byte, MaskLengthIndicatorLength)
	_, err = f.ReadAt(maskLenBytes, fileSize-int64(MaskLengthIndicatorLength))
	if err != nil {
		return fmt.Errorf("failed to read mask length indicator: %w", err)
	}
	maskHeadLength := int(binary.LittleEndian.Uint32(maskLenBytes))

	if maskHeadLength <= 0 || int64(maskHeadLength) > fileSize {
		return ErrNotDisguised
	}

	// 读取逆序的原始文件头
	var originalHead []byte
	actualOriginalLen := maskHeadLength
	bodySize := fileSize - int64(MaskLengthIndicatorLength) - int64(maskHeadLength)

	if int64(maskHeadLength) <= bodySize {
		// 正常情况：面具长度 <= 真实文件长度
		readPos := fileSize - int64(MaskLengthIndicatorLength) - int64(maskHeadLength)
		originalHead = make([]byte, maskHeadLength)
		_, err = f.ReadAt(originalHead, readPos)
	} else {
		// 非正常情况：面具长度 > 真实文件长度
		actualOriginalLen = int(fileSize - int64(MaskLengthIndicatorLength) - int64(maskHeadLength))
		if actualOriginalLen < 0 {
			return ErrNotDisguised
		}
		originalHead = make([]byte, actualOriginalLen)
		_, err = f.ReadAt(originalHead, int64(maskHeadLength))
	}
	if err != nil {
		return fmt.Errorf("failed to read original file header: %w", err)
	}

	// 截断文件末尾多余部分
	err = f.Truncate(fileSize - int64(maskHeadLength) - int64(MaskLengthIndicatorLength))
	if err != nil {
		return fmt.Errorf("failed to truncate file: %w", err)
	}

	// 写回逆序还原的原始文件头
	restoredHead := reverseBytes(originalHead)
	_, err = f.WriteAt(restoredHead, 0)
	if err != nil {
		return fmt.Errorf("failed to write back original file header: %w", err)
	}

	return nil
}

// RevealToFile 还原伪装文件并输出到新路径（不修改原文件）
func RevealToFile(srcPath, dstPath string) error {
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("failed to read source file: %w", err)
	}
	if err := os.WriteFile(dstPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write destination file: %w", err)
	}
	return Reveal(dstPath)
}

// Inspect 检测文件是否经过 apate 伪装，并返回分析结果
func Inspect(filePath string) (*InspectResult, error) {
	fi, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrFileNotExist
		}
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}
	if fi.IsDir() {
		return nil, ErrIsDirectory
	}

	result := &InspectResult{
		FilePath:    filePath,
		FileSize:    fi.Size(),
		IsDisguised: false,
	}

	if fi.Size() < int64(MaskLengthIndicatorLength) {
		return result, nil
	}

	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	// 读取最后4字节作为面具长度
	maskLenBytes := make([]byte, MaskLengthIndicatorLength)
	_, err = f.ReadAt(maskLenBytes, fi.Size()-int64(MaskLengthIndicatorLength))
	if err != nil {
		return result, nil
	}
	maskHeadLength := int(binary.LittleEndian.Uint32(maskLenBytes))

	// 合理性检查
	if maskHeadLength <= 0 || int64(maskHeadLength) > fi.Size()-int64(MaskLengthIndicatorLength) {
		return result, nil
	}

	// 读取当前文件头用于判断伪装类型
	currentHeadLen := maskHeadLength
	if currentHeadLen > 128 {
		currentHeadLen = 128
	}
	currentHead := make([]byte, currentHeadLen)
	_, err = f.ReadAt(currentHead, 0)
	if err != nil {
		return result, nil
	}

	// 检测伪装类型
	detectedType := detectMaskType(currentHead)

	// 读取逆序的原始文件头
	readPos := fi.Size() - int64(MaskLengthIndicatorLength) - int64(maskHeadLength)
	if readPos < 0 {
		return result, nil
	}
	reversedOriginal := make([]byte, maskHeadLength)
	_, err = f.ReadAt(reversedOriginal, readPos)
	if err != nil {
		return result, nil
	}

	originalHead := reverseBytes(reversedOriginal)

	result.IsDisguised = true
	result.MaskLength = maskHeadLength
	result.OriginalHeader = originalHead
	result.DetectedType = detectedType

	return result, nil
}

// GetAllFilesRecursively 递归遍历路径，获取所有文件
// 与 C# 版本 Program.GetAllFilesRecursively 一致
func GetAllFilesRecursively(path string) ([]string, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if !fi.IsDir() {
		return []string{path}, nil
	}

	var files []string
	err = filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			files = append(files, p)
		}
		return nil
	})
	return files, err
}

// FileToBytes 将文件转换为字节数组，大小受限于 MaximumMaskLength
func FileToBytes(filePath string) ([]byte, error) {
	fi, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}
	if fi.Size() <= 0 || fi.Size() >= int64(MaximumMaskLength) {
		return nil, ErrMaskTooLarge
	}
	return os.ReadFile(filePath)
}

// RenameDisguised 伪装后重命名文件（追加扩展名）
func RenameDisguised(filePath, maskExtension string) (string, error) {
	newPath := filePath + maskExtension
	return newPath, os.Rename(filePath, newPath)
}

// RenameRevealed 还原后重命名文件（去掉最后一个扩展名）
func RenameRevealed(filePath string) (string, error) {
	lastDot := strings.LastIndex(filePath, ".")
	if lastDot <= 0 {
		// 没有扩展名可移除，保持不变
		return filePath, nil
	}
	newPath := filePath[:lastDot]
	return newPath, os.Rename(filePath, newPath)
}

// reverseBytes 将字节数组逆序排列（与 C# ReverseByteArray 一致）
func reverseBytes(buf []byte) []byte {
	result := make([]byte, len(buf))
	for i, b := range buf {
		result[len(buf)-1-i] = b
	}
	return result
}

// detectMaskType 根据文件头字节检测伪装类型
func detectMaskType(head []byte) string {
	if len(head) == 0 {
		return "unknown"
	}

	// 检查 EXE (MZ header)
	if len(head) >= 2 && head[0] == 0x4D && head[1] == 0x5A {
		return "exe"
	}

	// 检查 JPG (FF D8 FF)
	if len(head) >= 3 && head[0] == 0xFF && head[1] == 0xD8 && head[2] == 0xFF {
		return "jpg"
	}

	// 检查 MP4 (ftyp)
	if len(head) >= 8 && head[4] == 0x66 && head[5] == 0x74 && head[6] == 0x79 && head[7] == 0x70 {
		return "mp4"
	}

	// 检查 MOV (moov)
	if len(head) >= 4 && head[0] == 0x6D && head[1] == 0x6F && head[2] == 0x6F && head[3] == 0x76 {
		return "mov"
	}

	return "unknown"
}

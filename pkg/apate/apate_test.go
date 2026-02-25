package apate

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

// TestDisguiseAndReveal 测试伪装后再还原，验证文件内容完全一致
func TestDisguiseAndReveal(t *testing.T) {
	tests := []struct {
		name     string
		maskHead []byte
		content  []byte
	}{
		{"EXE伪装-普通内容", ExeHead, []byte("Hello, World! This is a test file for apate Go CLI.")},
		{"JPG伪装-普通内容", JpgHead, []byte("Another test file with some data 1234567890")},
		{"MP4伪装-普通内容", Mp4Head, []byte("MP4 test content with various characters: !@#$%^&*()")},
		{"MOV伪装-普通内容", MovHead, []byte("MOV disguise test")},
		{"EXE伪装-长内容", ExeHead, bytes.Repeat([]byte("abcdefghij"), 1000)},
		{"JPG伪装-单字节", JpgHead, []byte("x")},
		{"MP4伪装-二进制内容", Mp4Head, []byte{0x00, 0x01, 0x02, 0xFF, 0xFE, 0xFD}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建临时文件
			tmpDir := t.TempDir()
			filePath := filepath.Join(tmpDir, "testfile.bin")
			if err := os.WriteFile(filePath, tt.content, 0644); err != nil {
				t.Fatalf("创建测试文件失败: %v", err)
			}

			// 伪装
			if err := Disguise(filePath, tt.maskHead); err != nil {
				t.Fatalf("伪装失败: %v", err)
			}

			// 验证文件头已被替换
			disguisedData, err := os.ReadFile(filePath)
			if err != nil {
				t.Fatalf("读取伪装文件失败: %v", err)
			}

			headCheckLen := len(tt.maskHead)
			if headCheckLen > len(disguisedData) {
				headCheckLen = len(disguisedData)
			}
			if !bytes.Equal(disguisedData[:headCheckLen], tt.maskHead[:headCheckLen]) {
				t.Error("伪装后文件头不匹配")
			}

			// 还原
			if err := Reveal(filePath); err != nil {
				t.Fatalf("还原失败: %v", err)
			}

			// 验证还原后内容一致
			restoredData, err := os.ReadFile(filePath)
			if err != nil {
				t.Fatalf("读取还原文件失败: %v", err)
			}

			if !bytes.Equal(restoredData, tt.content) {
				t.Errorf("还原后内容不一致\n原始: %v\n还原: %v", tt.content, restoredData)
			}
		})
	}
}

// TestDisguiseAndRevealSmallFile 测试文件比面具小的极端情况
func TestDisguiseAndRevealSmallFile(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tiny.bin")
	content := []byte("ab") // 只有2字节，比 ExeHead (119字节) 小得多

	if err := os.WriteFile(filePath, content, 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 伪装
	if err := Disguise(filePath, ExeHead); err != nil {
		t.Fatalf("伪装小文件失败: %v", err)
	}

	// 还原
	if err := Reveal(filePath); err != nil {
		t.Fatalf("还原小文件失败: %v", err)
	}

	restoredData, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("读取还原文件失败: %v", err)
	}

	if !bytes.Equal(restoredData, content) {
		t.Errorf("还原后内容不一致\n原始: %v\n还原: %v", content, restoredData)
	}
}

// TestOnekeyMask 测试一键伪装（使用内嵌的 mask.mp4）
func TestOnekeyMask(t *testing.T) {
	mask, err := GetOnekeyMask()
	if err != nil {
		t.Fatalf("获取内嵌面具失败: %v", err)
	}
	if len(mask) == 0 {
		t.Fatal("内嵌面具为空")
	}

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "testfile.zip")
	content := []byte("PK\x03\x04 This is a fake zip file content for testing purposes 1234567890")

	if err := os.WriteFile(filePath, content, 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 伪装
	if err := Disguise(filePath, mask); err != nil {
		t.Fatalf("一键伪装失败: %v", err)
	}

	// 还原
	if err := Reveal(filePath); err != nil {
		t.Fatalf("一键还原失败: %v", err)
	}

	restoredData, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("读取还原文件失败: %v", err)
	}

	if !bytes.Equal(restoredData, content) {
		t.Errorf("一键伪装还原后内容不一致")
	}
}

// TestInspect 测试文件伪装检测
func TestInspect(t *testing.T) {
	tmpDir := t.TempDir()

	// 测试正常文件
	normalFile := filepath.Join(tmpDir, "normal.txt")
	if err := os.WriteFile(normalFile, []byte("normal file content"), 0644); err != nil {
		t.Fatal(err)
	}
	result, err := Inspect(normalFile)
	if err != nil {
		t.Fatalf("检测正常文件失败: %v", err)
	}
	if result.IsDisguised {
		t.Error("正常文件不应被检测为伪装")
	}

	// 测试伪装文件
	disguisedFile := filepath.Join(tmpDir, "disguised.bin")
	content := []byte("This is secret content that should be hidden")
	if err := os.WriteFile(disguisedFile, content, 0644); err != nil {
		t.Fatal(err)
	}
	if err := Disguise(disguisedFile, ExeHead); err != nil {
		t.Fatalf("伪装失败: %v", err)
	}

	result, err = Inspect(disguisedFile)
	if err != nil {
		t.Fatalf("检测伪装文件失败: %v", err)
	}
	if !result.IsDisguised {
		t.Error("伪装文件应被检测为已伪装")
	}
	if result.MaskLength != len(ExeHead) {
		t.Errorf("面具长度不匹配: 期望 %d, 实际 %d", len(ExeHead), result.MaskLength)
	}
	if result.DetectedType != "exe" {
		t.Errorf("伪装类型不匹配: 期望 exe, 实际 %s", result.DetectedType)
	}
}

// TestDisguiseToFile 测试非原地伪装
func TestDisguiseToFile(t *testing.T) {
	tmpDir := t.TempDir()
	srcPath := filepath.Join(tmpDir, "source.txt")
	dstPath := filepath.Join(tmpDir, "source.txt.exe")
	content := []byte("source content for copy disguise test")

	if err := os.WriteFile(srcPath, content, 0644); err != nil {
		t.Fatal(err)
	}

	if err := DisguiseToFile(srcPath, dstPath, ExeHead); err != nil {
		t.Fatalf("DisguiseToFile 失败: %v", err)
	}

	// 验证源文件未被修改
	srcData, _ := os.ReadFile(srcPath)
	if !bytes.Equal(srcData, content) {
		t.Error("源文件不应被修改")
	}

	// 还原目标文件并验证
	if err := Reveal(dstPath); err != nil {
		t.Fatalf("还原失败: %v", err)
	}
	dstData, _ := os.ReadFile(dstPath)
	if !bytes.Equal(dstData, content) {
		t.Error("还原后内容不一致")
	}
}

// TestRevealToFile 测试非原地还原
func TestRevealToFile(t *testing.T) {
	tmpDir := t.TempDir()
	srcPath := filepath.Join(tmpDir, "disguised.bin")
	dstPath := filepath.Join(tmpDir, "restored.bin")
	content := []byte("content for RevealToFile test")

	if err := os.WriteFile(srcPath, content, 0644); err != nil {
		t.Fatal(err)
	}
	if err := Disguise(srcPath, Mp4Head); err != nil {
		t.Fatal(err)
	}

	if err := RevealToFile(srcPath, dstPath); err != nil {
		t.Fatalf("RevealToFile 失败: %v", err)
	}

	dstData, _ := os.ReadFile(dstPath)
	if !bytes.Equal(dstData, content) {
		t.Error("RevealToFile 还原后内容不一致")
	}
}

// TestGetAllFilesRecursively 测试递归文件遍历
func TestGetAllFilesRecursively(t *testing.T) {
	tmpDir := t.TempDir()

	// 创建目录结构
	os.MkdirAll(filepath.Join(tmpDir, "sub1", "sub2"), 0755)
	os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("1"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "sub1", "file2.txt"), []byte("2"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "sub1", "sub2", "file3.txt"), []byte("3"), 0644)

	files, err := GetAllFilesRecursively(tmpDir)
	if err != nil {
		t.Fatalf("递归遍历失败: %v", err)
	}

	if len(files) != 3 {
		t.Errorf("期望3个文件，实际 %d 个", len(files))
	}
}

// TestReverseBytes 测试字节数组逆序
func TestReverseBytes(t *testing.T) {
	input := []byte{1, 2, 3, 4, 5}
	expected := []byte{5, 4, 3, 2, 1}
	result := reverseBytes(input)
	if !bytes.Equal(result, expected) {
		t.Errorf("逆序结果不正确: 期望 %v, 实际 %v", expected, result)
	}
}

// TestErrorCases 测试各种错误情况
func TestErrorCases(t *testing.T) {
	// 空面具
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.bin")
	os.WriteFile(filePath, []byte("test"), 0644)

	err := Disguise(filePath, []byte{})
	if err != ErrInvalidMask {
		t.Errorf("空面具应返回 ErrInvalidMask, 实际: %v", err)
	}

	// 文件不存在
	err = Disguise("/nonexistent/file.bin", ExeHead)
	if err != ErrFileNotExist {
		t.Errorf("不存在的文件应返回 ErrFileNotExist, 实际: %v", err)
	}

	err = Reveal("/nonexistent/file.bin")
	if err != ErrFileNotExist {
		t.Errorf("不存在的文件应返回 ErrFileNotExist, 实际: %v", err)
	}
}

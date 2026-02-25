# Gopate

**Gopate** 是一款基于 Go 语言的纯命令行文件格式伪装工具，能够简洁、快速地对文件进行格式伪装与还原。

> 本项目基于 [rippod/apate](https://github.com/rippod/apate) (C#) 重构而来，**完全兼容**其伪装文件格式，可混合使用。

## ✨ 特性

- 🚀 **极速处理** — 支持超大文件，瞬间伪装/还原
- 🔒 **安全加密** — 原始文件头经过加密处理，不易被检测出原始格式
- 🌍 **跨平台** — 支持 Linux / macOS / Windows
- 🔄 **完全兼容** — 与 C# 版 apate 伪装的文件互通
- 📦 **批量处理** — 支持 glob 通配符、递归目录处理
- 🔍 **文件探测** — 检测文件是否经过伪装，查看原始文件头信息
- 📋 **预览模式** — `--dry-run` 预览操作，不实际修改文件
- 🛡️ **安全模式** — 默认生成新文件，`--in-place` 可选原地修改

## 安装

### 从源码编译

```bash
git clone https://github.com/maolei1024/gopate.git
cd gopate
go build -o gopate .
```

### 直接安装

```bash
go install github.com/maolei1024/gopate@latest
```

##  使用说明

### 伪装文件

```bash
# 一键伪装（默认模式，伪装为 MP4）
gopate disguise secret.zip

# 伪装为指定格式
gopate disguise secret.zip --mode exe    # 伪装为 EXE
gopate disguise secret.zip --mode jpg    # 伪装为 JPG
gopate disguise secret.zip --mode mp4    # 伪装为 MP4
gopate disguise secret.zip --mode mov    # 伪装为 MOV

# 使用自定义面具文件
gopate disguise secret.zip --mode mask --mask-file cover.png

# 原地修改（不保留原文件）
gopate disguise secret.zip --mode onekey --in-place

# 批量处理
gopate disguise *.zip --mode onekey --in-place

# 递归处理目录
gopate disguise ./mydir -r --mode mp4 --in-place

# 输出到指定目录
gopate disguise secret.zip -o ./output/
```

### 还原文件

```bash
# 还原文件（生成新文件）
gopate reveal secret.zip.mp4

# 原地还原
gopate reveal secret.zip.mp4 --in-place

# 跳过确认提示
gopate reveal secret.zip.mp4 --in-place -f

# 批量还原
gopate reveal *.mp4 --in-place -f
```

### 检测文件

```bash
# 检测文件是否经过伪装
gopate inspect suspicious.mp4

# 显示详细信息（含原始文件头）
gopate inspect suspicious.mp4 -v
```

### 其他选项

```bash
gopate version          # 显示版本
gopate --help           # 显示帮助
gopate disguise --help  # 显示伪装命令帮助

# 全局选项
-v, --verbose   显示详细输出
-q, --quiet     静默模式，仅显示错误
```

##  伪装模式说明

| 模式 | 说明 | 面具来源 |
|------|------|----------|
| `onekey` | 一键伪装为 MP4（默认，适用大部分场景） | 内嵌 MP4 文件 |
| `mask` | 自定义面具文件伪装 | 用户指定文件 |
| `exe` | 简易伪装为 EXE | PE 文件头 |
| `jpg` | 简易伪装为 JPG | JPEG 文件头 |
| `mp4` | 简易伪装为 MP4 | MP4 文件头 |
| `mov` | 简易伪装为 MOV | MOV 文件头 |

##  注意事项

1. **使用前请务必做好数据备份**
2. 对未经伪装的文件执行还原操作可能导致文件损坏
3. 本软件不得用于非法用途，使用者自行承担一切后果

## 与 apate 的兼容性

Gopate 与 [apate](https://github.com/rippod/apate) (C# 版本 v1.4.2) 采用完全一致的二进制文件格式：

-  Gopate 伪装的文件可以用 apate 还原
-  apate 伪装的文件可以用 Gopate 还原
-  所有面具字节与 apate 逐字节一致

##  许可证

本项目基于 [MIT License](LICENSE) 开源。

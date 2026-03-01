package i18n

var messagesZh = map[string]string{
	// ── 根命令 ───────────────────────────────────────────────
	"root.short": "Gopate - 文件格式伪装工具 (Go CLI)",
	"root.long": `Gopate 是一款基于 Go 语言的纯命令行文件格式伪装工具。
完全兼容 apate (C# 版本) 的文件格式，支持跨平台使用。

功能特性:
  • 支持超大文件，瞬间伪装/还原
  • 支持批量处理、递归目录
  • 原始文件头经过加密处理，不易被检测
  • 完全兼容 apate 伪装的文件，可混合使用`,
	"flag.verbose": "显示详细输出",
	"flag.quiet":   "静默模式，仅显示错误",
	"flag.lang":    "设置语言 (en, zh)",

	// ── 伪装命令 ─────────────────────────────────────────────
	"disguise.use":   "disguise [文件或目录...]",
	"disguise.short": "伪装文件",
	"disguise.long": `对指定文件进行格式伪装。

伪装模式:
  onekey  使用内嵌 MP4 面具文件伪装（默认，适用大部分场景）
  mask    使用自定义面具文件伪装（需配合 --mask-file 参数）
  exe     使用 EXE 文件头伪装
  jpg     使用 JPG 文件头伪装
  mp4     使用 MP4 文件头伪装
  mov     使用 MOV 文件头伪装

示例:
  gopate disguise file.zip                                    # 一键伪装为 MP4
  gopate disguise file.zip --mode exe                         # 伪装为 EXE
  gopate disguise file.zip --mode mask --mask-file cover.png  # 使用自定义面具
  gopate disguise *.zip --mode mp4 --in-place                 # 批量原地伪装
  gopate disguise ./mydir -r --mode onekey                    # 递归伪装目录下所有文件`,
	"flag.mode":       "伪装模式: onekey|mask|exe|jpg|mp4|mov",
	"flag.mask_file":  "自定义面具文件路径（仅 mask 模式需要）",
	"flag.output_dir": "输出目录（默认输出到源文件同目录）",
	"flag.in_place":   "原地修改文件（默认生成新文件）",
	"flag.recursive":  "递归处理目录",
	"flag.dry_run":    "预览模式，不实际修改文件",

	// ── 还原命令 ─────────────────────────────────────────────
	"reveal.use":   "reveal [文件或目录...]",
	"reveal.short": "还原伪装文件",
	"reveal.long": `将经过伪装的文件还原为原始格式。

注意: 如果对未经过伪装的文件执行还原操作，可能会导致文件损坏！
请务必做好数据备份。

示例:
  gopate reveal file.zip.mp4                # 还原单个文件
  gopate reveal *.mp4 --in-place            # 批量原地还原
  gopate reveal ./mydir -r --in-place       # 递归还原目录下所有文件
  gopate reveal file.mp4 -o ./restored/     # 还原到指定目录`,
	"flag.force": "跳过确认提示",

	// ── 检测命令 ──────────────────────────────────────────────
	"inspect.use":   "inspect [文件...]",
	"inspect.short": "检测文件是否经过 apate/gopate 伪装",
	"inspect.long": `分析文件是否经过 apate/gopate 伪装，并显示相关信息。

示例:
  gopate inspect file.mp4              # 检测单个文件
  gopate inspect *.mp4                 # 批量检测
  gopate inspect file.mp4 -v           # 显示详细信息（含原始文件头）`,

	// ── 版本命令 ──────────────────────────────────────────────
	"version.short":  "显示版本信息",
	"msg.compatible": "兼容 apate v1.4.2 文件格式",
	"msg.homepage":   "项目主页: https://github.com/maolei1024/gopate",

	// ── 运行时消息 ───────────────────────────────────────────
	"msg.no_files":             "未找到任何文件",
	"msg.dry_run_disguise":     "预览模式 - 将伪装 %d 个文件 (模式: %s)",
	"msg.dry_run_reveal":       "预览模式 - 将还原 %d 个文件",
	"msg.disguising":           "伪装: %s",
	"msg.revealing":            "还原: %s",
	"msg.failed":               "失败: %s - %v",
	"msg.rename_failed":        "重命名失败: %s - %v",
	"msg.create_outdir_failed": "创建输出目录失败: %w",
	"msg.done":                 "完成！成功 %d 个，失败 %d 个",
	"msg.some_failed":          "有 %d 个文件处理失败",
	"msg.load_mask_failed":     "加载内嵌面具失败: %w",
	"msg.mask_file_required":   "mask 模式需要指定 --mask-file 参数",
	"msg.read_mask_failed":     "读取面具文件失败: %w",
	"msg.unknown_mode":         "未知的伪装模式: %s",
	"msg.unknown_mode_opts":    "未知的伪装模式: %s\n可选: onekey, mask, exe, jpg, mp4, mov",
	"msg.glob_failed":          "路径匹配失败: %w",
	"msg.access_failed":        "无法访问: %s - %w",
	"msg.walk_dir_failed":      "遍历目录失败: %s - %w",
	"msg.skip_dir":             "跳过目录: %s (使用 -r 递归处理)",
	"msg.reveal_warning":       "⚠ 警告: 对未经伪装的文件执行还原可能导致文件损坏！请确保已备份。",
	"msg.confirm_proceed":      "即将处理 %d 个文件，是否继续? (y/N): ",
	"msg.cancelled":            "已取消",
	"msg.inspect_failed":       "检测失败: %s - %v",
	"msg.disguised_yes":        "✅ %s - 已伪装",
	"msg.file_size":            "   文件大小: %s",
	"msg.mask_length":          "   面具长度: %d 字节",
	"msg.disguise_type":        "   伪装类型: %s",
	"msg.original_header":      "   原始文件头: ",
	"msg.disguised_no":         "❌ %s - 未检测到伪装",
}

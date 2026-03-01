package i18n

var messagesEn = map[string]string{
	// ── Root command ──────────────────────────────────────────
	"root.short": "Gopate - File format disguise tool (Go CLI)",
	"root.long": `Gopate is a pure command-line file format disguise tool built with Go.
Fully compatible with apate (C# version) file format, supports cross-platform usage.

Features:
  • Ultra-fast processing for large files
  • Batch processing with glob patterns and recursive directories
  • Encrypted original file headers for stealth
  • Fully compatible with apate disguised files`,
	"flag.verbose": "Show verbose output",
	"flag.quiet":   "Quiet mode, show errors only",
	"flag.lang":    "Set language (en, zh)",

	// ── Disguise command ─────────────────────────────────────
	"disguise.use":   "disguise [files or directories...]",
	"disguise.short": "Disguise files",
	"disguise.long": `Disguise specified files by modifying their format headers.

Disguise modes:
  onekey  Use embedded MP4 mask file (default, suitable for most cases)
  mask    Use a custom mask file (requires --mask-file)
  exe     Use EXE file header
  jpg     Use JPG file header
  mp4     Use MP4 file header
  mov     Use MOV file header

Examples:
  gopate disguise file.zip                                    # One-key disguise as MP4
  gopate disguise file.zip --mode exe                         # Disguise as EXE
  gopate disguise file.zip --mode mask --mask-file cover.png  # Custom mask file
  gopate disguise *.zip --mode mp4 --in-place                 # Batch in-place disguise
  gopate disguise ./mydir -r --mode onekey                    # Recursive disguise`,
	"flag.mode":       "Disguise mode: onekey|mask|exe|jpg|mp4|mov",
	"flag.mask_file":  "Custom mask file path (required for mask mode)",
	"flag.output_dir": "Output directory (defaults to source file directory)",
	"flag.in_place":   "Modify files in-place (default: create new files)",
	"flag.recursive":  "Process directories recursively",
	"flag.dry_run":    "Preview mode, no actual file changes",

	// ── Reveal command ───────────────────────────────────────
	"reveal.use":   "reveal [files or directories...]",
	"reveal.short": "Reveal disguised files",
	"reveal.long": `Restore disguised files to their original format.

WARNING: Revealing files that were not disguised may cause data corruption!
Please make sure to back up your data first.

Examples:
  gopate reveal file.zip.mp4                # Reveal a single file
  gopate reveal *.mp4 --in-place            # Batch in-place reveal
  gopate reveal ./mydir -r --in-place       # Recursive reveal
  gopate reveal file.mp4 -o ./restored/     # Reveal to specified directory`,
	"flag.force": "Skip confirmation prompt",

	// ── Inspect command ──────────────────────────────────────
	"inspect.use":   "inspect [files...]",
	"inspect.short": "Detect if files are disguised by apate/gopate",
	"inspect.long": `Analyze files to detect if they were disguised by apate/gopate.

Examples:
  gopate inspect file.mp4              # Inspect a single file
  gopate inspect *.mp4                 # Batch inspect
  gopate inspect file.mp4 -v           # Show detailed info (with original header)`,

	// ── Version command ──────────────────────────────────────
	"version.short":  "Show version information",
	"msg.compatible": "Compatible with apate v1.4.2 file format",
	"msg.homepage":   "Homepage: https://github.com/maolei1024/gopate",

	// ── Runtime messages ─────────────────────────────────────
	"msg.no_files":             "no files found",
	"msg.dry_run_disguise":     "Preview mode - will disguise %d file(s) (mode: %s)",
	"msg.dry_run_reveal":       "Preview mode - will reveal %d file(s)",
	"msg.disguising":           "Disguising: %s",
	"msg.revealing":            "Revealing: %s",
	"msg.failed":               "Failed: %s - %v",
	"msg.rename_failed":        "Rename failed: %s - %v",
	"msg.create_outdir_failed": "failed to create output directory: %w",
	"msg.done":                 "Done! %d succeeded, %d failed",
	"msg.some_failed":          "%d file(s) failed to process",
	"msg.load_mask_failed":     "failed to load embedded mask: %w",
	"msg.mask_file_required":   "mask mode requires --mask-file parameter",
	"msg.read_mask_failed":     "failed to read mask file: %w",
	"msg.unknown_mode":         "unknown disguise mode: %s",
	"msg.unknown_mode_opts":    "unknown disguise mode: %s\nAvailable: onekey, mask, exe, jpg, mp4, mov",
	"msg.glob_failed":          "path matching failed: %w",
	"msg.access_failed":        "cannot access: %s - %w",
	"msg.walk_dir_failed":      "failed to traverse directory: %s - %w",
	"msg.skip_dir":             "Skipping directory: %s (use -r for recursive processing)",
	"msg.reveal_warning":       "⚠ WARNING: Revealing files that were not disguised may cause corruption! Make sure you have backups.",
	"msg.confirm_proceed":      "About to process %d file(s), continue? (y/N): ",
	"msg.cancelled":            "Cancelled",
	"msg.inspect_failed":       "Inspection failed: %s - %v",
	"msg.disguised_yes":        "✅ %s - Disguised",
	"msg.file_size":            "   File size: %s",
	"msg.mask_length":          "   Mask length: %d bytes",
	"msg.disguise_type":        "   Disguise type: %s",
	"msg.original_header":      "   Original header: ",
	"msg.disguised_no":         "❌ %s - Not disguised",
}

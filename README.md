# Granary

Exports meeting notes and transcripts from [Granola](https://www.granola.so)'s local cache to markdown files.

## 📤 What it exports

- 🤖 **AI-generated meeting notes** - Granola's AI summaries and notes
- 🎙️ **Full transcripts** - Complete word-for-word transcripts when available
- 📄 **Both together** - Files include both notes and transcripts when both exist

## ✨ Features

- ⚡ **Smart caching** - Only writes changed files (efficient for scheduled runs)
- 🔍 **Version detection** - Auto-detects latest Granola cache version (`cache-v3.json`, `cache-v4.json`, etc.)
- 🛡️ **Data protection** - Preserves transcripts even if Granola purges them from cache
- 🕐 **Background service** - Built-in macOS LaunchAgent for automatic exports every 6 hours

## 🛠️ Installation

### Homebrew

```bash
brew install wassimk/tap/granary
```

### From source

```bash
go install github.com/wassimk/granary@latest
```

## 💻 Usage

### Export meeting notes

```bash
granary run
```

#### Options

```
-o, --output-dir   Custom output directory (default: ~/.local/share/granola-transcripts)
```

### Background service (LaunchAgent)

Install a macOS LaunchAgent that automatically exports every 6 hours:

```bash
granary install
```

Check the service status:

```bash
granary status
```

Remove the background service:

```bash
granary uninstall
```

### Other commands

```bash
granary version    # Show version
granary help       # Show help
```

## ⏰ Alternative: cron

If you prefer cron over the built-in LaunchAgent:

```bash
crontab -e
```

Add this line:

```bash
0 */6 * * * /opt/homebrew/bin/granary run >> /tmp/granary.log 2>&1
```

## ⚙️ How it works

- **Reads from:** `~/Library/Application Support/Granola/cache-v*.json` (auto-detects version)
- **Exports to:** `~/.local/share/granola-transcripts/`
- **Format:** Markdown with AI notes section and transcript section (when available)
- **Filename:** `YYYY-MM-DD_Meeting_Title.md`

## ⚠️ Important: Transcript Availability

**Granola does not keep all transcripts in its local cache.** Transcripts are fetched from Granola's servers on-demand when you open a meeting, and older transcripts are periodically purged from the cache.

This means:

- **New meetings:** Transcripts appear in cache after you view them in Granola
- **Old meetings:** Even if you viewed them before, the transcript may no longer be in cache

**Data protection:** Once this tool exports a transcript, it preserves it forever. Even if Granola later purges it from cache. The tool merges the latest AI notes with any previously exported transcript, so you never lose data.

## 📄 Output format

Each exported file contains:

```markdown
# Meeting Title
Date: 2025-01-24 14:30
Meeting ID: abc-123

---

## AI-Generated Notes

[Granola's AI-generated meeting notes and summaries]

---

## Transcript

**Me:** [Your words]

**Them:** [Other participant's words]
```

## 📝 Disclaimer

This project is not affiliated with, endorsed by, or connected to [Granola](https://www.granola.so) in any way. I love Granola and use it every day. This is just a personal utility to export my meeting data.

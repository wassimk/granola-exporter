# Granola Meeting Exporter

> [!WARNING]
> This is early-stage software that works for my personal use case. Expect breaking changes, bugs, and rough edges as I continue to develop it.

Exports meeting notes and transcripts from [Granola](https://www.granola.so)'s local cache to markdown files.

## What it exports

- **AI-generated meeting notes** - Granola's AI summaries and notes
- **Full transcripts** - Complete word-for-word transcripts when available
- **Both together** - Files include both notes and transcripts when both exist

## Features

- **Smart caching** - Only writes changed files (efficient for hourly cron runs)
- **Version detection** - Auto-detects latest Granola cache version (`cache-v3.json`, `cache-v4.json`, etc.)
- **Data protection** - Preserves transcripts even if Granola purges them from cache

## Installation

Download the latest binary from the [Releases](https://github.com/wassimk/granola-exporter/releases) page.

Or build from source:

```bash
go install github.com/wassimk/granola-exporter@latest
```

## Usage

Run the exporter:

```bash
granola-exporter
```

### Options

```
-h, --help         Show help message
-V, --version      Show version number
-o, --output-dir   Custom output directory (default: ~/.local/share/granola-transcripts)
```

## Automated export with cron

Set up hourly automated exports:

```bash
crontab -e
```

Add this line:

```bash
0 * * * * /path/to/granola-exporter >> /tmp/granola-export.log 2>&1
```

## How it works

- **Reads from:** `~/Library/Application Support/Granola/cache-v*.json` (auto-detects version)
- **Exports to:** `~/.local/share/granola-transcripts/`
- **Format:** Markdown with AI notes section and transcript section (when available)
- **Filename:** `YYYY-MM-DD_Meeting_Title.md`

## Important: Transcript Availability

**Granola does not keep all transcripts in its local cache.** Transcripts are fetched from Granola's servers on-demand when you open a meeting, and older transcripts are periodically purged from the cache.

This means:

- **New meetings:** Transcripts appear in cache after you view them in Granola
- **Old meetings:** Even if you viewed them before, the transcript may no longer be in cache

**Data protection:** Once this tool exports a transcript, it preserves it forever—even if Granola later purges it from cache. The tool merges the latest AI notes with any previously exported transcript, so you never lose data.

## Output format

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

## Disclaimer

This project is not affiliated with, endorsed by, or connected to [Granola](https://www.granola.so) in any way. I love Granola and use it every day—this is just a personal utility to export my meeting data.

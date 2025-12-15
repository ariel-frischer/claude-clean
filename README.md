# 🧹 Claude Clean

[![CI](https://github.com/ariel-frischer/claude-clean/actions/workflows/release.yml/badge.svg)](https://github.com/ariel-frischer/claude-clean/actions)
[![Release](https://img.shields.io/github/v/release/ariel-frischer/claude-clean)](https://github.com/ariel-frischer/claude-clean/releases)
[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**Transform messy Claude Code JSON logs into beautiful terminal output.** ✨

<p align="center">
  <img src="https://img.shields.io/badge/colors-yes-brightgreen" alt="Colors">
  <img src="https://img.shields.io/badge/box%20drawing-yes-blue" alt="Box Drawing">
  <img src="https://img.shields.io/badge/real--time-streaming-orange" alt="Streaming">
</p>

---

## ⚡ Quick Start

```bash
# Install (one command!)
curl -fsSL https://raw.githubusercontent.com/ariel-frischer/claude-clean/main/install.sh | sh

# Use it
claude -p "your prompt" --verbose --output-format stream-json | cclean
```

That's it! 🎉

---

## 🤔 Why?

Claude Code's `stream-json` output is **unreadable**:

```json
{"type":"assistant","message":{"content":[{"type":"text","text":"Hello!"}],"usage":{"input_tokens":150,"output_tokens":45}}}
```

With `cclean`, you get **beautiful output**:

```
┌─ ASSISTANT
│ Hello!
└─
```

**Bonus:** Using `-p` with `--output-format stream-json` bypasses Claude's interactive UI, avoiding segfault issues some users experience.

---

## 📦 Installation

### 🚀 One-liner (recommended)
```bash
curl -fsSL https://raw.githubusercontent.com/ariel-frischer/claude-clean/main/install.sh | sh
```

### 📥 Manual Download

| Platform | Download |
|----------|----------|
| 🐧 Linux x64 | [cclean-linux-amd64](https://github.com/ariel-frischer/claude-clean/releases/latest/download/cclean-linux-amd64) |
| 🐧 Linux ARM | [cclean-linux-arm64](https://github.com/ariel-frischer/claude-clean/releases/latest/download/cclean-linux-arm64) |
| 🍎 macOS Intel | [cclean-darwin-amd64](https://github.com/ariel-frischer/claude-clean/releases/latest/download/cclean-darwin-amd64) |
| 🍎 macOS Apple Silicon | [cclean-darwin-arm64](https://github.com/ariel-frischer/claude-clean/releases/latest/download/cclean-darwin-arm64) |
| 🪟 Windows | [cclean-windows-amd64.exe](https://github.com/ariel-frischer/claude-clean/releases/latest/download/cclean-windows-amd64.exe) |

```bash
chmod +x cclean-* && sudo mv cclean-* /usr/local/bin/cclean
```

### 🔧 Build from source
```bash
git clone https://github.com/ariel-frischer/claude-clean && cd claude-clean
make install
```

### 🗑️ Uninstall
```bash
cclean --uninstall
```

---

## 🎯 Usage

```bash
# 📡 Live streaming (most common)
claude -p "your prompt" --verbose --output-format stream-json | cclean

# 📄 From a file
cclean logs.jsonl

# 📥 From stdin
cat logs.jsonl | cclean
```

### 🎨 Output Styles

| Style | Flag | Best For |
|-------|------|----------|
| **Default** | `-s default` | 📦 Boxed, colorful, easy to read |
| **Compact** | `-s compact` | 📊 Single-line, quick scanning |
| **Minimal** | `-s minimal` | 📝 No boxes, still colored |
| **Plain** | `-s plain` | 📋 No colors, great for logs |

```bash
cclean -s compact logs.jsonl  # Try different styles!
```

### 🔧 Options

| Flag | Description |
|------|-------------|
| `-s, --style` | Output style (default/compact/minimal/plain) |
| `-v, --verbose` | Show system reminders |
| `-l, --line-numbers` | Show source line numbers |
| `-V, --usage` | Show token usage stats |

---

## 🌈 What Gets Parsed

| Type | Color | Description |
|------|-------|-------------|
| 🔧 **SYSTEM** | Cyan | Init, config, session info |
| 💬 **ASSISTANT** | Green | Claude's responses |
| 🛠️ **TOOL** | Yellow | Tool calls (Bash, Read, Write...) |
| ✅ **RESULT** | Gray | Tool output |
| ❌ **ERROR** | Red | Tool errors |
| 📊 **USAGE** | Magenta | Token stats & costs |

---

## 📚 Documentation

- [Usage Guide](docs/USAGE.md) - Shell aliases, advanced examples, output styles
- [Installation](docs/INSTALL.md) - Detailed installation options
- [Development](docs/DEVELOPMENT.md) - Building from source, architecture
- [Contributing](docs/CONTRIBUTING.md) - How to contribute
- [Releases](docs/RELEASES.md) - Release process and versioning
- [Auto-bump](docs/AUTOBUMP.md) - Automatic version bumping

---

## 📜 License

MIT © [Ariel Frischer](https://github.com/ariel-frischer)

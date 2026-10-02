<div align="center">

```
┌─────────────────────────────────────────┐
│   🔍  f f                                │
│   one query. your whole filesystem.      │
└─────────────────────────────────────────┘
```

### A stupid-fast fuzzy file finder for the entire filesystem

*No live scanning. No lag. Type three letters, find the file.*

[![Go Version](https://img.shields.io/badge/Go-1.22%2B-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev)
[![Platform](https://img.shields.io/badge/Platform-Linux-1793D1?style=flat-square&logo=linux&logoColor=white)](https://archlinux.org)
[![License](https://img.shields.io/badge/License-GPLv3-blue.svg?style=flat-square)](./LICENSE)
[![Built with Bubble Tea](https://img.shields.io/badge/Built%20with-Bubble%20Tea-FF69B4?style=flat-square)](https://github.com/charmbracelet/bubbletea)

</div>

---

## Why this exists

`find` walks the disk live — slow on anything big. `locate`-family tools fix the speed problem with a pre-built index, but most fuzzy layers on top of them rank matches by fuzzy-matching the **whole path as one string** — which quietly breaks in two ways: a real match buried deep in a long path scores *worse* than a shallow coincidence, and letters from totally unrelated folders can combine into a false positive.

**`ff` fixes both by ranking each path *component* separately** — filename, each folder name — instead of the path as one blob, then rolls the best per-component scores back up with extra weight on filename hits.

<div align="center">

| 📇 plocate | 🧮 fuzzy ranking | 🖥️ Bubble Tea |
|:---:|:---:|:---:|
| instant, pre-indexed candidates | per-component scoring, no false positives | live-filtering terminal UI |

**↓ query in, ranked results out, in milliseconds ↓**

</div>

```
your query
    │
    ▼
 plocate -i <tokens>          →  candidate paths (fast, pre-indexed, no live scan)
    │
    ▼
 per-component fuzzy ranking  →  scores each path segment independently
    │
    ▼
 Bubble Tea TUI                →  you see it, you pick it, you're done
```

---

## 📑 Table of contents

- [Features](#-features)
- [Requirements](#-requirements)
- [Installation](#-installation)
- [Configuration from scratch](#-configuration-from-scratch)
- [Usage](#-usage)
- [Architecture](#-architecture)
- [Repo hygiene](#-repo-hygiene)
- [Troubleshooting](#-troubleshooting)
- [License](#-license)

---

## ✨ Features

| | |
|---|---|
| ⚡ **Instant, whole-filesystem search** | Backed by `plocate`'s pre-built index — no live directory walking, results come back in milliseconds even on huge disks |
| 🎯 **Per-component fuzzy ranking** | Each query word is matched against each path *segment* independently, not the path as one giant string — kills both the distance-from-start penalty and false-positive matches scattered across unrelated folders |
| 🔗 **AND-ed multi-word queries** | `ff cael notes` only shows paths matching *both* words — and typing an actual path (`/usr/share/sddm`) works the same way, split on `/` |
| 🗂️ **Smart open, not dumb open** | `Enter` asks *which app* to open with via `mimeopen` (if installed) instead of silently guessing the default — falls back cleanly to `xdg-open` if it's not |
| ✏️ **Editor integration** | `Ctrl+E` opens the selected path in `$EDITOR`, falling back to `nvim` then `vi` |
| 📋 **Clipboard integration** | `Ctrl+Y` copies the selected path — `wl-copy` on Wayland, falls back to `xclip`/`xsel` on X11 |
| 🔄 **On-demand reindex** | `Ctrl+U` force-refreshes the `plocate` index right from inside the TUI — no alt-tabbing to a shell |

---

## ✅ Requirements

| Requirement | Needed for | Notes |
|---|---|---|
| `plocate` | Everything — the actual search | `sudo pacman -S plocate` |
| `xdg-utils` | `Enter` opening files in their default app | `sudo pacman -S xdg-utils` |
| `wl-clipboard` | `Ctrl+Y` copy-to-clipboard on Wayland | `sudo pacman -S wl-clipboard` — falls back to `xclip`/`xsel` on X11 |
| `perl-file-mimeinfo` *(optional)* | Makes `Enter` **ask** which app to open with | `sudo pacman -S perl-file-mimeinfo` |
| Go 1.22+ *(optional)* | Building from source | Not needed if you use the prebuilt binary |

> **Package names above are Arch/pacman.** On Debian/Ubuntu/Fedora the same tools exist under identical or near-identical names (`plocate`, `xdg-utils`, `wl-clipboard`). Nothing about `ff` itself is Arch-specific — it just shells out to standard XDG tooling.

---

## 🚀 Installation

**Automatic (recommended)**
```bash
git clone https://github.com/yashwanth0072/ff.git
cd ff
./install.sh
```
`install.sh` checks your dependencies, builds from source if you have Go (falls back to the prebuilt `linux/amd64` binary if not), installs to `~/.local/bin/ff`, and warns you if that directory isn't on your `PATH` yet — with the exact line to add for your shell.

**Manual, step by step**

**1 · Clone the repo**
```bash
git clone https://github.com/yashwanth0072/ff.git
cd ff
```

**2 · Build**
```bash
go build -trimpath -ldflags="-s -w" -o ff .
```
> `-trimpath -ldflags="-s -w"` strips local build-machine paths and debug symbols — smaller, cleaner binary.

**3 · Install**
```bash
install -Dm755 ff ~/.local/bin/ff
```
<details>
<summary>No Go toolchain? Use the prebuilt binary instead</summary>

```bash
install -Dm755 ff-prebuilt-linux-amd64 ~/.local/bin/ff
```
</details>

**4 · Verify**
```bash
ff
```
You should land in an empty search box. Type a few letters of any file you know exists — results should filter live.

---

## ⚙️ Configuration from scratch

`ff` has **no config file** — zero-config by design. "Configuration" just means getting the tools it wraps installed. Do this once, top to bottom, on a fresh system.

<table>
<tr><td width="40"><b>1</b></td><td>

**Core search — plocate**
```bash
sudo pacman -S plocate
sudo systemctl enable --now plocate-updatedb.timer
```
That timer keeps the index fresh in the background — you don't touch it again. Force an immediate refresh any time with `sudo updatedb`, or `Ctrl+U` inside `ff`.
</td></tr>
<tr><td><b>2</b></td><td>

**Opening files — xdg-utils**
```bash
sudo pacman -S xdg-utils
```
Verify: `xdg-open --version`
</td></tr>
<tr><td><b>3</b></td><td>

**Clipboard — wl-clipboard**
```bash
sudo pacman -S wl-clipboard
```
On X11 instead of Wayland, install `xclip` or `xsel` instead — `ff` detects and uses whichever is present.
</td></tr>
<tr><td><b>4</b></td><td>

**App picker on Enter — perl-file-mimeinfo** *(optional but recommended)*
```bash
sudo pacman -S perl-file-mimeinfo
```
Without this, `Enter` silently opens the OS default app for that file type. With it, `Enter` shows a numbered list of every registered app and lets you pick — no config, it's auto-detected at runtime.
</td></tr>
<tr><td><b>5</b></td><td>

**Building from source — Go** *(optional)*
```bash
sudo pacman -S go
go version   # confirm 1.22+
```
Skip this entirely if you're fine with the prebuilt binary.
</td></tr>
<tr><td><b>6</b></td><td>

**PATH check**
```bash
# bash/zsh — add to your shell rc file
export PATH="$HOME/.local/bin:$PATH"

# fish
fish_add_path $HOME/.local/bin
```
</td></tr>
</table>

That's it — no YAML, no dotfiles, no environment variables to hand-edit. `ff` auto-detects every optional tool (`mimeopen`, `xclip`/`xsel`) at runtime.

---

## 🕹️ Usage

```bash
ff                # launch with an empty search box
ff cael notes     # launch with a prefilled, AND-ed query
```

| Key | Action |
|---|---|
| *(type)* | filter results live, debounced ~120ms |
| `↑` `↓` / `Ctrl+K` `Ctrl+J` | move the selection |
| `Enter` | open the selected path — `mimeopen` picker if installed, else `xdg-open` |
| `Ctrl+E` | open in `$EDITOR` (falls back `nvim` → `vi`) |
| `Ctrl+Y` | copy the selected path to your clipboard |
| `Ctrl+U` | force-refresh the `plocate` index (`sudo updatedb`) |
| `Esc` / `Ctrl+C` | quit |

> **Example flow:** Run `ff` → type `notes qml` → the list narrows to paths where *both* words fuzzy-match somewhere in the path → `↓` to the one you want → `Enter` → pick the app from the `mimeopen` list (or it just opens if you only have one handler registered). Total time: however fast you can type two words.

**New to fuzzy finders?** You never need the exact, full filename — a handful of letters from anywhere in the name is enough. Typing an actual path works too: `/usr/share/sddm` gets split into `usr`, `share`, `sddm` and matched the same way separate words would be.

---

## 🏗️ Architecture

```
main.go                       entry point
internal/search/
  ├─ plocate.go                shells out to plocate, builds the candidate list
  ├─ tokenize.go                splits queries into AND-ed tokens (on space and /)
  ├─ fuzzy.go                   per-component ranking — the core algorithm
  └─ *_test.go                  regression tests pinning down the ranking behavior
internal/tui/
  ├─ model.go                   Bubble Tea state
  ├─ update.go                  key handling / event loop
  └─ view.go                    rendering
internal/actions/
  ├─ open.go                    Enter: mimeopen (if present) → xdg-open fallback
  ├─ editor.go                  Ctrl+E: $EDITOR → nvim → vi
  ├─ clipboard.go               Ctrl+Y: wl-copy → xclip → xsel
  └─ updatedb.go                Ctrl+U: sudo updatedb
```

Search and ranking are fully decoupled from the TUI — `internal/search` has no Bubble Tea imports at all, so the ranking algorithm is independently testable (see `fuzzy_test.go`) without spinning up any UI.

### What's deliberately *not* in here

- **No live filesystem watching.** `inotify` on the whole tree blows past `fs.inotify.max_user_watches` almost immediately — this is exactly why `locate`-family tools re-index periodically instead of watching live. Need a just-created file findable *now*? `Ctrl+U`.
- **No "reveal in file manager."** Trivial to add later, not worth guessing at up front.
- **No direct reads of the plocate database.** Always shells out to the `plocate` binary — the DB format is internal, and the binary is `setgid` specifically for permission-aware search. Reading the DB directly would silently break that.

---

## 🧼 Repo hygiene

This repo ships **source + a prebuilt `linux/amd64` binary** (for people without a Go toolchain). If you're forking or building on top of this, keep your own local build artifacts out of git — already covered by the included [`.gitignore`](.gitignore):

```gitignore
/ff
/ff.exe
*.test
*.out
```

If you rebuild the prebuilt binary (`ff-prebuilt-linux-amd64`) after a source change, remember to actually commit the new one — a stale prebuilt binary silently ships old behavior to anyone who installs without Go.

---

## 🩺 Troubleshooting

| Problem | Fix |
|---|---|
| Results are empty or stale | `plocate` index isn't fresh — run `sudo updatedb` or press `Ctrl+U` inside `ff` |
| `Enter` never asks which app to use | `mimeopen` isn't installed — `sudo pacman -S perl-file-mimeinfo`, no config needed after |
| `Ctrl+Y` does nothing / errors | No clipboard tool found — install `wl-clipboard` (Wayland) or `xclip`/`xsel` (X11) |
| `ff: command not found` after install | `~/.local/bin` isn't on your `PATH` — see step 6 in [Configuration](#-configuration-from-scratch) |
| Files you can see elsewhere don't show up in `ff` | `plocate` is `setgid` and only shows files your user can read, and respects excludes in `/etc/updatedb.conf` |
| `go build` fails with missing module errors | Run `go mod download` and confirm `go version` is 1.22+ |

---

## 📄 License

<div align="center">

Released under **GPLv3** — see [`LICENSE`](./LICENSE) for full terms.

*For everyone tired of `find / -iname` taking forty seconds.*

</div>

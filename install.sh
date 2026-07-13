#!/usr/bin/env bash
# Installs ff into ~/.local/bin. Builds from source if a Go toolchain is
# present, otherwise falls back to the prebuilt linux/amd64 binary
# shipped alongside this script.
set -euo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"
INSTALL_DIR="${HOME}/.local/bin"
mkdir -p "$INSTALL_DIR"

echo "==> checking runtime deps"
missing=()
command -v plocate >/dev/null 2>&1 || missing+=("plocate")
command -v wl-copy  >/dev/null 2>&1 || missing+=("wl-clipboard")
command -v xdg-open >/dev/null 2>&1 || missing+=("xdg-utils")

if [ "${#missing[@]}" -gt 0 ]; then
    echo "    missing: ${missing[*]}"
    echo "    install with: sudo pacman -S ${missing[*]}"
else
    echo "    all good"
fi

if ! command -v mimeopen >/dev/null 2>&1; then
    echo "    (optional) mimeopen not found — Enter will silently xdg-open"
    echo "    instead of prompting which app to use. Install with:"
    echo "      sudo pacman -S perl-file-mimeinfo"
fi

if command -v go >/dev/null 2>&1; then
    echo "==> Go toolchain found, building from source"
    cd "$SCRIPT_DIR"
    go build -trimpath -ldflags="-s -w" -o "$INSTALL_DIR/ff" .
else
    echo "==> no Go toolchain, using prebuilt binary"
    cp "$SCRIPT_DIR/ff-prebuilt-linux-amd64" "$INSTALL_DIR/ff"
    chmod +x "$INSTALL_DIR/ff"
fi

echo "==> installed to $INSTALL_DIR/ff"

case ":$PATH:" in
    *":$INSTALL_DIR:"*) ;;
    *)
        echo
        echo "!! $INSTALL_DIR is not on your PATH yet."
        if [ -n "${FISH_VERSION:-}" ] || command -v fish >/dev/null 2>&1; then
            echo "    fish detected — run this once:"
            echo "      fish_add_path \$HOME/.local/bin"
            echo "    (persists automatically via fish_variables, no rc edit needed)"
            echo
            echo "    if your interactive shell is actually bash, add this to ~/.bashrc instead:"
            echo "      export PATH=\"\$HOME/.local/bin:\$PATH\""
        else
            echo "    add this to your shell rc (~/.bashrc / ~/.zshrc):"
            echo "      export PATH=\"\$HOME/.local/bin:\$PATH\""
        fi
        ;;
esac

echo
echo "==> checking plocate index freshness"
if command -v plocate >/dev/null 2>&1; then
    if systemctl is-enabled plocate-updatedb.timer >/dev/null 2>&1; then
        echo "    plocate-updatedb.timer is enabled"
    else
        echo "    plocate-updatedb.timer is NOT enabled — ff's results will be stale."
        echo "    run: sudo systemctl enable --now plocate-updatedb.timer"
    fi
fi

echo
echo "Done. Run 'ff' to start, or 'ff <initial query>' to prefill the search."

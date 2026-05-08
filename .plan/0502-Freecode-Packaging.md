# Freecode — Packaging

## 1.0 Purpose

This document specifies packaging for all supported platforms: FreeBSD, Linux (Flatpak), macOS (Homebrew), and IllumOS (OpenSolaris).

---

## 2.0 FreeBSD

### 2.1 Package Format

Use FreeBSD's native packaging system via `pkg` or build a tarball.

### 2.2 Directory Structure

```
packaging/freebsd/
├── freecode/
│   └── pkg/
│       └── +MANIFEST
├── freecode-server/
│   └── pkg/
│       └── +MANIFEST
├── scripts/
│   └── pre-install.sh
├── rc.d/
│   └── freecode
└── freecode.spec
```

### 2.3 pkg-plist

```
bin/freecode
bin/freecode-server
share/man/man1/freecode.1
share/man/man1/freecode-server.1
etc/freecode/config.yaml.example
```

### 2.4 RC Script

```sh
#!/bin/sh
# PROVIDE: freecode
# REQUIRE: NETWORKING
# KEYWORD: shutdown

. /etc/rc.subr

name="freecode"
rcvar="freecode_enable"
command="/usr/local/bin/freecode-server"
pidfile="/var/run/freecode.pid"

load_rc_config $name
run_rc_command "$1"
```

### 2.5 Installation

```bash
# Build package
cd packaging/freebsd
make package

# Install
sudo pkg install ./freecode-*.pkg

# Enable service
sudo sysrc freecode_enable="YES"
sudo service freecode start
```

---

## 3.0 Linux (Flatpak)

### 3.1 Flatpak Manifest

```yaml
# com.freecode.Freecode.yml
app-id: com.freecode.Freecode
runtime: org.freedesktop.Platform
runtime-version: '23.08'
sdk: org.freedesktop.Sdk
sdk-extensions:
  - org.freedesktop.Sdk.Extension.rust
  - org.freedesktop.Sdk.Extension.node20
command: freecode

finish-args:
  - --share=ipc
  - --socket=x11
  - --socket=wayland
  - --socket=pulseaudio
  - --share=network
  - --device=dri
  - --filesystem=home
  - --talk-name=org.freedesktop.Notifications
  - --env=OPENCODE_CONFIG_DIR=/app/config

modules:
  - name: freecode
    buildsystem: simple
    build-commands:
      - install -Dm755 freecode /app/bin/freecode
      - install -Dm644 config.yaml /app/config/config.yaml.example
      - install -Dm644 com.freecode.Freecode.desktop /app/share/applications/com.freecode.Freecode.desktop
      - install -Dm644 com.freecode.Freecode.metainfo.xml /app/share/metainfo/com.freecode.Freecode.metainfo.xml
    sources:
      - type: archive
        url: https://github.com/freecode/releases/latest/download/freecode-linux-amd64.tar.gz
        sha256: ...
```

### 3.2 Desktop Entry

```ini
[Desktop Entry]
Name=Freecode
Comment=AI Coding Assistant
Exec=freecode
Icon=com.freecode.Freecode
Terminal=true
Type=Application
Categories=Development;
Keywords=ai;assistant;coding;
```

### 3.3 AppStream Metadata

```xml
<?xml version="1.0" encoding="UTF-8"?>
<component type="console-application">
  <id>com.freecode.Freecode</id>
  <name>Freecode</name>
  <summary>Platform-independent AI coding assistant</summary>
  <description>
    <p>
      Freecode is a Go-based AI coding assistant that provides
      intelligent code completion, refactoring, and more.
    </p>
  </description>
  <launchable type="command">freecode</launchable>
  <url type="homepage">https://github.com/freecode/freecode</url>
  <provides>
    <binary>freecode</binary>
  </provides>
</component>
```

### 3.4 Build and Install

```bash
# Build Flatpak
flatpak-builder --user --install build-dir com.freecode.Freecode.yml

# Or install from repo
flatpak install flathub com.freecode.Freecode
```

---

## 4.0 macOS (Homebrew)

### 4.1 Tap Structure

```
packaging/macos/
├── Formula/
│   └── freecode.rb
├── freecode-server.plist
└── freecode-launchd.plist
```

### 4.2 Homebrew Formula

```ruby
class Freecode < Formula
  desc "Platform-independent AI coding assistant"
  homepage "https://github.com/freecode/freecode"
  version "1.0.0"

  on_macos do
    on_intel do
      url "https://github.com/freecode/releases/download/v1.0.0/freecode-darwin-amd64.tar.gz"
      sha256 "..."
    end
    on_arm do
      url "https://github.com/freecode/releases/download/v1.0.0/freecode-darwin-arm64.tar.gz"
      sha256 "..."
    end
  end

  def install
    bin.install "freecode"
    bin.install "freecode-server"
    etc.install "config.yaml" => "freecode.yaml.example"
  end

  plist_options startup: true, manual: "freecode-server"

  def plist
    <<~PLIST
      <?xml version="1.0" encoding="UTF-8"?>
      <!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "...">
      <plist version="1.0">
        <dict>
          <key>Label</key>
          <string>com.freecode.server</string>
          <key>ProgramArguments</key>
          <array>
            <string>#{opt_bin}/freecode-server</string>
            <string>--port=18792</string>
            <string>--host=127.0.0.1</string>
          </array>
          <key>RunAtLoad</key>
          <true/>
          <key>KeepAlive</key>
          <true/>
        </dict>
      </plist>
    PLIST
  end

  test do
    system "#{bin}/freecode", "--version"
  end
end
```

### 4.3 Installation

```bash
# Add tap
brew tap freecode/tap https://github.com/freecode/homebrew-tap

# Install
brew install freecode

# Start service
brew services start freecode
```

### 4.4 Service Binding

The macOS service MUST bind to localhost only:

```ruby
def plist
  <<~PLIST
    <array>
      <string>#{opt_bin}/freecode-server</string>
      <string>--port=18792</string>
      <string>--host=127.0.0.1</string>
    </array>
  PLIST
end
```

---

## 5.0 IllumOS (OpenSolaris)

### 5.1 Package Structure

```
packaging/illuminos/
├── freecode/
│   ├── prototype
│   ├── pkginfo
│   └── postinstall
├── freecode-server/
│   └── prototype
└── freecode.p5m
```

### 5.2 IPS Manifest (p5m)

```xml
<?xml version="1.0" encoding="UTF-8"?>
<manifest xmlns:p5m="http://xmlns.projectpantheon.org/5manifest/">
  <transform file => pathvar \
    add_privs=(basic) \
    owner=root \
    group=bin />

  <file path="usr/bin/freecode" mode="0555"/>
  <file path="usr/bin/freecode-server" mode="0555"/>
  <file path="etc/freecode/config.yaml" preserve="none"/>
  <file path="etc/smf/manifests/freecode-server.xml"/>
</manifest>
```

### 5.3 SMF Manifest

```xml
<?xml version="1.0"?>
<service manifest-type="service" name="application/freecode" version="1">
  <dependency grouping="require_all" restart_on="refresh" type="service">
    <service_fmri value="svc:/networkloopback"/>
  </dependency>
  <exec_method type="method" exec="/usr/bin/freecode-server --port=18792 --host=127.0.0.1" timeout_seconds="60"/>
  <property_instance fault_segment="" restart_segment="">
    <property name="host" type="astring">
      <astring_value>127.0.0.1</astring_value>
    </property>
  </property_instance>
  <stability value="Unstable"/>
  <template>
    <common_name>
      <loctx_translate value="Freecode Server"/>
    </common_name>
  </template>
</service>
```

### 5.4 Installation

```bash
# Build package
pkgsend publish -d packaging/illuminos/freecode -p prototype

# Install
pkg install freecode

# Enable
svcadm enable freecode-server
```

---

## 6.0 GoReleaser Configuration

### 6.1 .goreleaser.yaml

```yaml
projectName: freecode

before:
  hooks:
    - go mod download
    - go generate ./...

builds:
  - id: freecode-cli
    main: ./cmd/freecode
    binary: freecode
    env:
      - CGO_ENABLED=0
    goos:
      - freebsd
      - linux
      - darwin
    goarch:
      - amd64
      - arm64

  - id: freecode-server
    main: ./cmd/freecode-server
    binary: freecode-server
    env:
      - CGO_ENABLED=0
    goos:
      - freebsd
      - linux
      - darwin
    goarch:
      - amd64
      - arm64

archives:
  - id: default
    format: tar.gz
    format_overrides:
      - goos: windows
        format: zip
    files:
      - README.md
      - LICENSE
      - config.yaml.example

checksum:
  name_template: 'checksums.txt'

snapshot:
  name_template: "{{ .Tag }}-next"

release:
  github:
    owner: freecode
    name: freecode
  prerelease: auto
```

---

## 7.0 Versioning

- Follow Semantic Versioning (SemVer)
- Use `v1.0.0` tag format
- Generate checksums for all artifacts
- Sign releases with GPG (optional)

---

## 8.0 Installation Scripts

### 8.1 Install Script (Unix)

```bash
#!/bin/sh
set -e

INSTALL_DIR="${HOME}/.local/bin"
REPO="freecode/freecode"

# Detect OS
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
esac

# Download
URL="https://github.com/${REPO}/releases/latest/download/freecode-${OS}-${ARCH}.tar.gz"
curl -sL "$URL" | tar -xz -C "$INSTALL_DIR"

# Make executable
chmod +x "${INSTALL_DIR}/freecode"
chmod +x "${INSTALL_DIR}/freecode-server"

echo "Installed to ${INSTALL_DIR}"
```

---

## 9.0 Post-Installation

### 9.1 First Run Setup

```bash
# Create config directory
mkdir -p ~/.config/freecode

# Copy example config
cp /etc/freecode/config.yaml.example ~/.config/freecode/config.yaml

# Or let freecode create default config
freecode config init
```

### 9.2 Config Location

- Global: `~/.config/freecode/config.yaml`
- Project: `.freecode/config.yaml`
- Legacy (read-only): `~/.config/opencode/config.json`

---

## 10.0 Documentation

### 10.1 Man Pages

All platforms ship man pages in section 1 (user commands):

| File | Section | Description |
|------|---------|-------------|
| `freecode.1` | 1 | Main CLI command |
| `freecode-run.1` | 1 | `freecode run` subcommand |
| `freecode-config.1` | 1 | Configuration guide |
| `freecode-agents.1` | 1 | Built-in agents reference |
| `freecode-config.5` | 5 | Config file format |
| `freecode-yolo.1` | 1 | YOLO mode guide |

### 10.2 Man Page Format

Use mdoc (BSD man) format:

```md
.\" Section 1 - User Commands
.TH FREECODE 1 "2026-05-01" "1.0.0" "Freecode Manual"
.SH NAME
freecode \- AI coding assistant
.SH SYNOPSIS
.B freecode
[\fIcommand\fR] [\fIoptions\fR]
.SH DESCRIPTION
.B freecode
is a platform-independent AI coding assistant.
.SH OPTIONS
.TP
.B \fB\-y\fR, \fB\-\-yolo\fR
Enable YOLO mode (skip confirmations).
.TP
.B \fB\-m\fR, \fB\-\-model\fR \fImodel\fR
Specify the model to use.
.SH FILES
.TP
.B ~/.config/freecode/config.yaml
User configuration file.
.SH SEE ALSO
.BR freecode-config (5),
.BR freecode-agents (1)
```

### 10.3 Platform Documentation Locations

| Platform | Man Pages | Other Docs |
|---------|-----------|------------|
| FreeBSD | `/usr/share/man/man1/` | `/usr/share/doc/freecode/` |
| Linux | `/usr/share/man/man1/` | `/usr/share/doc/freecode/` |
| macOS | `/usr/local/share/man/` | `/usr/local/share/doc/freecode/` |
| IllumOS | `/usr/share/man/` | `/usr/share/doc/freecode/` |

### 10.4 Manual Sections

| Section | Content |
|---------|---------|
| 1 | User commands (freecode, freecode-run, etc.) |
| 5 | Configuration file format (freecode-config.5) |
| 7 | Miscellaneous (freecode-hacking, freecode-license) |

### 10.5 Building Man Pages

```bash
# Generate man pages from source
go generate ./docs/man

# Or use pandoc
pandoc README.md -t man -o freecode.1

# Install man pages
sudo make install-man
```

### 10.6 Generated Documentation

| Source | Output |
|--------|--------|
| `docs/*.md` | HTML in `/usr/share/doc/freecode/` |
| `CHANGELOG.md` | HTML in `/usr/share/doc/freecode/` |
| `README.md` | HTML in `/usr/share/doc/freecode/` |

---

**Author:** Mark LaPointe <mark@cloudbsd.org>
**Last Updated:** 2026-05-01

---

## Author Policy

- **Author:** Mark LaPointe <mark@cloudbsd.org>
- **No trailers**: No `Co-authored-by`, `Sponsored-by`, or similar trailers
- **No sponsorships**: No funding acknowledgments
- **No co-authors**: All commits made solely by the author

# Freecode — Platform-Specific Code

## 1.0 Purpose

This document details platform-specific implementations for FreeBSD 16, macOS, Linux, and IllumOS (OpenSolaris).

---

## 2.0 FreeBSD 16 (Primary)

### 2.1 Platform Identification

```go
// internal/platform/freebsd.go
package platform

const (
    OS     = "freebsd"
    GOOS   = "freebsd"
    Kernel = "FreeBSD"
)
```

### 2.2 Shell Detection

```go
func DetectShell() string {
    // FreeBSD defaults
    if os.Getenv("SHELL") != "" {
        return os.Getenv("SHELL")
    }

    // Check common shells
    for _, shell := range []string{
        "/usr/local/bin/bash",
        "/usr/local/bin/zsh",
        "/bin/bash",
        "/bin/sh",
    } {
        if _, err := os.Stat(shell); err == nil {
            return shell
        }
    }

    return "/bin/sh"
}
```

### 2.3 PTY Handling

```go
import "golang.org/x/sys/unix"

func OpenPTY() (master, slave int, err error) {
    // FreeBSD uses ptsname(3)
    return unix.Openpty(0, 0, nil, nil, nil)
}
```

### 2.4 File Watcher

```go
// FreeBSD uses kqueue
func NewWatcher() (*Watcher, error) {
    kq, err := unix.Kqueue()
    if err != nil {
        return nil, err
    }
    return &Watcher{kq: kq}, nil
}
```

### 2.5 Path Handling

```go
var PathSeparator = "/"
var PathListSeparator = ":"

// Realpath handles FreeBSD-specific symlink resolution
func Realpath(path string) (string, error) {
    return filepath.EvalSymlinks(path)
}
```

### 2.6 Home Directory

```go
func UserHomeDir() string {
    if os.Getenv("HOME") != "" {
        return os.Getenv("HOME")
    }
    // FreeBSD uses pw utility
    out, _ := exec.Command("pw", "user", "show", "-n", os.Getenv("USER")).Output()
    // Parse homedir from output
    return "/home/" + os.Getenv("USER")
}
```

---

## 3.0 macOS

### 3.1 Platform Identification

```go
const (
    OS     = "darwin"
    GOOS   = "darwin"
    Kernel = "Darwin"
)
```

### 3.2 Shell Detection

```go
func DetectShell() string {
    if os.Getenv("SHELL") != "" {
        return os.Getenv("SHELL")
    }

    // macOS Catalina+ uses zsh by default
    if _, err := os.Stat("/bin/zsh"); err == nil {
        return "/bin/zsh"
    }
    return "/bin/bash"
}
```

### 3.3 PTY Handling

```go
import "golang.org/x/sys/unix"

func OpenPTY() (master, slave int, err error) {
    // macOS uses ptsname(3)
    return unix.Openpty(0, 0, nil, nil, nil)
}
```

### 3.4 File Watcher

```go
// macOS uses FSEvents
func NewWatcher() (*Watcher, error) {
    // Use FSEvents API via CGO or wrapper library
    return fsevents.NewWatcher()
}
```

### 3.5 Notifications

```go
// macOS uses Notification Center
func SendNotification(title, body string) error {
    script := fmt.Sprintf(`display notification "%s" with title "%s"`, body, title)
    return exec.Command("osascript", "-e", script).Run()
}
```

---

## 4.0 Linux

### 4.1 Platform Identification

```go
const (
    OS     = "linux"
    GOOS   = "linux"
    Kernel = "Linux"
)
```

### 4.2 Shell Detection

```go
func DetectShell() string {
    if os.Getenv("SHELL") != "" {
        return os.Getenv("SHELL")
    }

    // Check /etc/shells
    data, _ := os.ReadFile("/etc/shells")
    shells := strings.Split(string(data), "\n")
    for _, s := range shells {
        if strings.HasPrefix(s, "#") || s == "" {
            continue
        }
        if _, err := os.Stat(s); err == nil {
            return s
        }
    }

    return "/bin/bash"
}
```

### 4.3 PTY Handling

```go
import "golang.org/x/sys/unix"

func OpenPTY() (master, slave int, err error) {
    // Linux uses devpts
    return unix.Openpty(0, 0, nil, nil, nil)
}
```

### 4.4 File Watcher

```go
// Linux uses inotify
func NewWatcher() (*Watcher, error) {
    fd, err := unix.InotifyInit()
    if err != nil {
        return nil, err
    }
    return &Watcher{fd: fd}, nil
}
```

### 4.5 Notifications

```go
// Linux uses libnotify
func SendNotification(title, body string) error {
    return exec.Command("notify-send", title, body).Run()
}
```

---

## 5.0 IllumOS (OpenSolaris)

### 5.1 Platform Identification

```go
const (
    OS     = "illumos"
    GOOS   = "illumos"
    Kernel = "SunOS"
)
```

### 5.2 Shell Detection

```go
func DetectShell() string {
    if os.Getenv("SHELL") != "" {
        return os.Getenv("SHELL")
    }

    // IllumOS uses ksh93 by default
    for _, shell := range []string{
        "/usr/bin/ksh93",
        "/usr/bin/ksh",
        "/bin/sh",
    } {
        if _, err := os.Stat(shell); err == nil {
            return shell
        }
    }

    return "/usr/bin/ksh93"
}
```

### 5.3 PTY Handling

```go
import "golang.org/x/sys/unix"

func OpenPTY() (master, slave int, err error) {
    // IllumOS uses ptmx/pts
    return unix.Openpty(0, 0, nil, nil, nil)
}
```

### 5.4 File Watcher

```go
// IllumOS uses portfs (Event Ports)
func NewWatcher() (*Watcher, error) {
    port, err := unix.PortCreate()
    if err != nil {
        return nil, err
    }
    return &Watcher{port: port}, nil
}
```

### 5.5 Notifications

```go
// IllumOS uses dtrace or custom notification
// No standard notification daemon, so this is a no-op or uses DBus
func SendNotification(title, body string) error {
    // Try DBus first
    if err := exec.Command("dbus-send", "--session", "--dest=org.freedesktop.Notifications",
        "--type=method_call", "/org/freedesktop/Notifications",
        "org.freedesktop.Notifications.Notify", "string:Freecode",
        "uint32:0", "string:", "string:"+title, "string:"+body).Run(); err != nil {
        return nil // Silently fail on IllumOS
    }
    return nil
}
```

---

## 6.0 Cross-Platform Utilities

### 6.1 Platform Detection

```go
// internal/platform/platform.go
package platform

import (
    "runtime"
    "strings"
)

var (
    GOOS   = runtime.GOOS
    GOARCH = runtime.GOARCH
)

func IsFreeBSD() bool { return GOOS == "freebsd" }
func IsDarwin() bool  { return GOOS == "darwin" }
func IsLinux() bool   { return GOOS == "linux" }
func IsIllumOS() bool { return GOOS == "illumos" || GOOS == "solaris" }
```

### 6.2 Path Operations

```go
import (
    "path/filepath"
    "strings"
)

func ExpandHome(path string) string {
    if strings.HasPrefix(path, "~/") {
        return filepath.Join(os.Getenv("HOME"), path[2:])
    }
    return path
}

func IsAbs(path string) bool {
    return filepath.IsAbs(path)
}

func Abs(path string) (string, error) {
    return filepath.Abs(path)
}
```

### 6.3 OS-Related Exec

```go
import (
    "os/exec"
    "runtime"
)

func Command(name string, args ...string) *exec.Cmd {
    cmd := exec.Command(name, args...)

    // Set environment based on platform
    cmd.Env = os.Environ()

    // macOS specific
    if IsDarwin() {
        cmd.Env = append(cmd.Env, "TERM=xterm-256color")
    }

    return cmd
}
```

---

## 7.0 Platform-Specific Features

### 7.1 Feature Matrix

| Feature | FreeBSD | macOS | Linux | IllumOS |
|---------|---------|-------|-------|---------|
| PTY | ✓ | ✓ | ✓ | ✓ |
| inotify/FSEvents/kqueue | kqueue | FSEvents | inotify | portfs |
| Notifications | ✓ | ✓ | ✓ | ✗ |
| Socket Activation | ✗ | ✗ | ✓ | ✓ |
| systemd | ✗ | ✗ | Optional | ✗ |
| SMF | ✗ | ✗ | ✗ | ✓ |

### 7.2 Fallback Behavior

```go
// If a feature is not available, use fallback
func NewWatcher() (Watcher, error) {
    switch {
    case IsLinux():
        return newInotifyWatcher()
    case IsDarwin():
        return newFSEventsWatcher()
    case IsFreeBSD():
        return newKqueueWatcher()
    case IsIllumOS():
        return newPortWatcher()
    default:
        return newPollWatcher() // Generic fallback
    }
}
```

---

## 8.0 Build Tags

### 8.1 Build Tag Files

```go
// internal/platform/freebsd.go
//go:build freebsd
// +build freebsd

package platform
```

```go
// internal/platform/darwin.go
//go:build darwin
// +build darwin

package platform
```

```go
// internal/platform/linux.go
//go:build linux
// +build linux

package platform
```

```go
// internal/platform/illuminos.go
//go:build illumos || solaris
// +build illumos solaris

package platform
```

### 8.2 Common Code

```go
// internal/platform/common.go
//go:build !freebsd && !darwin && !linux && !illumos && !solaris
// +build !freebsd,!darwin,!linux,!illumos,!solaris

package platform

// Generic fallback implementations
```

---

## 9.0 Testing

### 9.1 Platform Tests

```go
// internal/platform/platform_test.go
func TestShellDetection(t *testing.T) {
    shell := DetectShell()
    if shell == "" {
        t.Error("shell should not be empty")
    }
}
```

### 9.2 Skip Tests on Wrong Platform

```go
// internal/platform/freebsd_test.go
//go:build freebsd

package platform

func TestFreeBSDSpecific(t *testing.T) {
    // Only runs on FreeBSD
}
```

---

**Author:** Mark LaPointe <mark@cloudbsd.org>
**Last Updated:** 2026-05-01

---

## Author Policy

- **Author:** Mark LaPointe <mark@cloudbsd.org>
- **No trailers**: No `Co-authored-by`, `Sponsored-by`, or similar trailers
- **No sponsorships**: No funding acknowledgments
- **No co-authors**: All commits made solely by the author

# Linux Flatpak Packaging

## Status: Planned

## Package Name
`com.freecode.Freecode`

## Manifest
`com.freecode.Freecode.yaml`

## Build Instructions
```bash
flatpak-builder --user --install build-dir com.freecode.Freecode.yaml
```

## Runtime
- FLATPAK_VERSION: 23.08
- SDK: org.freedesktop.Sdk
- Runtime: org.freedesktop.Platform

## Notes
Flatpak provides sandboxed execution on Linux.

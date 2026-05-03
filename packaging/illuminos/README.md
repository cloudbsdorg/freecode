# IllumOS Tarball Packaging

## Status: Planned

## Package Name
`freecode`

## Build Instructions
```bash
cd packaging/illuminos
make package
```

## Installation
```bash
pkgadd -d freecode.pkg
```

## Dependencies
- Go 1.24+

## Notes
IllumOS uses the same packaging system as FreeBSD (pkg).

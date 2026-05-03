# FreeBSD Packaging

## Status: Planned

## Package Name
`freecode`

## Dependencies
- Go 1.24+

## Build Instructions
```bash
cd packaging/freebsd
make pkg
```

## Installation
```bash
pkg install freecode
```

## Notes
FreeBSD packages are built using poudriere or Synth.

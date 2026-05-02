# Main Makefile for Freecode
# Detects OS and delegates to platform-specific Makefiles in .plan/

UNAME_S := $(shell uname -s)

ifeq ($(UNAME_S),Linux)
    PLATFORM_MAKEFILE := .plan/Makefile.linux
endif
ifeq ($(UNAME_S),FreeBSD)
    PLATFORM_MAKEFILE := .plan/Makefile.freebsd
endif
ifeq ($(UNAME_S),Darwin)
    PLATFORM_MAKEFILE := .plan/Makefile.macos
endif
ifeq ($(UNAME_S),SunOS)
    PLATFORM_MAKEFILE := .plan/Makefile.illumos
endif

# Default if not detected
PLATFORM_MAKEFILE ?= .plan/Makefile.linux

.PHONY: all build test clean install uninstall fmt tidy package

all build test clean install uninstall fmt tidy package:
	@$(MAKE) -f $(PLATFORM_MAKEFILE) $@
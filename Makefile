# Main Makefile for Freecode
# Detects OS and delegates to platform-specific Makefiles in .plan/

UNAME_S := $(shell uname -s)

ifeq ($(UNAME_S),Linux)
    PLATFORM_MAKEFILE := make/linux.mk
endif
ifeq ($(UNAME_S),FreeBSD)
    PLATFORM_MAKEFILE := make/freebsd.mk
endif
ifeq ($(UNAME_S),Darwin)
    PLATFORM_MAKEFILE := make/macos.mk
endif
ifeq ($(UNAME_S),SunOS)
    PLATFORM_MAKEFILE := make/illumos.mk
endif

# Default if not detected
PLATFORM_MAKEFILE ?= make/linux.mk

.PHONY: all build test clean install uninstall fmt tidy package

all build test clean install uninstall fmt tidy package:
	@$(MAKE) -f $(PLATFORM_MAKEFILE) $@
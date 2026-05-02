#!/bin/sh

PORTNAME=freecode
PORTVERSION=0.1.0
CATEGORIES=devel ai chat assistant
MASTER_SITES=https://github.com/cloudbsdorg/freecode/releases/download/v${PORTVERSION}/
DISTNAME=freecode-${PORTVERSION}

MAINTAINER=mark@cloudbsd.org
COMMENT="Unified AI coding assistant"

RUN_DEPENDS=
BUILD_DEPENDS=go124:lang/go124

USE_GITHUB=yes
GH_ACCOUNT=cloudbsdorg
GH_PROJECT=freecode
GH_TAGNAME=v${PORTVERSION}

USES=go: modules
GO_TARGET=./cmd/freecode
GO_BUILDFLAGS=-ldflags="-s -w"

.include <bsd.port.mk>

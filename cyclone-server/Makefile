#Makefile for Jenkins
#Copyright (C) to cloudboot
PREFIX=$(shell pwd)
BINDIR=$(PREFIX)/bin
CONFDIR=$(PREFIX)/doc/conf
HOSTOS=$(shell go env GOHOSTOS)
HOSTARCH=$(shell go env GOHOSTARCH)
GOOSS=$(GOOS)

DST=$(PREFIX)/cmd


BINDIRLNAMD64=$(BINDIR)/linux_amd64
BINDIRLNARM64=$(BINDIR)/linux_arm64

ifeq (${HOSTOS}_${HOSTARCH}, linux_amd64) 
	BINDIRLNAMD64=$(BINDIR)
else ifeq (${HOSTOS}_${HOSTARCH}, linux_arm64) 
	BINDIRLNARM64=$(BINDIR)
endif

build:
	GOOS=linux GOARCH=amd64 gbb --debug

	@mkdir -p ${DST}/bin/
	@mkdir -p ${DST}/etc/

	@cp -rf $(CONFDIR)/*  ${DST}/etc/
	@cp -rf $(BINDIRLNAMD64)/cloudboot-server ${DST}/bin/
	@cp -rf $(BINDIRLNAMD64)/cloudboot-agent ${DST}/bin/
	@cp -rf $(BINDIRLNAMD64)/hw-server ${DST}/bin/

	GOOS=linux GOARCH=arm64 gbb --debug

	@cp -rf $(BINDIRLNARM64)/cloudboot-agent ${DST}/bin/cloudboot-agent-arm64
	@cp -rf $(BINDIRLNARM64)/hw-server ${DST}/bin/hw-server-arm64

clean:
	rm -rf $(BINDIR)
	rm -rf ${PREFIX}/cmd/

.PHONY: build  clean
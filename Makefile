go_src				:= $(shell find . -type f -name '*.go' -not -path 'vendor/*')
lintfolders			:= cmd internal pkg
cmd 				?= tinyrstats
skip_dep			?= true
go_bin				:= ${GOPATH}/bin

include scripts/Makefile.inc

PHONY: usage clean linter format build dep_ensure $(cmd)

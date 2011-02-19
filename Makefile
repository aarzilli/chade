include $(GOROOT)/src/Make.inc

TARG=chade
GOFILES=chade.go interpreters.go encoders.go decoders.go entities.go unicode_base.go tests.go

include $(GOROOT)/src/Make.cmd

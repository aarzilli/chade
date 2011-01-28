include $(GOROOT)/src/Make.inc

TARG=chade
GOFILES=chade.go interpreters.go encoders.go

include $(GOROOT)/src/Make.cmd

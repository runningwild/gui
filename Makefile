# Copyright 2010 David Roundy, roundyd@physics.oregonstate.edu.
# All rights reserved.

include $(GOROOT)/src/Make.inc

TARG=github.com/droundy/widgets

DEPS=websocket

GOFILES=\
	widgets.go

include $(GOROOT)/src/Make.pkg


include $(GOROOT)/src/Make.inc

TARG=gettingStarted testdraw2d testPath
OFILES=$(TARG:%=%.$O)

all: $(TARG)

$(TARG): %: %.$O
        $(LD) -o $@ $<

$(OFILES): %.$O: %.go Makefile
        $(GC) -o $@ $<

clean:
        rm -f *.[$(OS)] $(TARG) $(CLEANFILES)

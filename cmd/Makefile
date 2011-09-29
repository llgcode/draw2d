
include $(GOROOT)/src/Make.inc

TARG=gettingStarted testdraw2d testX11draw testandroid testgopher testimage testpostscript #draw2dgl

OFILES=$(TARG:%=%.$O)

all: $(TARG)

$(TARG): %: %.$O
	$(LD) -o $@ $<

$(OFILES): %.$O: %.go Makefile
	$(GC) -o $@ $<

clean:
	rm -f *.[$(OS)] $(TARG) $(CLEANFILES)


include $(GOROOT)/src/Make.inc

TARG=testpostscript
OFILES=$(TARG:%=%.$O)

all: $(TARG)

$(TARG): %: %.$O
	$(LD) -o $@ $<

$(OFILES): %.$O: %.go Makefile
	$(GC) -o $@ $<

clean:
	rm -f *.[$(OS)] $(TARG) $(CLEANFILES)

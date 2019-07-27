TOPTARGETS := all clean build test cross-compile
SUBDIRS := $(wildcard */.)

$(TOPTARGETS): $(SUBDIRS)

all: $(SUBDIRS)
$(SUBDIRS):
	$(MAKE) -C $@ $(MAKECMDGOALS)

.PHONY: all $(SUBDIRS)

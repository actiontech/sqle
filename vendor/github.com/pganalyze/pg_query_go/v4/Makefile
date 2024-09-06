.PHONY: default build test benchmark update_source clean

default: test

build:
	go build

test: build
	go test -v ./

benchmark:
	go build -a
	go test -test.bench=. -test.run=XXX -test.benchtime 10s -test.benchmem -test.cpu=4
	#go test -c -o benchmark
	#GODEBUG=schedtrace=100 ./benchmark -test.bench=BenchmarkRawParseCreateTableParallel -test.run=XXX -test.benchtime 20s -test.benchmem -test.cpu=16

# --- Below only needed for releasing new versions

LIB_PG_QUERY_TAG = 15-4.2.1

root_dir := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
LIB_TMPDIR = $(root_dir)/tmp
LIBDIR = $(LIB_TMPDIR)/libpg_query
LIBDIRGZ = $(TMPDIR)/libpg_query-$(LIB_PG_QUERY_TAG).tar.gz

$(LIBDIR): $(LIBDIRGZ)
	mkdir -p $(LIBDIR)
	cd $(LIB_TMPDIR); tar -xzf $(LIBDIRGZ) -C $(LIBDIR) --strip-components=1

$(LIBDIRGZ):
	mkdir -p $(LIB_TMPDIR)
	curl -o $(LIBDIRGZ) https://codeload.github.com/pganalyze/libpg_query/tar.gz/$(LIB_PG_QUERY_TAG)

update_source: $(LIBDIR)
	rm -f parser/*.{c,h}
	rm -fr parser/include
	# Reduce everything down to one directory
	cp -a $(LIBDIR)/src/* parser/
	rm parser/pg_query_outfuncs_protobuf_cpp.cc
	mv parser/postgres/* parser/
	rmdir parser/postgres
	cp -a $(LIBDIR)/pg_query.h parser/include
	# Make sure every .c in the top-level directory is its own translation unit
	mv parser/*{_conds,_defs,_helper}.c parser/include
	# Protobuf definitions
	mkdir -p $(PWD)/bin
	GOBIN=$(PWD)/bin go install github.com/golang/protobuf/protoc-gen-go
	PATH=$(PWD)/bin:$(PATH) protoc --proto_path=$(LIBDIR)/protobuf --go_out=. $(LIBDIR)/protobuf/pg_query.proto
	mkdir -p parser/include/protobuf
	cp -a $(LIBDIR)/protobuf/*.h parser/include/protobuf
	cp -a $(LIBDIR)/protobuf/*.c parser/
	# Protobuf library code
	mkdir -p parser/include/protobuf-c
	cp -a $(LIBDIR)/vendor/protobuf-c/*.h parser/include
	cp -a $(LIBDIR)/vendor/protobuf-c/*.h parser/include/protobuf-c
	cp -a $(LIBDIR)/vendor/protobuf-c/*.c parser/
	# xxhash library code
	mkdir -p parser/include/xxhash
	cp -a $(LIBDIR)/vendor/xxhash/*.h parser/include
	cp -a $(LIBDIR)/vendor/xxhash/*.h parser/include/xxhash
	cp -a $(LIBDIR)/vendor/xxhash/*.c parser/
	# Other support files
	rm -fr testdata
	cp -a $(LIBDIR)/testdata testdata

clean:
	-@ $(RM) -r $(LIB_TMPDIR)

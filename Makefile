version=1.2.6

all: fmt linux windows macosx

linux: 
	export GOOS=linux && export GOARCH=amd64 && go build ./gxs-ssh-ddl-gen.go && mv gxs-ssh-ddl-gen gxs-ddl-gen_linux64_v${version}

windows:
	export GOOS=windows && export GOARCH=amd64 && go build ./gxs-ssh-ddl-gen.go && mv gxs-ssh-ddl-gen.exe gxs-ddl-gen_win64_v${version}.exe

macosx:
	export GOOS=darwin && export GOARCH=amd64 && go build ./gxs-ssh-ddl-gen.go && mv gxs-ssh-ddl-gen gxs-ddl-gen_macosx_v${version}

fmt:
	@echo "+ $@" ;
	@echo "+ please format Go code with 'gofmt -s'"


clean:
	rm -rf gxs-ddl-gen_* gxs-ddl-gen_*.exe


.PHONY: all linux windows fmt gxs-ddl-gen_v${version} gxs-ddl-gen_v${version}.exe clean

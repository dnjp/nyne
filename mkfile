ALL=nyne nynetab com xcom a+ a- save md move f+ f- font

all:V: ${ALL:%=bin/%}

bin:
	mkdir bin

bin/%:V: ./cmd/% bin
	go build -o bin/$stem cmd/$stem/*.go

install:
	go install ./...

MKSHELL=$PLAN9/bin/rc
uninstall-rc:V:
	for(cmd in $ALL) rm -f $GOPATH/bin/$cmd

MKSHELL=sh
uninstall-sh:V:
	for cmd in $ALL; do
		rm -f $GOPATH/bin/$cmd
	done

uninstall:V: uninstall-sh uninstall-rc

check: $ALL
	go test -count=1 ./...
	go fmt ./...
	go vet ./...
	golint ./...
	staticcheck ./...

clean tidy nuke:V:
	rm -rf ./bin

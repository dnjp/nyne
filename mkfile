ALL=nyne nynetab com xcom a+ a- save md

all:V: $ALL

bin:
	mkdir bin

nyne: bin
	go build -o bin/nyne cmd/nyne/*.go

save: bin
	go build -o bin/save cmd/save/*.go

nynetab: bin
	go build -o bin/nynetab cmd/nynetab/*.go

md: bin
	go build -o bin/md cmd/md/*.go

com: bin
	go build -o bin/com cmd/com/*.go

xcom: bin
	go build -o bin/xcom cmd/xcom/*.go

a+: bin
	go build -o bin/a+ cmd/a+/*.go

a-: bin
	go build -o bin/a- cmd/a-/*.go

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

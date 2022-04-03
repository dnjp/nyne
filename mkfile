ALL=nyne nynetab com xcom a+ a- save md move f+ f- font xec

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

check:
	go test -count=1 ./...
	go fmt ./...
	go vet ./...
	golint ./...
	staticcheck ./...

MKSHELL=$PLAN9/bin/rc
readmes:
	for(f in `{ls cmd}){
		cd $f && \
			goreadme -credit=false -title=`{echo $f | sed 's/cmd\///'} \
			| sed 's/```go/```/g' \
			>README.md && \
			cd ../..
	}

clean tidy nuke:V:
	rm -rf ./bin

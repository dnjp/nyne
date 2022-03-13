
installdir=${installdir}

ALL=nyne nynetab com a+ a-

all:V: $ALL

bin:
	mkdir bin

nyne: bin
	go build -o bin/nyne cmd/nyne/nyne.go

nynetab: bin
	go build -o bin/nynetab cmd/nynetab/nynetab.go

com: bin
	go build -o bin/com cmd/com/com.go

a+: bin
	go build -o bin/a+ cmd/a+/*.go

a-: bin
	go build -o bin/a- cmd/a-/*.go

install:
	cp bin/* $installdir

uninstall:
	rm $installdir/nyne
	rm $installdir/nynetab
	rm $installdir/com
	rm $installdir/a+
	rm $installdir/a-

check: $ALL
	go test -count=1 ./...
	go fmt ./...
	go vet ./...
	golint ./...
	staticcheck ./...

clean tidy nuke:V:
	rm -rf ./bin

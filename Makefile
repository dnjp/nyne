
all: nyne nynetab com ind ui
.PHONY: all


.PHONY: nyne
nyne:
	go build cmd/nyne/nyne.go

.PHONY: nynetab
nynetab:
	go build cmd/nynetab/nynetab.go

.PHONY: com
com:
	go build cmd/com/com.go

.PHONY: ind
ind:
	go build cmd/ind/ind.go

.PHONY: ui
ui:
	go build cmd/ui/ui.go

install:
	go install cmd/nynetab/nynetab.go
	go install cmd/nyne/nyne.go
	
check:
	go test -count=1 ./...
	go fmt ./...
	go vet ./...
	golint ./...
	staticcheck ./...	

clean:
	rm -f nyne nynetab

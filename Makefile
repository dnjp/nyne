
all: nyne nynetab com ind ui gen
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
	
.PHONY: gen
gen:
	mkdir -p ./gen
	go run cmd/gen/gen.go > ./gen/gen.go

install:
	go install cmd/nynetab/nynetab.go
	go install cmd/nyne/nyne.go
	go install cmd/com/com.go
	go install cmd/ind/ind.go
	go install cmd/ui/ui.go
	
check:
	go test -count=1 ./...
	go fmt ./...
	go vet ./...
	golint ./...
	staticcheck ./...	

clean:
	rm -f nyne nynetab

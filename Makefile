
.PHONY: nyne
nyne:
	go build cmd/nyne/nyne.go

.PHONY: nynetab
nynetab:
	go build cmd/nynetab/nynetab.go

clean:
	rm -f nyne nynetab

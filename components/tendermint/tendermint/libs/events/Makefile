.PHONY: docs
REPO:=web3space/ethermint/components/tendermint/tendermint/libs/events

docs:
	@go get github.com/davecheney/godoc2md
	godoc2md $(REPO) > README.md

test:
	go test -v ./...

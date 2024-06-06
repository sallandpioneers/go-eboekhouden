.PHONY: wsdl
wsdl:
	rm -f generated/*.go
	docker run --rm -v ./:/src -w /src sallandpioneers/go-wsdl:0.5.2 gowsdl -o eboekhouden.go -p generated ./templates/eboekhouden.wsdl

.PHONY: lint
lint:
	golangci-lint run --new-from-rev $(git rev-parse origin/master)

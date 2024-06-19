
.PHONY: test
test:
	go test -count=1 bremlin/parser

.PHONE: cover
cover:
	go test -coverprofile=coverage.out bremlin/parser
	go tool cover -html=coverage.out
	rm coverage.out
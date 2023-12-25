# start server
.PHONY: server test cover
server:
	@go run cmd/server/main.go

# run test
test:
	@go test ./... -v

# run test with collecting coverage
cover:
	@go test ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out
	@rm coverage.out
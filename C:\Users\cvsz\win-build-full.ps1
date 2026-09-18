cd C:\Users\cvsz\zneondrive\services\game-api
go build -trimpath ./cmd/server
echo BUILD_OK
go vet ./...
echo VET_OK
go test ./...
echo TEST_OK

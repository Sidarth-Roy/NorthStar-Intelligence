go work init ./Backend ./DB
go mod download
go clean -modcache
go mod tidy
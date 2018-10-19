cover:
	go tool cover -html=c.out

unittest:
	go test -timeout=30s -coverpkg=./... -coverprofile c.out github.com/trungnn/sialia/test


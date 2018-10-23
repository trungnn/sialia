test:
	go test -v -timeout=30s -coverpkg=./promise -coverprofile c.out github.com/trungnn/sialia/tests

cover:
	go tool cover -html=c.out


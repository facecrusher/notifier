test:
	go test ./... -coverpkg=./... -race;

cover:
	go test -coverprofile="/tmp/go-cover.$$.tmp" -coverpkg=./... ./... -race && go tool cover -html="/tmp/go-cover.$$.tmp" && unlink "/tmp/go-cover.$$.tmp";
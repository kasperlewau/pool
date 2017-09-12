test:
	go test -v -race
cover:
	go test -coverprofile=cov.out && go tool cover -html=cov.out
bench:
	go test -run= -bench=. -benchmem

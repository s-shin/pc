pkg_base := github.com/s-shin/pc

test_pc:
	go test -v -cover $(pkg_base)

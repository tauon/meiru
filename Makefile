meiru: meiru.go
	go build meiru.go

test: meiru
	createdb meiru_test 2> /dev/null || true
	go test

clean:
	rm -f meiru
	dropdb meiru_test 2> /dev/null || true

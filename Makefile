all:
	go build

run:
	./weatherdata ksql

fmt:
	go fmt

packages:
	#go get "github.com/cfdrake/go-ystocks"
	go get code.google.com/p/go.net/html

clean:
	go clean

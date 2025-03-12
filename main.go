package main

type TransportMethod interface {
	// test
	DeliverPackage(destination string) (string, error)
	GetStatus() string
}

func main() {

}

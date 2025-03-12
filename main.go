package main

type TransportMethod interface {
	DeliverPackage(destination string) (string, error)
	GetStatus() string
}

func main() {

}

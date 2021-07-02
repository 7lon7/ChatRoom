package main

func main() {
	svr := NewServer("127.0.0.1",11111)
	svr.Start()
}
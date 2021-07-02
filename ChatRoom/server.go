package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip string
	Port int
	Msg chan string
	OnlineMap map[string]*User
	MapMutex sync.RWMutex
}

func NewServer(ip string,port int) *Server {
	server := &Server{
		Ip: ip,
		Port: port,
		Msg: make(chan string),
		OnlineMap: make(map[string]*User),
	}
	return server
}

func (i *Server)AssignMsg()  {
	for  {
		msg := <- i.Msg
		i.MapMutex.Lock()
		for _,cli := range i.OnlineMap {
			cli.Chnl <- msg
		}
		i.MapMutex.Unlock()
	}
}

func (i *Server)Handler(conn net.Conn)  {
	user := NewUser(conn, i)
	user.Online()
	go func() {
		rcv := make([]byte,4096)
		for  {
			n,err := conn.Read(rcv)
			if n == 0 {
				user.Outline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Handler : conn.Read error : ",err)
				return
			}
			user.DoMsg(string(rcv[0:n-1]))
		}
	}()
	select {

	}
}

func (i *Server)BroadCast(user *User,str string) {
	msg := "["+user.Name+"]:"+str
	i.Msg <- msg
}

func (i *Server)Start()  {
	listener,err := net.Listen("tcp",fmt.Sprintf("%s:%d", i.Ip, i.Port))
	if err != nil {
		fmt.Println("net.Listen error : ",err)
		return
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			fmt.Println("listen.Close() error : ",err)
		}
	}(listener)
	go i.AssignMsg()
	for{
		conn,err := listener.Accept()
		if err != nil {
			continue
		}
		go i.Handler(conn)
	}
}
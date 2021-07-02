package main

import (
	"fmt"
	"io"
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	Chnl chan string
	conn net.Conn
	server *Server
}

func NewUser(conn net.Conn,server *Server) *User {
	user := &User{
		Name: conn.RemoteAddr().String(),
		Addr: conn.RemoteAddr().String(),
		Chnl: make(chan string),
		conn: conn,
		server: server,
	}
	go user.ListenAndSendMsg()
	return user
}

func (i *User)Online() {
	i.server.MapMutex.Lock()
	i.server.OnlineMap[i.Name] = i
	i.server.MapMutex.Unlock()
	i.server.BroadCast(i,"Online!")
}

func (i *User)Outline() {
	i.server.MapMutex.Lock()
	delete(i.server.OnlineMap, i.Name)
	i.server.MapMutex.Unlock()
	i.server.BroadCast(i,"Outline!")
}

func (i *User)DoMsg(msg string) {
	index := strings.Index(msg,"|")
	if index > 0 {
		cmd := msg[:index]
		if cmd == "to" {
			lastindex := strings.LastIndex(msg,"|")
			if lastindex > 0 && lastindex - index > 1  {
				name := msg[index+1:lastindex]
				realmsg := msg[lastindex+1:]+"\n"
				user,ok := i.server.OnlineMap[name]
				if ok {
					_,err := user.conn.Write([]byte("solo msg from "+i.Name+":"+realmsg))
					_,err = i.conn.Write([]byte("solo msg to "+name+":"+realmsg))
					if err != nil {
						fmt.Println(err)
					}
					return
				} else {
					_,err := i.conn.Write([]byte("unknown user\n"))
					if err != nil {
						fmt.Println(err)
					}
					return
				}
			}

		} else if cmd == "rename" {
			newname := msg[index+1:]
			i.server.BroadCast(i,"rename->["+newname+"]")
			i.server.MapMutex.Lock()
			delete(i.server.OnlineMap, i.Name)
			i.server.OnlineMap[newname] = i
			i.server.MapMutex.Unlock()
			i.Name = newname
			return
		} else if cmd == "list" {
			i.server.MapMutex.Lock()
			for key,_ := range i.server.OnlineMap {
				_,err := i.conn.Write([]byte(key+"\n"))
				if err != nil {
					fmt.Println(err)
				}
			}
			i.server.MapMutex.Unlock()
			return
		}
	}
	i.server.BroadCast(i,msg)
}


func (i *User)ListenAndSendMsg() {
	for {
		msg := <-i.Chnl
		_,err := i.conn.Write([]byte(msg + "\n"))
		if err != nil {
			if err == io.ErrClosedPipe {
				fmt.Println(err)
				return
			} else {
				continue
			}
		}
	}
}

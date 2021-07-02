package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIP string
	ServerPort int
	Name string
	Conn net.Conn
	Flag int
}

func NewClient(ip string,port int) *Client {
	client := &Client{
		ServerIP: ip,
		ServerPort: port,
		Flag: 999,
	}
	conn,err := net.Dial("tcp",fmt.Sprintf("%s:%d",ip,port))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	client.Conn = conn
	return client
}

func (i *Client)Recv() {
	_,err := io.Copy(os.Stdout,i.Conn)
	if err != nil {
		panic(err)
	}
}

func (i *Client)SelectFunc() bool {
	var f int
	fmt.Println("1.群聊")
	fmt.Println("2.私聊")
	fmt.Println("3.改名")
	fmt.Println("0.退出")
	_,err := fmt.Scanln(&f)
	if err != nil {
		fmt.Println("SelectFunc() error ",err)
		return false
	}
	if f >= 0 && f <= 3 {
		i.Flag = f
		return true
	}
	return false
}

func (i *Client)GroupChat() {
	var msg string
	fmt.Println("input chat msg(exit)")
	_,err := fmt.Scanln(&msg)
	if err != nil {
		fmt.Println("GroupChat()/Scanf() error ",err)
		return
	}
	for msg != "exit" {
		if len(msg) != 0 {
			_,err = i.Conn.Write([]byte(msg+"\n"))
			if err != nil {
				fmt.Println("GroupChar()/Conn.Write() error ",err)
				break
			}
		}
		msg = ""
		fmt.Println("input chat msg(exit)")
		_,err := fmt.Scanln(&msg)
		if err != nil {
			fmt.Println("GroupChat()/Scanf() error ",err)
			return
		}
	}
}

func (i *Client)SelectUser()  {
	_,err := i.Conn.Write([]byte("list|\n"))
	if err != nil {
		fmt.Println("SoloChat()/Conn.Write() error ",err)
	}
}

func (i *Client)SoloChat() {
	var who string
	var msg string
	i.SelectUser()
	fmt.Println("input who you wanna chat(exit)")
	_,err := fmt.Scanln(&who)
	if err != nil {
		fmt.Println("SoloChat()/Scanln() error ",err)
		return
	}
	for who != "exit" {
		fmt.Println("input msg to "+who+"(exit)")
		_,err = fmt.Scanln(&msg)
		if err != nil {
			fmt.Println("SoloChat()/Scanln() error ",err)
			return
		}
		for msg != "exit" {
			if len(msg) != 0 {
				_,err = i.Conn.Write([]byte("to|"+who+"|"+msg+"\n"))
				if err != nil {
					fmt.Println("GroupChat()/Conn.Write() error ",err)
					break
				}
			}
			msg = ""
			fmt.Println("input chat msg(exit)")
			_,err := fmt.Scanln(&msg)
			if err != nil {
				fmt.Println("GroupChat()/Scanf() error ",err)
				return
			}
		}
		i.SelectUser()
		fmt.Println("input who you wanna chat(exit)")
		_,err := fmt.Scanln(&who)
		if err != nil {
			fmt.Println("SoloChat()/Scanln() error ",err)
			return
		}
	}
}

func (i *Client)UpdateName() bool {
	var newname string
	fmt.Println("input new name:")
	_,err := fmt.Scanln(&newname)
	if err != nil {
		fmt.Println("UpdateName()/Scanf() error ",err)
		return false
	}
	i.Name = newname
	_,err = i.Conn.Write([]byte("rename|"+newname+"\n"))
	if err != nil {
		fmt.Println("UpdateName()/Conn.Write() error ",err)
		return false
	}
	return true
}

func (i *Client)Run() {
	for i.Flag != 0 {
		for i.SelectFunc() != true {

		}
		switch i.Flag {
		case 1:
			i.GroupChat()
			break
		case 2:
			i.SoloChat()
			break
		case 3:
			i.UpdateName()
			break
		case 0:
			os.Exit(0)
		}
	}
}


var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp,"ip","127.0.0.1","设置服务器IP(默认 127.0.0.1)")
	flag.IntVar(&serverPort,"port",11111,"设置服务器端口(默认 11111)")
}

func main() {
	flag.Parse()
	client := NewClient(serverIp,serverPort)
	if client == nil {
		fmt.Println("NewClient() error")
		return
	}
	go client.Recv()
	client.Run()
	select {

	}
}

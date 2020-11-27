package main

import (
	"bufio"
	"flag"
	"log"
	"net"
	"reflect"
	"strings"
)

var flagIP = flag.String("ip", "127.0.0.1", "IP address to listen on")
var flagPort = flag.String("port", "7300", "Port to listen on")

type client struct {
	Conn     net.Conn
	CallSign string
	Message  chan string
}

var connections map[string]*client

func main() {
	flag.Parse()
	connections = make(map[string]*client, 0)

	//start listener
	listener, err := net.Listen("tcp", *flagIP+":"+*flagPort)
	if err != nil {
		log.Fatalf("could not listen on interface %v:%v error: %v ", *flagIP, *flagPort, err)
	}
	defer listener.Close()
	log.Println("listening on: ", listener.Addr())
	//main listen accept loop
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("could not accept connection %v ", err)
		}
		//create new client on connection
		go createclient(conn)

		for _, c := range connections {
			c.Message <- "someone has joined!"
		}
	}
}

func (c *client) close() {
	c.Conn.Close()
	c.Message <- "\\quit"
	delete(connections, c.CallSign)
}

func (c *client) recieve() {
	for {
		msg := <-c.Message
		log.Printf("recieve: client(%v|%v) recvd msg: %s ", c.Conn.RemoteAddr(), c.CallSign, msg)
		writeFormattedMsg(c.Conn, msg)
	}
}

func createclient(conn net.Conn) {

	log.Printf("createclient: remote connection from: %v", conn.RemoteAddr())

	callSign, err := readInput(conn, "Login: ")
	if err != nil {
		panic(err)
	}

	writeFormattedMsg(conn, "Welcome "+callSign)

	//init client struct
	client := &client{
		Message:  make(chan string),
		Conn:     conn,
		CallSign: callSign,
	}

	log.Printf("new client created: %v %v", client.Conn.RemoteAddr(), client.CallSign)

	connections[client.CallSign] = client

	//spin off seperate send, recieve
	go client.recieve()
}

func writeFormattedMsg(conn net.Conn, msg interface{}) error {
	_, err := conn.Write([]byte("---------------------------\n"))
	t := reflect.ValueOf(msg)
	switch t.Kind() {
	case reflect.Map:
		for k, v := range msg.(map[string]string) {
			_, err = conn.Write([]byte(k + " : " + v))
		}
		break
	case reflect.String:
		v := reflect.ValueOf(msg).String()
		_, err = conn.Write([]byte(v + "\n"))
		break
	} //switch
	conn.Write([]byte("---------------------------\n"))

	if err != nil {
		return err
	}
	return nil //todo
}

func readInput(conn net.Conn, qst string) (string, error) {
	conn.Write([]byte(qst))
	s, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Printf("readinput: could not read input from stdin: %v from client %v", err, conn.RemoteAddr().String())
		return "", err
	}
	s = strings.Trim(s, "\r\n")
	return s, nil
}

package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"reflect"
	"sort"
	"strings"
	"time"
)

var flagIP = flag.String("ip", "127.0.0.1", "IP address to listen on")
var flagPort = flag.String("port", "7300", "Port to listen on")

type client struct {
	Conn    net.Conn
	User    user
	Message chan string
}

var (
	connections map[string]*client
	spots       Spots
	maxSpotID   int
)

func blast(msg string) {
	for _, c := range connections {
		writeFormattedMsg(c.Conn, msg)
	}
}

func main() {
	flag.Parse()
	connections = make(map[string]*client, 0)

	go func() {
		for {
			log.Println("fetching spots")
			err := getSpots(&spots)
			if err != nil {
				fmt.Println(err)
			}
			sort.SliceStable(spots, func(i, j int) bool {
				return spots[i].SpotID < spots[j].SpotID
			})

			if len(spots) > 1 {

				if spots[0].SpotID > maxSpotID {
					for _, s := range spots {
						if s.SpotID > maxSpotID {
							blast(s.ToClusterFormat())
						}
					}
				}

				maxSpotID = spots[0].SpotID
			}
			time.Sleep(30 * time.Second)
		}
	}()

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
		time.Sleep(time.Millisecond)
	}
}

func (c *client) close() {
	c.Conn.Close()
	c.Message <- "\\quit"
	delete(connections, c.id())
}

func (c *client) id() string {
	return c.Conn.RemoteAddr().String()
}

func (c *client) receive() {
	defer c.close()

	for {
		userInput, err := bufio.NewReader(c.Conn).ReadString('\n')
		if err != nil {
			log.Printf("receive: client(%v|%v) recvd msg with error: %s ", c.Conn.RemoteAddr(), c.User.CallSign, err.Error())
			return
		}
		userInput = strings.TrimSuffix(userInput, "\n")
		log.Printf("receive: client(%v|%v) recvd msg: %s ", c.Conn.RemoteAddr(), c.User.CallSign, userInput)

		time.Sleep(time.Millisecond)
	}
}

func getSpots(target interface{}) error {

	var myClient = &http.Client{Timeout: 10 * time.Second}

	r, err := myClient.Get("https://api.pota.app/spot/activator")
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func createclient(conn net.Conn) {

	log.Printf("createclient: remote connection from: %v", conn.RemoteAddr())

	callSign, err := readInput(conn, "Login: ")
	if err != nil {
		log.Printf("%s", err)
	}

	writeFormattedMsg(conn, "Welcome "+callSign)

	user := &user{
		CallSign: callSign,
	}
	//init client struct
	client := &client{
		Message: make(chan string),
		Conn:    conn,
		User:    *user,
	}

	log.Printf("new client created: %v %v", client.Conn.RemoteAddr(), client.User.CallSign)

	connections[client.id()] = client

	for _, s := range spots {
		writeFormattedMsg(conn, s.ToClusterFormat())
	}

	//spin off seperate send, receive
	go client.receive()
}

func writeFormattedMsg(conn net.Conn, msg interface{}) error {
	var err error
	t := reflect.ValueOf(msg)
	switch t.Kind() {
	case reflect.Map:
		for k, v := range msg.(map[string]string) {
			_, err = conn.Write([]byte(k + " : " + v))
		}
		break
	case reflect.String:
		v := reflect.ValueOf(msg).String()
		if !strings.HasSuffix(v, "\n") {
			v += "\n"
		}
		_, err = conn.Write([]byte(v))
		break
	} //switch

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

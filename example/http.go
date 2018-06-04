package main

import (
	"flag"
	"github.com/libp2p/go-libp2p"
	"fmt"
	"context"

	"github.com/upperwal/go-libp2p-protocol"
	"github.com/libp2p/go-libp2p-net"
	"strings"
	"errors"
	"log"
)

/*
 This is a stripped down version of HTTP
*/

const(
	GET = "GET"
	POST = "POST"

	STATUS_OK = "200"
	STATUS_NOT_FOUND = "400"

	HTTP_VERSION = "HTTP/1.0"
)

var protoHTTP = protocol.Protocol{
	Name: "ping",
	Version: 1,
	Run: func(stream net.Stream, init bool) {

		fmt.Println("Stream opened to ", stream.Conn().RemotePeer())

		if init == true {
			fmt.Println("Client Started")
			go client(stream)
		} else {
			fmt.Println("Server Started")
			go server(stream)
		}
	},
}

func server(stream net.Stream) {
	var data [100]byte
	for {
		n, err := stream.Read(data[:])
		if err != nil {
			log.Println(err)
			return
		}

		method, path, err := parseMethod(data[:n])

		sendData, err := fetchDataToSend(method, path)

		if err == nil {
			stream.Write([]byte(HTTP_VERSION + " " + STATUS_OK + " OK\n" +  sendData))
		} else {
			stream.Write([]byte(HTTP_VERSION + " " + STATUS_NOT_FOUND + " Not_Found\nsome"))
		}

	}

}

func client(stream net.Stream) {
	var data [100]byte
	//for {
		_, err := stream.Write([]byte(GET + " / " + HTTP_VERSION))
		if err != nil {
			panic(err)
		}

		n, err := stream.Read(data[:])
		if err != nil {
			panic(err)
		}

		res, err := parseResponse(data[:n])

		if err == nil {
			fmt.Println("Data: ", res)
		} else {
			panic(err)
		}

	//}
}

func parseResponse(data []byte) (string, error) {
	segrigateHeader := strings.Split(string(data), "\n")

	if len(segrigateHeader) != 2 {
		return "", errors.New("One of header or body received")
	}
	headerItems := strings.Split(segrigateHeader[0], " ")
	if len(headerItems) != 3 {
		return "", errors.New("Header not complete")
	} else if headerItems[0] != HTTP_VERSION {
		return "", errors.New("This version of HTTP not supported")
	} else if headerItems[1] != STATUS_OK {
		return "", errors.New("HTTP status not OK")
	}

	return segrigateHeader[1], nil
}

func fetchDataToSend(method, path string) (string, error) {
	switch path {

		case "/":
			return "Hello", nil
			break

	}

	return "some", errors.New("404 not found")
}

func parseMethod(bytedata []byte) (string, string, error) {

	splitString := strings.Split(string(bytedata), " ")

	if len(splitString) != 3 {
		return "", "", errors.New("No HTTP protocol")
	} else if splitString[2] != "HTTP/1.0" {
		return "", "", errors.New("This version of HTTP not supported")
	}

	return splitString[0], splitString[1], nil
}

func main() {
	dest := flag.String("d", "", "Dest MultiAddr String")
	flag.Parse()

	host, err := libp2p.New(context.Background(), libp2p.Defaults)
	if(err != nil) {
		panic(err)
	}

	xHost := protocol.SetProtocol(host, []protocol.Protocol{protoHTTP})
	fmt.Printf("This peer address: %s\n", xHost.GetCompleteAddr())

	if *dest != "" {
		err := xHost.AddPeerWithAddr(*dest)
		if err != nil {
			panic(err)
		}
	}
	select{}
}

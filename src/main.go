package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"parseData"
	"strconv"
)

func main() {
	parseData.Printmessage()
	service := ":5432"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	CheckError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	CheckError(err)
	for {
		conn, err := listener.Accept()
		//fmt.Println("success connect")
		if err != nil {
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	var buf []byte = make([]byte, 64, 64)
	var newBuff []byte
	var data []byte
	var dataLength int
	var lengthEnd int
	//var data_list [] string
	for {
		recvLen, err := conn.Read(buf)
		CheckError(err)
		fmt.Println("buf", buf)
		buf = bytes.TrimRight(buf, "\x00")
		newBuff = append(newBuff, buf[0:recvLen]...)
		if bytes.Contains(newBuff, []byte("\x0A")) == false { // if we receive the length
			continue
		} else {
			lengthEnd = bytes.Index((newBuff), []byte("\x0A")) // locate the length
			dataLength, _ = strconv.Atoi(string(newBuff[0:lengthEnd]))
			for len(newBuff)-len(newBuff[0:lengthEnd]) >= dataLength {
				data = append(data, newBuff[lengthEnd+1:lengthEnd+1+dataLength]...)
				newBuff = newBuff[lengthEnd+1+dataLength:] //new packet
				//data_list.append(data_list, data)
				//if strings.Contains(string(newBuff),"\n") == true {
				//	lengthEnd = strings.Index(string(newBuff), "\n")
				//	dataLength = strconv.Atoi(string(newbuff[0, lengthEnd]))
				//	continue;
				//} else {
				//	break
				//}
				if len(data) == dataLength { //receive finished
					//fmt.Println("data=", data)
					break
				}
			}
		}
		if len(data) == dataLength {
			break
		}
	}
	fmt.Println("data=", string(data[:]))
}

func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

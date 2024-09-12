// Adaptado de Alan A. A. Donovan & Brian W. Kernighan.
// a TCP server that periodically writes the time.
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"
	"bufio"
)

type IpsConfigs struct{
	Ips      []string `json:"ips"`
}

func main() {

	//escuta na porta 8000 (pode ser monitorado com lsof -Pn -i4 | grep 8000)
	listener, err := net.Listen("tcp", "150.165.42.171:8000")

	jsonFile, err := os.Open(`ips.json`)

	byteValueJSON, _:= ioutil.ReadAll(jsonFile)

	//Declaração abreviada de um objeto do tipo Book
	objIps := IpsConfigs{}

	//Conversão da variável byte em um objeto do tipo struct Book
	json.Unmarshal(byteValueJSON, &objIps)

	fmt.Println(objIps.Ips)

	if err != nil {
		log.Fatal(err)
	}
	for {
		//aceita uma conexão criada por um cliente
		conn, err := listener.Accept()
		if err != nil {
			// falhas na conexão. p.ex abortamento
			log.Print(err)
			continue
		}
		// serve a conexão estabelecida
		go handleConn(conn)
		
	}
}

func handleConn(c net.Conn) {

	defer c.Close()
	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		fmt.Println(netData)

		if err != nil {
			return
		}
		time.Sleep(1 * time.Second)
	}
}
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

type IpsConfigs struct {
	Ips []string `json:"ips"`
}

type HashResponse struct {
	Hash int
	IP   string
}

var hashs []HashResponse
var ipsConfigs IpsConfigs

func main() {
	loadIpsConfigs()

	listener, err := net.Listen("tcp", "150.165.42.171:8001")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func loadIpsConfigs() {
	jsonFile, err := os.Open("ips.json")
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	byteValueJSON, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValueJSON, &ipsConfigs)
	if err != nil {
		log.Fatal(err)
	}
}

func handleConn(c net.Conn) {

	reader := bufio.NewReader(c)
	
	for {
		netData, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("Connection closed by client")
				return
			}
			log.Println("Error reading from connection:", err)
			return
		}
		
		netData = strings.TrimSpace(netData)
		parts := strings.SplitN(netData, " ", 2)

		if len(parts) < 2 {
			fmt.Fprintln(c, "Invalid command format")
			continue
		}

		switch parts[0] {
		case "search":
			hash, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Fprintln(c, "Invalid hash format")
				continue
			}
			search(c, hash)
		case "update":
			hashParts := strings.SplitN(parts[1], " ", 2)
			if len(hashParts) != 2 {
				fmt.Fprintln(c, "Invalid update format. Expected format: <hash> <ip>")
				continue
			}
			hash, err := strconv.Atoi(hashParts[0])
			if err != nil {
				fmt.Fprintln(c, "Invalid hash format")
				continue
			}
			ip := hashParts[1]
			appendHash(HashResponse{Hash: hash, IP: ip})
			fmt.Fprintln(c, "Hash updated successfully")
		default:
			fmt.Fprintln(c, "Unknown command")
		}
	}
}

func search(conn net.Conn, hash int) {
	var response []string
	for _, hashResponse := range hashs {
		if hashResponse.Hash == hash {
			response = append(response, hashResponse.IP)
		}
	}
	if len(response) == 0 {
		fmt.Fprintln(conn, "No matches found")
	} else {
		fmt.Fprintln(conn, strings.Join(response, ", "))
	}
}

func appendHash(response HashResponse) {
	fmt.Println(response)
	hashs = append(hashs, response)
}

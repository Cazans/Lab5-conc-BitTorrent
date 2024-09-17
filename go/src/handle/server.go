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

var hashs = make(map[int][]string)
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
	defer c.Close()

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
			appendHash(hash, ip)
			fmt.Fprintln(c, "Hash updated successfully")
		default:
			fmt.Fprintln(c, "Unknown command")
		}
	}
}

func search(conn net.Conn, hash int) {
	ips, found := hashs[hash]
	if !found {
		fmt.Fprintln(conn, "No matches found")
		return
	}
	fmt.Fprintln(conn, strings.Join(ips, ", "))
}

func appendHash(hash int, ip string) {
	hashs[hash] = append(hashs[hash], ip)
	fmt.Println("Updated hashs map:", hashs)
}

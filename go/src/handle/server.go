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

	listener, err := net.Listen("tcp", "127.0.0.1:8001")
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
		fmt.Println("Received command:", netData)
		parts := strings.SplitN(netData, " ", 2)

		if len(parts) < 2 {
			fmt.Println("Invalid command format")
			continue
		}

		switch parts[0] {
		case "search":
			hash, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Invalid hash format")
				continue
			}
			search(c, hash)
		case "update":
			stringIp := c.RemoteAddr().String()
			ipClient := strings.Split(stringIp, ":")[0]
			hashParts := parts[1]
			hash, err := strconv.Atoi(hashParts)
			if err != nil {
				fmt.Println("Invalid hash format")
				continue
			}
			appendHash(hash, ipClient)
			fmt.Println("Hash updated successfully")
		default:
			fmt.Println("Unknown command")
		}
	}
}

func search(conn net.Conn, hash int) {
	ips, found := hashs[hash]
	if !found || len(ips) == 0 {
		fmt.Fprintln(conn, "No matches found")
		return
	}

	// Retornar os IPs associados ao hash, cada um em uma nova linha
	for _, ip := range ips {
		fmt.Fprintln(conn, ip)
	}
}

func appendHash(hash int, ip string) {
	// Verifica se o IP já está presente para o hash
	ipSet := make(map[string]struct{})
	for _, existingIP := range hashs[hash] {
		ipSet[existingIP] = struct{}{}
	}

	if _, exists := ipSet[ip]; !exists {
		hashs[hash] = append(hashs[hash], ip)
		fmt.Println("Updated hashs map:", hashs)
	} else {
		fmt.Printf("IP %s already exists for hash %d\n", ip, hash)
	}
}

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

var hashs = make(map[int][]string)

func main() {

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

func handleConn(c net.Conn) {
	defer c.Close()

	reader := bufio.NewReader(c)

	for {
		netData, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("Conexão fechada pelo cliente")
				return
			}
			log.Println("Erro ao ler da conexão:", err)
			return
		}

		netData = strings.TrimSpace(netData)
		fmt.Println("Comando recebido:", netData)
		parts := strings.SplitN(netData, " ", 2)

		if len(parts) < 2 && parts[0] != "update" {
			fmt.Fprintln(c, "Formato de comando inválido")
			continue
		}

		switch parts[0] {
		case "search":
			hash, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Formato de hash inválido")
				continue
			}
			search(c, hash)

		case "update":
			// Imprime o mapa antes de atualizar
			fmt.Println("Mapa de hashs ANTES da atualização:", hashs)

			stringIp := c.RemoteAddr().String()
			ipClient := strings.Split(stringIp, ":")[0]

			hash, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Fprintln(c, "Formato de hash inválido")
				return
			}
			appendHash(hash, ipClient)

			// Imprime o mapa depois de atualizar
			fmt.Println("Mapa de hashs DEPOIS da atualização:", hashs)

		default:
			fmt.Fprintln(c, "Comando desconhecido")
		}
	}
}

func search(conn net.Conn, hash int) {
	ips, found := hashs[hash]
	if !found || len(ips) == 0 {
		fmt.Fprintln(conn, "Nenhuma correspondência encontrada")
		return
	}

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
		fmt.Println("IP adicionado ao hash", hash)
	} else {
		fmt.Printf("IP %s já existe para o hash %d\n", ip, hash)
	}
}

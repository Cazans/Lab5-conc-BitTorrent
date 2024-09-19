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

	listener, err := net.Listen("tcp", "150.165.74.99:8001")
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
			stringIp := c.RemoteAddr().String()
			ipClient := strings.Split(stringIp, ":")[0]
			parts[1] = strings.Trim(parts[1], "[]")
			strSlice := strings.Fields(parts[1])
			intSlice := make([]int, 0, len(strSlice))
			for _, s := range strSlice {
				num, err := strconv.Atoi(s)
				if err != nil {
					fmt.Println("Erro ao converter string para inteiro:", err)
					continue
				}
				intSlice = append(intSlice, num)
			}
			handleUpdate(intSlice, ipClient)
			fmt.Fprintln(c, "Update realizado com sucesso!")
		default:
			fmt.Fprintln(c, "Comando desconhecido")
		}
	}
}

func handleUpdate(intSlice []int, ip string) {
	ipSet := make(map[string]struct{})
	for _, hash := range intSlice {
		for _, existingIP := range hashs[hash] {
			ipSet[existingIP] = struct{}{}
		}

		if _, exists := hashs[hash]; !exists && ipSet[ip] == struct{}{} {
			hashs[hash] = append(hashs[hash], ip)
			fmt.Println("IP adicionado ao hash", hash)
		} else if _, exists := ipSet[ip]; !exists{
			hashs[hash] = append(hashs[hash], ip)
			fmt.Println("IP adicionado ao hash", hash)
		} else {
			fmt.Printf("IP %s já existe para o hash %d\n", ip, hash)
		}
	}

	var diferenca []int

	for hash := range hashs {
		if !contem(intSlice, hash) {
			diferenca = append(diferenca, hash)
		}
	}

	for _, hash := range diferenca {
		hashs[hash] = removerValor(hashs[hash], ip)
		if len(hashs[hash]) == 0 {
			delete(hashs, hash)
		}
	}
	fmt.Println(hashs)
}

func contem[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

func removerValor(slice []string, valor string) []string {
	for i, v := range slice {
		if v == valor {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func search(conn net.Conn, hash int) {
	ips, found := hashs[hash]
	if !found || len(ips) == 0 {
		fmt.Fprintln(conn, "Nenhuma correspondência encontrada")
		return
	}
	fmt.Fprintln(conn, ips)
}

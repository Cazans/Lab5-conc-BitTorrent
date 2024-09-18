package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"lab5-conc/hashFiles" // Importando o pacote para calcular o hash dos arquivos
)

func main() {

	conn, err := net.Dial("tcp", "localhost:8001")
	if err != nil {
		fmt.Println("Erro ao conectar-se ao servidor:", err)
		return
	}
	defer conn.Close()

	for {
		fmt.Print("Digite o comando (search <hash> ou update): ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		// Verifica se o comando começa com "update"
		if strings.HasPrefix(input, "update") {
			// Se for "update" sem nada mais
			if input == "update" {
				hashFiles.SendHash(conn, "localhost") // Calcula e envia os hashes
			} else {
				fmt.Println("Resposta do servidor: Formato de comando inválido")
			}
			continue
		}

		// Envia o comando para o servidor
		fmt.Fprintf(conn, input+"\n")

		// Lê a resposta do servidor
		serverResponse, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Println("Resposta do servidor: ", serverResponse)
	}
}

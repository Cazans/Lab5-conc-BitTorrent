package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// Conectar ao servidor
	conn, err := net.Dial("tcp", "150.165.42.171:8000")
	if err != nil {
		fmt.Println("Erro ao conectar-se ao servidor:", err)
		return
	}
	defer conn.Close()

	// Ler entrada do usuário (hash para busca)
	fmt.Print("Digite o comando de busca (ex: search <hash>): ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')

	// Enviar comando para o servidor
	fmt.Fprintf(conn, input)

	// Receber e imprimir a resposta do servidor
	serverResponse, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Println("Resposta do servidor:", serverResponse)
}

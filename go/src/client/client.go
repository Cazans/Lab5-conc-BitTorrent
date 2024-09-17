package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"

	"lab5-conc/hashfiles" // Ajuste o caminho do módulo conforme necessário
)

func main() {
	// Conectar ao servidor
	conn, err := net.Dial("tcp", "150.165.42.171:8000")
	if err != nil {
		fmt.Println("Erro ao conectar-se ao servidor:", err)
		return
	}
	defer conn.Close()

	var wg sync.WaitGroup

	// Iniciar goroutine para enviar hashes
	wg.Add(1)
	go func() {
		defer wg.Done()
		hashfiles.SendHash(conn, "clientIP")
	}()

	// Ler entrada do usuário (hash para busca)
	fmt.Print("Digite o comando de busca (ex: search <hash>): ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	fmt.Fprintf(conn, input)

	serverResponse, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Println("Resposta do servidor:", serverResponse)

	// Aguarda a goroutine terminar antes de encerrar o programa
	wg.Wait()
}

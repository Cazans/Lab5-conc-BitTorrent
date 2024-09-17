package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"

	"lab5-conc/hashFiles"
)

func main() {
	
	conn, err := net.Dial("tcp", "localhost:8001")
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
		hashFiles.SendHash(conn, "localhost")
	}()

	for{
		fmt.Print("Digite o comando de busca (ex: search <hash>): ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		fmt.Fprintf(conn, input)

		serverResponse, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Println("Resposta do servidor:", serverResponse)
		wg.Wait()
	}
	
	
}

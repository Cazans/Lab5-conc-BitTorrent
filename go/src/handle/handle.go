// Adaptado de Alan A. A. Donovan & Brian W. Kernighan.
// a TCP server that periodically writes the time.
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
	"time"
)

type IpsConfigs struct {
	Ips []string `json:"ips"`
}

func main() {

	//escuta na porta 8000 (pode ser monitorado com lsof -Pn -i4 | grep 8000)
	listener, err := net.Listen("tcp", "0.0.0.0:8000")

	jsonFile, err := os.Open(`ips.json`)

	byteValueJSON, _ := ioutil.ReadAll(jsonFile)

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
func search(hash int) {

}

func handleConn(c net.Conn) {

	defer c.Close()
	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		partes := strings.SplitN(netData, " ", 2)
		if partes[0] == "search" {
			num, err := strconv.Atoi(partes[1])
			if err != nil {
				files := sendHash()
				fmt.Println(files)
				search(num)
			}
		}
		fmt.Println(netData)

		if err != nil {
			return
		}
		time.Sleep(1 * time.Second)
	}
}

// HashFiles
// Estrutura que armazena o hash e a data de modificação de um arquivo
type fileInfo struct {
	hash         int
	lastModified string
}

// Lê um arquivo e retorna o conteúdo em bytes e a data de modificação
func readFile(filePath string) ([]byte, string, error) {
	data, readErr := os.ReadFile(filePath)
	if readErr != nil {
		fmt.Printf("Error reading file %s: %v", filePath, readErr)
		return nil, " ", readErr
	}

	fileInfo, infoErr := os.Stat(filePath)

	if infoErr != nil {
		fmt.Printf("Error reading file %s: %v", filePath, infoErr)
		return nil, " ", infoErr
	}

	lastModified := fileInfo.ModTime().String()

	return data, lastModified, nil
}

// Retorna o valor refetente ao hash de um arquivo
func fileToHash(filePath string, hashes chan fileInfo) (int, string, error) {
	data, lastModified, err := readFile(filePath)

	if err != nil {
		return 0, " ", err
	}

	hash := 0

	for _, _byte := range data {
		hash += int(_byte)
	}

	hashes <- fileInfo{hash, lastModified}

	return hash, lastModified, nil
}

// Envia o hash de todos os arquivos de um diretório para o servidor
func sendHash() fileInfo {
	dirPath := "/tmp/dataset"

	files, err := os.ReadDir(dirPath)

	if err != nil {
		fmt.Printf("Error reading directory %s: %v", dirPath, err)
	}

	hashes := make(chan fileInfo, len(files))

	for _, file := range files {
		filePath := dirPath + "/" + file.Name()
		go fileToHash(filePath, hashes)
	}
	var processedFile fileInfo
	for range files {
		processedFile := <-hashes
		//enviar a variável processedFile para o servidor
		fmt.Println(processedFile)
	}
	return processedFile
}

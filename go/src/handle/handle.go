// Adaptado de Alan A. A. Donovan & Brian W. Kernighan.
// Um servidor TCP que escuta por requisições de busca de arquivos e retorna IPs das máquinas que possuem o arquivo solicitado.

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

var ipsConfigs IpsConfigs

func main() {
	// Carregar IPs das máquinas
	loadIpsConfigs()

	// Escuta na porta 8000
	listener, err := net.Listen("tcp", "150.165.74.47:8000")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		// Aceita uma conexão criada por um cliente
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		// Serve a conexão estabelecida
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
	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			log.Print(err)
			return
		}

		netData = strings.TrimSpace(netData)
		partes := strings.SplitN(netData, " ", 2)
		if len(partes) < 2 {
			continue
		}

		if partes[0] == "search" {
			hash, err := strconv.Atoi(partes[1])
			if err != nil {
				fmt.Fprintln(c, "Invalid hash")
				continue
			}
			// Realiza a busca e envia os IPs das máquinas que possuem o arquivo
			result := search(hash)
			if len(result) == 0 {
				fmt.Fprintln(c, "Não achamos nenhum arquivo com o mesmo hash")
			} else {
				for _, ip := range result {
					fmt.Fprintln(c, ip)
				}
			}
		}
	}
}

func search(hash int) []string {
	var result []string

	// Caminho para o diretório local
	dirPath := "/tmp/dataset"
	files, err := os.ReadDir(dirPath)
	if err != nil {
		fmt.Printf("Error reading directory %s: %v", dirPath, err)
		return result
	}

	// Verifica os arquivos locais primeiro
	for _, file := range files {
		filePath := dirPath + "/" + file.Name()
		fileHash, err := fileToHash(filePath)
		if err != nil {
			continue
		}
		// Se o hash do arquivo corresponde ao hash pesquisado
		if fileHash == hash {
			// Adiciona o IP da máquina local aos resultados
			result = append(result, ipsConfigs.Ips[0]) // Aqui você pode usar o IP real da máquina local, se necessário
		}
	}

	// Agora vamos verificar as outras máquinas listadas em ips.json
	for _, ip := range ipsConfigs.Ips {
		if ip == "localhost" {
			// Já verificamos os arquivos locais, então pulamos o IP da máquina local
			continue
		}

		// Abre uma conexão TCP com a máquina remota
		conn, err := net.Dial("tcp", ip+":8000")
		if err != nil {
			fmt.Printf("Error connecting to %s: %v\n", ip, err)
			continue
		}
		defer conn.Close()

		// Envia a requisição de busca para a máquina remota
		message := fmt.Sprintf("search %d\n", hash)
		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Printf("Error sending message to %s: %v\n", ip, err)
			continue
		}

		// Lê a resposta da máquina remota
		netData, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading response from %s: %v\n", ip, err)
			continue
		}

		// Se a máquina remota tiver o arquivo, adiciona o IP dela à lista de resultados
		if strings.TrimSpace(netData) == "found" {
			result = append(result, ip)
		}
	}

	// Se nenhum IP foi encontrado com o hash solicitado
	if len(result) == 0 {
		fmt.Println("Não achamos nenhum arquivo com o mesmo hash")
	}

	return result
}

func fileToHash(filePath string) (int, error) {
	data, _, err := readFile(filePath)
	if err != nil {
		return 0, err
	}

	hash := 0
	for _, _byte := range data {
		hash += int(_byte)
	}

	return hash, nil
}

func readFile(filePath string) ([]byte, string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file %s: %v", filePath, err)
		return nil, "", err
	}

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		fmt.Printf("Error reading file %s: %v", filePath, err)
		return nil, "", err
	}

	lastModified := fileInfo.ModTime().String()
	return data, lastModified, nil
}

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

var ipsConfigs IpsConfigs

func main() {
	// Carregar IPs das máquinas
	loadIpsConfigs()

	// Escuta na porta 8000
	listener, err := net.Listen("tcp", "150.165.42.168:8000")
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
		// Ler a entrada do cliente
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			return
		}
		netData = strings.TrimSpace(netData)
		parts := strings.SplitN(netData, " ", 2)
		if len(parts) < 2 {
			continue
		}

		// Comando "search" recebido
		if parts[0] == "search" {
			hash, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Fprintln(c, "Invalid hash")
				continue
			}

			// Verifica localmente se o arquivo com o hash existe
			localIP := getLocalIP()
			foundMachines := []string{}

			if searchForHash(hash) {
				foundMachines = append(foundMachines, localIP)
			}

			// Busca nas outras máquinas da rede
			otherMachines := broadcastSearch(hash)

			// Adiciona as outras máquinas encontradas
			foundMachines = append(foundMachines, otherMachines...)

			// Envia a resposta para o cliente
			if len(foundMachines) > 0 {
				for _, ip := range foundMachines {
					fmt.Fprintln(c, ip)
				}
			} else {
				fmt.Fprintln(c, "Nenhuma máquina possui o arquivo.")
			}
		}
	}
}

func searchForHash(hash int) bool {
	// Diretório onde estão os arquivos
	dirPath := "/tmp/dataset"

	// Lê todos os arquivos do diretório
	files, err := os.ReadDir(dirPath)
	if err != nil {
		fmt.Printf("Error reading directory %s: %v", dirPath, err)
		return false
	}

	// Itera sobre cada arquivo no diretório
	for _, file := range files {
		filePath := dirPath + "/" + file.Name()

		// Calcula o hash do arquivo
		fileHash, err := fileToHash(filePath)
		if err != nil {
			continue
		}

		// Verifica se o hash do arquivo é igual ao hash procurado
		if fileHash == hash {
			return true
		}
	}

	// Nenhum arquivo com o hash correspondente foi encontrado
	return false
}

func broadcastSearch(hash int) []string {
	foundMachines := []string{}

	// Itera sobre todos os IPs
	for _, ip := range ipsConfigs.Ips {
		// Pula o IP da máquina local
		if ip == getLocalIP() {
			continue
		}

		// Conecta-se a cada IP
		conn, err := net.DialTimeout("tcp", ip+":8000", 2*time.Second)
		if err != nil {
			fmt.Println("Erro ao conectar-se a", ip, ":", err)
			continue
		}
		defer conn.Close()

		// Envia comando "search" seguido pelo hash
		fmt.Fprintf(conn, "search "+strconv.Itoa(hash)+"\n")

		// Lê a resposta do servidor
		message, _ := bufio.NewReader(conn).ReadString('\n')
		message = strings.TrimSpace(message)

		if message != "Nenhuma máquina possui o arquivo." {
			foundMachines = append(foundMachines, ip)
		}
	}

	return foundMachines
}

func getLocalIP() string {
	// Substitua por uma função que retorna o IP real da máquina local
	return "150.165.42.168" // Exemplo estático, ajuste conforme necessário
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

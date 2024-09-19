package hashFiles

import (
	"fmt"
	"net"
	"os"
)

// Função que lê o conteúdo de um arquivo e retorna seu hash
func fileToHash(filePath string) (int, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return 0, err
	}

	hash := 0
	for _, b := range data {
		hash += int(b)
	}

	return hash, nil
}

// Envia o hash de todos os arquivos de um diretório para o servidor
func SendHash(conn net.Conn, clientIP string) {
	dirPath := "/tmp/dataset" // Caminho para a pasta que contém os arquivos

	files, err := os.ReadDir(dirPath)
	if err != nil {
		fmt.Printf("Erro ao ler o diretório %s: %v\n", dirPath, err)
		return
	}

	var list []int
	for _, file := range files {
		filePath := dirPath + "/" + file.Name()
		hash, err := fileToHash(filePath)
		if err != nil {
			fmt.Printf("Erro ao calcular o hash do arquivo %s: %v\n", filePath, err)
			continue
		}
		list = append(list, hash)
		if err != nil {
			fmt.Printf("Erro ao enviar hash para o servidor: %v\n", err)
			return
		}
	}

	// Enviar o hash para o servidor no formato "update <hash>"
	message := fmt.Sprintf("update %d\n", list)
	_, err = conn.Write([]byte(message))
}
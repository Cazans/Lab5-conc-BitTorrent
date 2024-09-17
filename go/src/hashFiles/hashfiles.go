package hashFiles

import (
	"fmt"
	"net"
	"os"
)

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

// Retorna o valor referente ao hash de um arquivo
func fileToHash(filePath string) (int, string, error) {
	data, lastModified, err := readFile(filePath)

	if err != nil {
		return 0, " ", err
	}

	hash := 0
	for _, _byte := range data {
		hash += int(_byte)
	}

	return hash, lastModified, nil
}

// Envia o hash de todos os arquivos de um diretório para o servidor
func SendHash(conn net.Conn, clientIP string) {
	dirPath := "/tmp/dataset" // Atualize conforme necessário

	files, err := os.ReadDir(dirPath)
	if err != nil {
		fmt.Printf("Error reading directory %s: %v\n", dirPath, err)
		return
	}

	for _, file := range files {
		filePath := dirPath + "/" + file.Name()
		hash, _, err := fileToHash(filePath)
		if err != nil {
			fmt.Printf("Error hashing file %s: %v\n", filePath, err)
			continue
		}

		// Enviar a mensagem ao servidor no formato "update <hash>"
		message := fmt.Sprintf("update %d\n", hash)
		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Printf("Error sending hash to server: %v\n", err)
			return
		}
	}
}

package hashfiles

import (
	"fmt"
	"os"
)

// Estrutura que armazena o hash e a data de modificação de um arquivo
type fileInfo struct {
	hash int
	lastModified string
}

// Lê um arquivo e retorna o conteúdo em bytes e a data de modificação
func readFile(filePath string) ([]byte ,string, error) {
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
func sendHash() {
	dirPath := "tmp/dataset"

	files, err := os.ReadDir(dirPath)

	if err != nil {
		fmt.Printf("Error reading directory %s: %v", dirPath, err)
	}

	hashes := make(chan fileInfo, len(files))

	for _, file := range files {
		filePath := dirPath + "/" + file.Name()
		go fileToHash(filePath, hashes)	
	}

	for range files {
		processedFile := <-hashes
		//enviar a variável processedFile para o servidor
	}
}

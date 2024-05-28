package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	// Verifica se o link foi fornecido como argumento
	if len(os.Args) < 2 {
		fmt.Println("Por favor, forneça um link.")
		os.Exit(1)
	}

	// Obtém o link do argumento da linha de comando
	link := os.Args[1]
	arquivoTemporario := "/tmp/statusNFE.txt"

	// Obtendo conteúdo HTML da página
	fmt.Println("Obtendo conteúdo HTML da página:", link)

	// Executa o comando wget para obter o código-fonte da página
	cmd := exec.Command("wget", "-qO-", link)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("0") // Falha ao realizar o download
		os.Exit(1)
	}

	// Cria o arquivo temporário
	file, err := os.Create(arquivoTemporario)
	if err != nil {
		fmt.Println("0") // Falha ao criar o arquivo
		os.Exit(1)
	}
	defer file.Close()

	// Salva o conteúdo da página no arquivo temporário
	_, err = file.Write(output)
	if err != nil {
		fmt.Println("0") // Falha ao salvar o conteúdo
		os.Exit(1)
	}

	// Sucesso
	fmt.Println("1") // Download realizado com sucesso
}

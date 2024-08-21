package main

import (
	"fmt"
	"os"
	"os/exec"
)

// Função para baixar a tabela e salvar em um arquivo temporário
func downloadTable() error {
	url := "https://www.nfe.fazenda.gov.br/portal/disponibilidade.aspx?versao=0.00&tipoConteudo=P2c98tUpxrI="
	outputFile := "/tmp/statusNFE.txt"
	// Remover o arquivo existente, se houver
	if _, err := os.Stat(outputFile); err == nil {
		err := os.Remove(outputFile)
		if err != nil {
			fmt.Println("Erro ao remover arquivo existente:", err)
			return err
		}
	}

	// Baixar a tabela e salvar no arquivo
	cmd := exec.Command("curl", "-b", "session=", "-s", "-k", "-o", outputFile, url)
	err := cmd.Run()
	if err != nil {
		fmt.Println("0") // Falha ao realizar o download
		return err
	}

	fmt.Println("Tabela baixada com sucesso.")
	return nil
}

func main() {
	if err := downloadTable(); err != nil {
		fmt.Println("Erro ao baixar a tabela:", err)
	}
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Por favor, forneça um estado e um status.")
		os.Exit(1)
	}

	estado := os.Args[1]
	status := os.Args[2]

	// Lê a URL a partir do arquivo de configuração
	url, err := lerConfiguracao("config.cfg")
	if err != nil {
		fmt.Println("0") // Falha ao ler o arquivo de configuração
		os.Exit(1)
	}

	arquivoTemporario := "/tmp/statusNFE.txt"

	// Baixa o conteúdo da página
	fmt.Println("Obtendo conteúdo HTML da página:", url)
	cmd := exec.Command("wget", "-qO-", url)
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

	fmt.Println("Arquivo salvo em:", arquivoTemporario)

	// Analisa o arquivo temporário
	result := consultarServico(arquivoTemporario, estado, status)
	fmt.Println(result)
}

func lerConfiguracao(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "url=") {
			return strings.TrimPrefix(line, "url="), nil
		}
	}

	return "", fmt.Errorf("URL não encontrada no arquivo de configuração")
}

func consultarServico(filePath, estado, status string) string {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("0") // Falha ao abrir o arquivo
		os.Exit(1)
	}
	defer file.Close()

	fmt.Println("Analisando o arquivo:", filePath)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, fmt.Sprintf("<td>%s</td>", estado)) {
			fmt.Println("Linha encontrada:", line)
			var imageURL string
			switch status {
			case "AUTORIZACAO":
				imageURL = extractImageURL(line, "imagens/bola_")
			case "RETORNO.AUT":
				imageURL = extractImageURL(line, "imagens/bola_")
			case "INUTILIZACAO":
				imageURL = extractImageURL(line, "imagens/bola_")
			case "CONSULTA.PROTOCOLO":
				if estado == "SVC-AN" || estado == "SVC-RS" {
					imageURL = extractImageURL(line, "imagens/bola_")
				} else {
					imageURL = extractImageURL(line, "imagens/bola_")
				}
			case "SERVICO":
				if estado == "SVC-AN" || estado == "SVC-RS" {
					imageURL = extractImageURL(line, "imagens/bola_")
				} else {
					imageURL = extractImageURL(line, "imagens/bola_")
				}
			case "CONSULTA.CADASTRO":
				imageURL = extractImageURL(line, "imagens/bola_")
			case "RECEPCAO.EVENTO":
				switch estado {
				case "AM":
					imageURL = extractImageURL(line, "imagens/bola_")
				case "SVC-AN", "SVC-RS":
					imageURL = extractImageURL(line, "imagens/bola_")
				default:
					imageURL = extractImageURL(line, "imagens/bola_")
				}
			default:
				return "5" // SEM DADOS
			}

			fmt.Println("URL da imagem extraída:", imageURL)
			return checkServiceStatus(imageURL)
		}
	}

	return "5" // SEM DADOS
}

func extractImageURL(line string, token string) string {
	startIndex := strings.Index(line, token)
	if startIndex == -1 {
		return ""
	}
	startIndex += len(token)
	endIndex := strings.Index(line[startIndex:], "\"")
	if endIndex == -1 {
		return ""
	}
	endIndex += startIndex
	imageURL := line[startIndex:endIndex]
	return token + imageURL
}

func checkServiceStatus(imageURL string) string {
	fmt.Println("Verificando status para a URL da imagem:", imageURL)
	switch imageURL {
	case "imagens/bola_verde_P.png":
		return "1" // DISPONIVEL
	case "imagens/bola_amarela_P.png":
		return "2" // INDISPONIVEL
	case "imagens/bola_vermelho_P.png":
		return "10" // OFFLINE
	default:
		return "5" // SEM DADOS
	}
}

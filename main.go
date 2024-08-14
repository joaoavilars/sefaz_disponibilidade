package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/net/html"
)

// Estrutura para mapear os serviços
type EstadoStatus struct {
	Nome    string
	Estados string
	Indice  int
}

var estadoStatus = []EstadoStatus{
	{"AUTORIZACAO", "AM,BA,CE,GO,MG,MS,MT,PE,PR,RS,SP,SVAN,SVRS,SVC-AN,SVC-RS", 2},
	{"RETORNO.AUT", "AM,BA,CE,GO,MG,MS,MT,PE,PR,RS,SP,SVAN,SVRS,SVC-AN,SVC-RS", 4},
	{"INUTILIZACAO", "AM,BA,CE,GO,MG,MS,MT,PE,PR,RS,SP,SVAN,SVRS", 6},
	{"CONSULTA.PROTOCOLO", "AM,BA,CE,GO,MG,MS,MT,PE,PR,RS,SP,SVAN,SVRS", 8},
	{"SERVICO", "AM,BA,CE,GO,MG,MS,MT,PE,PR,RS,SP,SVAN,SVRS", 10},
	{"CONSULTA.CADASTRO", "BA,CE,GO,MG,MS,MT,PE,PR,RS,SP,SVRS", 12},
	{"RECEPCAO.EVENTO", "BA,CE,GO,MG,MS,MT,PE,PR,RS,SP,SVRS", 14},
}

// Função para baixar a tabela e salvar em um arquivo temporário
func downloadTable() error {
	url := "https://www.nfe.fazenda.gov.br/portal/disponibilidade.aspx?versao=0.00&tipoConteudo=P2c98tUpxrI="
	outputFile := "/tmp/statusNFE.txt"

	cmd := exec.Command("curl", "-b", "session=", "-s", "-k", "-o", outputFile, url)
	err := cmd.Run()
	if err != nil {
		fmt.Println("0") // Falha ao realizar o download
		return err
	}

	return nil
}

// Função para normalizar o texto, removendo acentos, caracteres especiais e transformando em maiúsculas
func normalizeString(s string) string {
	s = strings.ToUpper(s)
	replacer := strings.NewReplacer(
		"Ç", "C", "Ã", "A", "Õ", "O",
		"Á", "A", "É", "E", "Í", "I", "Ó", "O", "Ú", "U",
		"À", "A", "È", "E", "Ì", "I", "Ò", "O", "Ù", "U",
		"Â", "A", "Ê", "E", "Î", "I", "Ô", "O", "Û", "U",
		" ", ".",
	)
	return replacer.Replace(s)
}

// Função para verificar o status do serviço com base na URL da imagem
func checkServiceStatus(imageURL string) string {
	switch imageURL {
	case "imagens/bola_verde_P.png":
		return "1" // DISPONÍVEL
	case "imagens/bola_amarela_P.png":
		return "2" // INDISPONÍVEL
	case "imagens/bola_vermelho_P.png":
		return "0" // OFFLINE
	default:
		return "5" // SEM DADOS
	}
}

// Função para percorrer a árvore HTML e extrair os cabeçalhos da tabela
func extractTableHeaders(n *html.Node) []string {
	var headers []string
	if n.Type == html.ElementNode && n.Data == "th" {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.TextNode {
				headers = append(headers, strings.TrimSpace(c.Data))
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		headers = append(headers, extractTableHeaders(c)...)
	}
	return headers
}

// Função para localizar a linha correspondente à UF
func findUFRow(n *html.Node, uf string) *html.Node {
	if n.Type == html.ElementNode && n.Data == "tr" {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode && c.Data == "td" && c.FirstChild != nil && c.FirstChild.Type == html.TextNode && strings.TrimSpace(c.FirstChild.Data) == uf {
				return n
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if found := findUFRow(c, uf); found != nil {
			return found
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: ./programa UF SERVICO")
		return
	}

	uf := os.Args[1]
	servico := os.Args[2]

	// Mensagem de depuração: Mostrar UF e SERVIÇO digitados
	//fmt.Printf("UF digitada: %s, SERVIÇO digitado: %s\n", uf, servico)

	// Baixar a tabela e salvar no arquivo
	if err := downloadTable(); err != nil {
		return
	}

	// Ler o conteúdo do arquivo
	content, err := ioutil.ReadFile("/tmp/statusNFE.txt")
	if err != nil {
		fmt.Println("0") // Falha ao ler o arquivo
		return
	}

	// Parse do conteúdo HTML
	doc, err := html.Parse(strings.NewReader(string(content)))
	if err != nil {
		fmt.Println("Erro ao fazer o parse do HTML:", err)
		return
	}

	// Extrair cabeçalhos da tabela
	headers := extractTableHeaders(doc)
	//fmt.Println("Cabeçalhos extraídos da tabela:", headers)

	// Normalizar os cabeçalhos
	var normalizedHeaders []string
	for _, header := range headers {
		normalizedHeader := normalizeString(header)
		normalizedHeaders = append(normalizedHeaders, normalizedHeader)
	}

	//fmt.Println("Cabeçalhos pós-normalização:", normalizedHeaders)

	// Ajuste para maiúsculas e remoção do número 4
	for i, header := range normalizedHeaders {
		if strings.HasSuffix(header, "4") {
			normalizedHeaders[i] = header[:len(header)-1]
		}
	}
	//fmt.Println("Cabeçalhos após remoção do '4':", normalizedHeaders)

	// Localizar o índice do serviço correspondente
	columnIndex := -1
	for i, header := range normalizedHeaders {
		if header == servico {
			columnIndex = i
			break
		}
	}

	if columnIndex == -1 {
		fmt.Println("Serviço não encontrado:", servico)
		fmt.Println("5")
		return
	}

	//fmt.Printf("Serviço '%s' encontrado na coluna %d.\n", servico, columnIndex)

	// Localizar a linha da UF
	ufRow := findUFRow(doc, uf)
	if ufRow == nil {
		fmt.Printf("UF '%s' não encontrada na tabela.\n", uf)
		fmt.Println("5") // UF não encontrada
		return
	}

	// Percorrer a linha para encontrar a coluna correspondente ao serviço
	colCount := 0
	var serviceStatusURL string
	for c := ufRow.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "td" {
			if colCount == columnIndex {
				if img := c.FirstChild; img != nil && img.Type == html.ElementNode && img.Data == "img" {
					for _, attr := range img.Attr {
						if attr.Key == "src" {
							serviceStatusURL = attr.Val
							break
						}
					}
				}
				break
			}
			colCount++
		}
	}

	if serviceStatusURL == "" {
		//fmt.Printf("Serviço '%s' não encontrado na UF '%s'.\n", servico, uf)
		fmt.Println("5") // Serviço não encontrado ou erro na leitura do status
		return
	}

	status := checkServiceStatus(serviceStatusURL)
	//fmt.Printf("Status do serviço '%s' para a UF '%s': %s\n", servico, uf, status)
	fmt.Println(status)
}

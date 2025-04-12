package codesp

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type Codesp struct {
	date string
	html string
}

func NewCodesp() *Codesp {
	html, err := getCodespData()
	if err != nil {
		return nil
	}

	newDate, err := extractDate(html)
	if err != nil {
		return nil
	}

	return &Codesp{
		date: newDate,
		html: html,
	}
}

func (c *Codesp) ExtractTable() map[string]map[string]string {

	html := c.html

	tableRegexp := regexp.MustCompile(`(?i)<table[\s\S]*?</table>`)
	rowRegexp := regexp.MustCompile(`(?i)<tr[\s\S]*?</tr>`)
	headerRegexp := regexp.MustCompile(`(?i)<th[^>]*>([^<]+)*?</th>`)
	rowDataRegexp := regexp.MustCompile(`(?i)<td[^>]*>(.*?)</td>`)
	commentRegex := regexp.MustCompile(`<!--([\s\S]*?)-->`)

	table := tableRegexp.FindString(html)

	if table == "" {
		fmt.Println("Nenhuma tabela encontrada")
		return nil
	}

	rows := rowRegexp.FindAllString(table, -1)

	if len(rows) < 3 {
		fmt.Println("Nenhuma linha encontrada")
		return nil
	}

	headers := headerRegexp.FindAllStringSubmatch(rows[1], -1)

	if len(headers) < 1 {
		fmt.Println("Nenhum cabeçalho encontrado")
		return nil
	}

	filteredHeaders := make([]string, 0, len(headers))

	for i, header := range headers {
		if i == 0 {
			continue
		}
		filteredHeaders = append(filteredHeaders, header[1])
	}

	rowsData := make(map[string]map[string]string)

	for _, row := range rows[2:] {

		cleanCells := commentRegex.ReplaceAllString(row, "")

		cells := rowDataRegexp.FindAllStringSubmatch(cleanCells, -1)

		if len(cells) < 1 {
			continue
		}

		berth := cells[0][1]

		cols := make([]string, 0, len(headers))

		for i, cell := range cells {
			if i > 0 {
				if strings.Contains(cell[1], "<!--") {
					continue
				}
				if strings.Contains(cell[1], "Trecho") {
					continue
				}
				if strings.Contains(cell[1], "strong") {
					cols = append(cols, cell[1])
					cols = append(cols, cell[1])
					continue
				}
				cols = append(cols, cell[1])
			}
		}

		if len(cols) < 1 {
			continue
		}

		if _, exists := rowsData[berth]; !exists {
			rowsData[berth] = make(map[string]string)
		}

		for i := range filteredHeaders {
			rowsData[berth][filteredHeaders[i]] = cols[i]
		}
	}

	return rowsData
}

func (c *Codesp) GetDate() string {
	return c.date
}

func getCodespData() (string, error) {
	url := "https://www.portodesantos.com.br/informacoes-operacionais/operacoes-portuarias/calados-operacionais-dos-bercos-de-atracacao/"

	// Disabling SSL verification
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := httpClient.Get(url)
	if err != nil {
		fmt.Println("Erro ao fazer request:", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erro ao ler resposta:", err)
		return "", err
	}

	// fmt.Println("Resposta do servidor:", string(body))s

	return string(body), nil
}

func extractDate(html string) (string, error) {
	dateRegexp := regexp.MustCompile(`(?i)Data:([\s\S]*?)</`)

	date := dateRegexp.FindStringSubmatch(html)

	if date == nil {
		fmt.Println("Nenhuma data encontrada")
		return "", errors.New("nenhuma data encontrada")
	}

	return strings.Trim(date[1], " "), nil
}

// func extractObservation(html string) string {

// 	sectionRegexp := regexp.MustCompile(`(?i)Observações([\s\S]*?)</ol>`)
// 	listRegexp := regexp.MustCompile(`(?i)<li>([\s\S]*?)</li>`)
// 	rowRegexp := regexp.MustCompile(`(?i)<tr[\s\S]*?</tr>`)
// 	headerRegexp := regexp.MustCompile(`(?i)<th[^>]*>([^<]+)*?</th>`)
// 	rowDataRegexp := regexp.MustCompile(`(?i)<td[^>]*>(.*?)</td>`)
// 	commentRegex := regexp.MustCompile(`<!--([\s\S]*?)-->`)
// 	rowSpanRegexp := regexp.MustCompile(`(?i)rowspan="([0-9])+"`)

// 	section := sectionRegexp.FindString(html)

// 	if section == "" {
// 		fmt.Println("Nenhuma observação encontrada")
// 		return ""
// 	}

// 	list := listRegexp.FindAllStringSubmatch(section, -1)

// 	if len(list) < 1 {
// 		fmt.Println("Nenhuma observação encontrada")
// 		return ""
// 	}

// 	for _, item := range list {
// 		if strings.Contains(item[1], "<table") {
// 			rows := rowRegexp.FindAllString(item[1], -1)

// 			if len(rows) < 1 {
// 				fmt.Println("Nenhuma linha encontrada")
// 				return ""
// 			}

// 			headers := headerRegexp.FindAllStringSubmatch(rows[0], -1)

// 			if len(headers) < 1 {
// 				fmt.Println("Nenhum cabeçalho encontrado")
// 				return ""
// 			}

// 			filteredHeaders := make(map[string][]string)

// 			// filteredHeaders are a map of the berths containing the headers

// 			for i, header := range headers {
// 				if i == 0 {
// 					continue
// 				}
// 				filteredHeaders[headers[0][1]] = append(filteredHeaders[headers[0][1]], header[1])
// 			}

// 			rowsData := make(map[string]map[string]map[string]string)

// 			additionalCells := make(map[int]int)

// 			// [ cellIndex: cellToEnd ]

// 			for rowIndex, row := range rows[1:] {

// 				// Iterating over the rows of the table

// 				cleanCells := commentRegex.ReplaceAllString(row, "")

// 				cells := rowDataRegexp.FindAllStringSubmatch(cleanCells, -1)

// 				if len(cells) < 1 {
// 					continue
// 				}

// 				cols := make([]string, 0, len(headers))

// 				// Cols represents the individual values of each column

// 				for cellIndex, cell := range cells {

// 					if rowIndex > 0 && rowIndex <= additionalCells[cellIndex] {
// 						// fmt.Println("BATATAAAAAAAAAAA", rowIndex, cellIndex, additionalCells[cellIndex])
// 						cols = append(cols, cell[1])
// 					}

// 					cols = append(cols, cell[1])

// 					teste := rowSpanRegexp.FindStringSubmatch(cell[0])

// 					if len(teste) < 1 {
// 						continue
// 					}
// 					//string to int
// 					num, err := strconv.Atoi(teste[1])

// 					if err != nil {
// 						fmt.Println(err)
// 						continue
// 					}

// 					if num > 1 {
// 						additionalCells[cellIndex] = rowIndex + num - 1
// 					}
// 				}

// 				for berth, berthHeaders := range filteredHeaders {

// 					if _, exists := rowsData[berth]; !exists {
// 						rowsData[berth] = make(map[string]map[string]string) // Inicializa o mapa interno
// 					}

// 					if _, exists := rowsData[berth][berthHeaders[0]]; !exists {
// 						rowsData[berth][cols[0]] = make(map[string]string) // Inicializa o mapa interno
// 					}

// 					for i, header := range berthHeaders {
// 						if i == 0 {
// 							continue
// 						}
// 						rowsData[berth][cols[0]][header] = cols[i]
// 					}
// 				}
// 			}

// 			for key, item := range rowsData {
// 				fmt.Println("Berço:", key)

// 				for itemKey, itemValue := range item {
// 					fmt.Println("	Tipo:", itemKey)

// 					for dataKey, dataValue := range itemValue {
// 						fmt.Println("	\t", dataKey+":", dataValue)
// 					}
// 				}
// 			}
// 		}
// 	}

// 	// fmt.Println(list)
// 	return ""
// }

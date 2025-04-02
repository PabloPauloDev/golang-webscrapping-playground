package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
)

const connStr = ""

func getLastChange(actualDate string) bool {
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		fmt.Println("Erro ao conectar ao banco de dados:", err)
	}
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), "SELECT date FROM last_changes")

	if err != nil {
		fmt.Println("Erro ao executar consulta:", err)
	}
	defer rows.Close()

	var lastDate string

	for rows.Next() {
		if err := rows.Scan(&lastDate); err != nil {
			fmt.Println("Erro ao escanear a linha:", err)
		}
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Erro ao iterar resultados:", err)
	}

	fmt.Println("Actual change:", actualDate)
	fmt.Println("Last change:", lastDate)

	hasChanged := actualDate != lastDate

	if hasChanged {
		fmt.Println("Change found")
	}

	return hasChanged
}
func setLastChange(actualDate string) bool {

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		fmt.Println("Erro ao conectar ao banco de dados:", err)
	}
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), "INSERT INTO last_changes (date) VALUES ($1)", actualDate)

	if err != nil {
		fmt.Println("Erro ao executar consulta:", err)
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		fmt.Println("Erro ao iterar resultados:", err)
	} else {
		fmt.Println("Last change updated")
	}

	return err == nil
}
func getHeader(header string) string {

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		fmt.Println("Erro ao conectar ao banco de dados:", err)
	}
	defer conn.Close(context.Background())

	var result string
	err = conn.QueryRow(context.Background(), "SELECT id_header FROM headers WHERE nm_header = $1", header).Scan(&result)

	if err != nil {
		fmt.Println("Erro ao executar consulta(Header): ", err)
		return ""
	}

	return result
}
func insertHeader(header string) string {

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		fmt.Println("Erro ao conectar ao banco de dados:", err)
	}
	defer conn.Close(context.Background())

	idHeader := ""

	err = conn.QueryRow(context.Background(), "INSERT INTO headers (nm_header) VALUES ($1) RETURNING id_header", header).Scan(&idHeader)

	if err != nil {
		fmt.Println("Erro ao executar consulta:", err)
	}

	return idHeader
}
func compareData(newData, oldData map[string]map[string]string) {
	changesText := ""
	for key, value := range newData {
		if _, exists := oldData[key]; !exists {
			fmt.Println("Novo dado:", key, value)
			for k, v := range value {
				insertData(key, k, v)
			}
		} else {
			for k, v := range value {
				if oldData[key][k] != v {
					changesText += fmt.Sprintf("<strong>%s</strong>: %s -> %s<br>", key, k, v)
					updateData(key, k, v)
				}
			}
		}
	}

	if changesText != "" {
		sendEmail(changesText)
	}
}
func sendEmail(message string) {
	fmt.Println("Enviando email...")
	url := ""

	emailBody := map[string]string{
		"from":    "sistema@sts.amorion.com.br",
		"to":      "pablo.carpanedo@sts.amorion.com.br",
		"subject": "Alteração nos berços - CODESP",
		"body":    `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd"><html dir="ltr" xmlns="http://www.w3.org/1999/xhtml" xmlns:o="urn:schemas-microsoft-com:office:office" lang="pt"><head><meta charset="UTF-8"><meta content="width=device-width, initial-scale=1" name="viewport"><meta name="x-apple-disable-message-reformatting"><meta http-equiv="X-UA-Compatible" content="IE=edge"><meta content="telephone=no" name="format-detection"><title>Empty template</title> <!--[if (mso 16)]><style type="text/css"> a {text-decoration: none;}  </style><![endif]--><!--[if gte mso 9]><style>sup { font-size: 100% !important; }</style><![endif]--><!--[if gte mso 9]><noscript> <xml> <o:OfficeDocumentSettings> <o:AllowPNG></o:AllowPNG> <o:PixelsPerInch>96</o:PixelsPerInch> </o:OfficeDocumentSettings> </xml> </noscript><![endif]--><!--[if mso]><xml> <w:WordDocument xmlns:w="urn:schemas-microsoft-com:office:word"> <w:DontUseAdvancedTypographyReadingMail/> </w:WordDocument> </xml><![endif]--><style type="text/css">.rollover:hover .rollover-first { max-height:0px!important; display:none!important;}.rollover:hover .rollover-second { max-height:none!important; display:block!important;}.rollover span { font-size:0px;}u + .body img ~ div div { display:none;}#outlook a { padding:0;}span.MsoHyperlink,span.MsoHyperlinkFollowed { color:inherit; mso-style-priority:99;}a.es-button { mso-style-priority:100!important; text-decoration:none!important;}a[x-apple-data-detectors],#MessageViewBody a { color:inherit!important; text-decoration:none!important; font-size:inherit!important; font-family:inherit!important; font-weight:inherit!important; line-height:inherit!important;}.es-desk-hidden { display:none; float:left; overflow:hidden; width:0; max-height:0; line-height:0; mso-hide:all;}@media only screen and (max-width:600px) {.es-m-p20b { padding-bottom:20px!important } .es-p-default { }*[class="gmail-fix"] { display:none!important } p, a { line-height:150%!important } h1, h1 a { line-height:120%!important } h2, h2 a { line-height:120%!important } h3, h3 a { line-height:120%!important } h4, h4 a { line-height:120%!important } h5, h5 a { line-height:120%!important } h6, h6 a { line-height:120%!important } .es-header-body p { } .es-content-body p { } .es-footer-body p { } .es-infoblock p { } h1 { font-size:40px!important; text-align:left } h2 { font-size:32px!important; text-align:left } h3 { font-size:28px!important; text-align:left } h4 { font-size:24px!important; text-align:left } h5 { font-size:20px!important; text-align:left } h6 { font-size:16px!important; text-align:left } .es-header-body h1 a, .es-content-body h1 a, .es-footer-body h1 a { font-size:40px!important } .es-header-body h2 a, .es-content-body h2 a, .es-footer-body h2 a { font-size:32px!important }.es-header-body h3 a, .es-content-body h3 a, .es-footer-body h3 a { font-size:28px!important } .es-header-body h4 a, .es-content-body h4 a, .es-footer-body h4 a { font-size:24px!important } .es-header-body h5 a, .es-content-body h5 a, .es-footer-body h5 a { font-size:20px!important } .es-header-body h6 a, .es-content-body h6 a, .es-footer-body h6 a { font-size:16px!important } .es-menu td a { font-size:14px!important } .es-header-body p, .es-header-body a { font-size:14px!important } .es-content-body p, .es-content-body a { font-size:14px!important } .es-footer-body p, .es-footer-body a { font-size:14px!important } .es-infoblock p, .es-infoblock a { font-size:12px!important } .es-m-txt-c, .es-m-txt-c h1, .es-m-txt-c h2, .es-m-txt-c h3, .es-m-txt-c h4, .es-m-txt-c h5, .es-m-txt-c h6 { text-align:center!important }.es-m-txt-r, .es-m-txt-r h1, .es-m-txt-r h2, .es-m-txt-r h3, .es-m-txt-r h4, .es-m-txt-r h5, .es-m-txt-r h6 { text-align:right!important } .es-m-txt-j, .es-m-txt-j h1, .es-m-txt-j h2, .es-m-txt-j h3, .es-m-txt-j h4, .es-m-txt-j h5, .es-m-txt-j h6 { text-align:justify!important } .es-m-txt-l, .es-m-txt-l h1, .es-m-txt-l h2, .es-m-txt-l h3, .es-m-txt-l h4, .es-m-txt-l h5, .es-m-txt-l h6 { text-align:left!important } .es-m-txt-r img, .es-m-txt-c img, .es-m-txt-l img { display:inline!important } .es-m-txt-r .rollover:hover .rollover-second, .es-m-txt-c .rollover:hover .rollover-second, .es-m-txt-l .rollover:hover .rollover-second { display:inline!important } .es-m-txt-r .rollover span, .es-m-txt-c .rollover span, .es-m-txt-l .rollover span { line-height:0!important; font-size:0!important; display:block } .es-spacer { display:inline-table }a.es-button, button.es-button { font-size:14px!important; padding:10px 20px 10px 20px!important; line-height:120%!important } a.es-button, button.es-button, .es-button-border { display:inline-block!important } .es-m-fw, .es-m-fw.es-fw, .es-m-fw .es-button { display:block!important } .es-m-il, .es-m-il .es-button, .es-social, .es-social td, .es-menu { display:inline-block!important } .es-adaptive table, .es-left, .es-right { width:100%!important } .es-content table, .es-header table, .es-footer table, .es-content, .es-footer, .es-header { width:100%!important; max-width:600px!important } .adapt-img { width:100%!important; height:auto!important } .es-mobile-hidden, .es-hidden { display:none!important } .es-desk-hidden { width:auto!important; overflow:visible!important; float:none!important; max-height:inherit!important; line-height:inherit!important } tr.es-desk-hidden { display:table-row!important }table.es-desk-hidden { display:table!important } td.es-desk-menu-hidden { display:table-cell!important } .es-menu td { width:1%!important } table.es-table-not-adapt, .esd-block-html table { width:auto!important } .h-auto { height:auto!important } }@media screen and (max-width:384px) {.mail-message-content { width:414px!important } }</style></head> <body class="body" style="width:100%;height:100%;-webkit-text-size-adjust:100%;-ms-text-size-adjust:100%;padding:0;Margin:0"><div dir="ltr" class="es-wrapper-color" lang="pt" style="background-color:#F6F6F6"><!--[if gte mso 9]><v:background xmlns:v="urn:schemas-microsoft-com:vml" fill="t"> <v:fill type="tile" color="#f6f6f6"></v:fill> </v:background><![endif]--><table width="100%" cellspacing="0" cellpadding="0" class="es-wrapper" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px;padding:0;Margin:0;width:100%;height:100%;background-color:#F6F6F6"><tr><td valign="top" style="padding:0;Margin:0"><table cellspacing="0" cellpadding="0" align="center" class="es-header" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px;width:100%;table-layout:fixed !important;background-color:transparent"><tr><td align="center" style="padding:0;Margin:0"><table cellspacing="0" cellpadding="0" bgcolor="#ffffff" align="center" class="es-header-body" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px;background-color:#FFFFFF;width:600px"><tr><td align="left" style="padding:0;Margin:0;padding-top:20px;padding-right:20px;padding-left:20px"><!--[if mso]><table style="width:560px" cellpadding="0" cellspacing="0"><tr><td style="width:180px" valign="top"><![endif]--><table cellspacing="0" cellpadding="0" align="left" class="es-left" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px;float:left"><tr><td align="center" valign="top" class="es-m-p20b" style="padding:0;Margin:0;width:180px"><table cellpadding="0" cellspacing="0" width="100%" role="presentation" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px"><tr><td align="center" style="padding:0;Margin:0;font-size:0"><img src="https://www.portodesantos.com.br/wp-content/themes/Tema%20SPA/assets/img/LogoSPA2023.jpg" alt="" width="180" class="adapt-img" style="display:block;font-size:14px;border:0;outline:none;text-decoration:none"></td> </tr></table></td></tr></table><!--[if mso]></td><td style="width:20px"></td><td style="width:360px" valign="top"><![endif]--><table cellspacing="0" cellpadding="0" align="right" class="es-right" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px;float:right"><tr><td align="left" style="padding:0;Margin:0;width:360px"><table role="presentation" cellpadding="0" cellspacing="0" width="100%" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px"><tr><td align="center" style="padding:0;Margin:0;font-size:0"><img src="https://www.portodesantos.com.br/wp-content/themes/Tema%20SPA/assets/img/foto_homepage_nova.jpg" alt="" width="360" class="adapt-img" style="display:block;font-size:14px;border:0;outline:none;text-decoration:none"></td></tr></table></td></tr></table> <!--[if mso]></td></tr></table><![endif]--></td></tr></table></td></tr></table> <table cellspacing="0" cellpadding="0" align="center" class="es-content" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px;width:100%;table-layout:fixed !important"><tr><td align="center" style="padding:0;Margin:0"><table cellspacing="0" cellpadding="0" bgcolor="#ffffff" align="center" class="es-content-body" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px;background-color:#FFFFFF;width:600px"><tr><td align="left" style="padding:0;Margin:0;padding-top:20px;padding-right:20px;padding-left:20px"><table width="100%" cellspacing="0" cellpadding="0" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px"><tr><td valign="top" align="center" style="padding:0;Margin:0;width:560px"><table width="100%" cellspacing="0" cellpadding="0" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px"><tr><td align="left" style="padding:0;Margin:0"><h3 style="Margin:0;font-family:arial, 'helvetica neue', helvetica, sans-serif;mso-line-height-rule:exactly;letter-spacing:0;font-size:28px;font-style:normal;font-weight:normal;line-height:33.6px;color:#333333"><strong>Altreração nos berços - CODESP</strong></h3> </td></tr></table></td></tr></table></td></tr></table></td></tr></table><table cellspacing="0" cellpadding="0" align="center" class="es-footer" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px;width:100%;table-layout:fixed !important;background-color:transparent"><tr><td align="center" style="padding:0;Margin:0"><table cellspacing="0" cellpadding="0" bgcolor="#ffffff" align="center" class="es-footer-body" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px;background-color:#FFFFFF;width:600px"><tr><td align="left" style="padding:0;Margin:0;padding-top:20px;padding-right:20px;padding-left:20px"><table width="100%" cellpadding="0" cellspacing="0" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px"><tr><td align="left" style="padding:0;Margin:0;width:560px"><table cellpadding="0" cellspacing="0" width="100%" role="presentation" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px"><tr><td align="left" style="padding:0;Margin:0"><p style="Margin:0;mso-line-height-rule:exactly;font-family:arial, 'helvetica neue', helvetica, sans-serif;line-height:21px;letter-spacing:0;color:#333333;font-size:14px">` + message + `</p> </td></tr></table></td></tr></table></td></tr> <tr><td align="left" style="padding:0;Margin:0;padding-top:20px;padding-right:20px;padding-left:20px"><table width="100%" cellpadding="0" cellspacing="0" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px"><tr><td align="left" style="padding:0;Margin:0;width:560px"><table cellpadding="0" cellspacing="0" width="100%" role="presentation" style="mso-table-lspace:0pt;mso-table-rspace:0pt;border-collapse:collapse;border-spacing:0px"><tr><td align="center" style="padding:0;Margin:0"><span class="es-button-border es-fw" style="border-style:solid;border-color:#2CB543;background:#31CB4B;border-width:0px 0px 2px 0px;display:block;border-radius:15px;width:auto"><a href="https://www.portodesantos.com.br/informacoes-operacionais/operacoes-portuarias/calados-operacionais-dos-bercos-de-atracacao/" target="_blank" class="es-button" style="mso-style-priority:100 !important;text-decoration:none !important;mso-line-height-rule:exactly;color:#FFFFFF;font-size:14px;padding:10px 20px 10px 20px;display:block;background:#31CB4B;border-radius:15px;font-family:arial, 'helvetica neue', helvetica, sans-serif;font-weight:normal;font-style:normal;line-height:16.8px;width:auto;text-align:center;letter-spacing:0;mso-padding-alt:0;mso-border-alt:10px solid #31CB4B">Acessar na Codesp</a> </span></td></tr></table></td></tr></table></td></tr></table></td></tr></table></td></tr></table></div></body></html>`}

	jsonData, err := json.Marshal(emailBody)
	if err != nil {
		fmt.Println("Erro ao criar JSON:", err)
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil || resp.StatusCode != 200 {
		fmt.Println("Erro ao enviar email:", err)
		return
	}
	fmt.Println("Email enviado com sucesso")
}
func updateData(berth, header, value string) {

	conn, err := pgx.Connect(context.Background(), connStr)

	if err != nil {
		fmt.Println("Erro ao conectar ao banco de dados:", err)
	}
	defer conn.Close(context.Background())

	_, err = conn.Exec(context.Background(), "UPDATE data SET value = $1 WHERE cod_header = $2 AND cod_berth = $3", value, getHeader(header), getHeader(berth))

	if err != nil {
		fmt.Println("Erro ao executar update:", err)
		return
	}

	fmt.Println("Dado atualizado:", berth, header, value)

}
func insertData(berth, header, value string) {
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		fmt.Println("Erro ao conectar ao banco de dados:", err)
	}
	defer conn.Close(context.Background())

	idBerth, idHeader := "", ""

	err = conn.QueryRow(context.Background(), "SELECT id_header FROM headers WHERE nm_header = $1", berth).Scan(&idBerth)
	if err != nil {
		fmt.Println("Erro ao buscar idBerth:", err)
		idBerth = insertHeader(berth)
	}

	err = conn.QueryRow(context.Background(), "SELECT id_header FROM headers WHERE nm_header = $1", header).Scan(&idHeader)
	if err != nil {
		fmt.Println("Erro ao buscar idHeader:", err)
		idHeader = insertHeader(header)
	}

	_, err = conn.Exec(context.Background(), "INSERT INTO data (cod_berth, cod_header, value) VALUES ($1, $2, $3)", idBerth, idHeader, value)

	if err != nil {
		fmt.Println("Erro ao executar consulta:", err)
	}
}
func getLastData() map[string]map[string]string {
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		fmt.Println("Erro ao conectar ao banco de dados:", err)
	}
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), "SELECT h1.nm_header AS header, h2.nm_header AS berth, d.value FROM data d LEFT JOIN headers h1 ON h1.id_header = d.cod_header LEFT JOIN headers h2 ON h2.id_header = d.cod_berth; ")
	if err != nil {
		fmt.Println("Erro ao executar consulta:", err)
	}
	defer rows.Close()

	data := make(map[string]map[string]string)

	for rows.Next() {
		var codBerth, codHeader, value string
		if err := rows.Scan(&codHeader, &codBerth, &value); err != nil {
			fmt.Println("Erro ao escanear a linha:", err)
		}
		if _, exists := data[codBerth]; !exists {
			data[codBerth] = make(map[string]string)
		}
		data[codBerth][codHeader] = value
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Erro ao iterar resultados:", err)
	}

	return data
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
func extractTable(html string) map[string]map[string]string {

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
func extractObservation(html string) string {

	sectionRegexp := regexp.MustCompile(`(?i)Observações([\s\S]*?)</ol>`)
	listRegexp := regexp.MustCompile(`(?i)<li>([\s\S]*?)</li>`)
	rowRegexp := regexp.MustCompile(`(?i)<tr[\s\S]*?</tr>`)
	headerRegexp := regexp.MustCompile(`(?i)<th[^>]*>([^<]+)*?</th>`)
	rowDataRegexp := regexp.MustCompile(`(?i)<td[^>]*>(.*?)</td>`)
	commentRegex := regexp.MustCompile(`<!--([\s\S]*?)-->`)
	rowSpanRegexp := regexp.MustCompile(`(?i)rowspan="([0-9])+"`)

	section := sectionRegexp.FindString(html)

	if section == "" {
		fmt.Println("Nenhuma observação encontrada")
		return ""
	}

	list := listRegexp.FindAllStringSubmatch(section, -1)

	if len(list) < 1 {
		fmt.Println("Nenhuma observação encontrada")
		return ""
	}

	for _, item := range list {
		if strings.Contains(item[1], "<table") {
			rows := rowRegexp.FindAllString(item[1], -1)

			if len(rows) < 1 {
				fmt.Println("Nenhuma linha encontrada")
				return ""
			}

			headers := headerRegexp.FindAllStringSubmatch(rows[0], -1)

			if len(headers) < 1 {
				fmt.Println("Nenhum cabeçalho encontrado")
				return ""
			}

			filteredHeaders := make(map[string][]string)

			// filteredHeaders are a map of the berths containing the headers

			for i, header := range headers {
				if i == 0 {
					continue
				}
				filteredHeaders[headers[0][1]] = append(filteredHeaders[headers[0][1]], header[1])
			}

			rowsData := make(map[string]map[string]map[string]string)

			additionalCells := make(map[int]int)

			// [ cellIndex: cellToEnd ]

			for rowIndex, row := range rows[1:] {

				// Iterating over the rows of the table

				cleanCells := commentRegex.ReplaceAllString(row, "")

				cells := rowDataRegexp.FindAllStringSubmatch(cleanCells, -1)

				if len(cells) < 1 {
					continue
				}

				cols := make([]string, 0, len(headers))

				// Cols represents the individual values of each column

				for cellIndex, cell := range cells {

					if rowIndex > 0 && rowIndex <= additionalCells[cellIndex] {
						// fmt.Println("BATATAAAAAAAAAAA", rowIndex, cellIndex, additionalCells[cellIndex])
						cols = append(cols, cell[1])
					}

					cols = append(cols, cell[1])

					teste := rowSpanRegexp.FindStringSubmatch(cell[0])

					if len(teste) < 1 {
						continue
					}
					//string to int
					num, err := strconv.Atoi(teste[1])

					if err != nil {
						fmt.Println(err)
						continue
					}

					if num > 1 {
						additionalCells[cellIndex] = rowIndex + num - 1
					}
				}

				for berth, berthHeaders := range filteredHeaders {

					if _, exists := rowsData[berth]; !exists {
						rowsData[berth] = make(map[string]map[string]string) // Inicializa o mapa interno
					}

					if _, exists := rowsData[berth][berthHeaders[0]]; !exists {
						rowsData[berth][cols[0]] = make(map[string]string) // Inicializa o mapa interno
					}

					for i, header := range berthHeaders {
						if i == 0 {
							continue
						}
						rowsData[berth][cols[0]][header] = cols[i]
					}
				}
			}

			for key, item := range rowsData {
				fmt.Println("Berço:", key)

				for itemKey, itemValue := range item {
					fmt.Println("	Tipo:", itemKey)

					for dataKey, dataValue := range itemValue {
						fmt.Println("	\t", dataKey+":", dataValue)
					}
				}

			}
		}
	}

	// fmt.Println(list)
	return ""
}
func extractDate(html string) string {
	dateRegexp := regexp.MustCompile(`(?i)Data:([\s\S]*?)</`)

	date := dateRegexp.FindStringSubmatch(html)

	if date == nil {
		fmt.Println("Nenhuma data encontrada")
		return ""
	}

	return strings.Trim(date[1], " ")

}
func main() {
	codesp, err := getCodespData()

	if err != nil {
		fmt.Println(err)
		return
	}

	date := extractDate(codesp)
	hasChanged := getLastChange(date)

	if !hasChanged {
		fmt.Println("Nenhuma alteração encontrada")
		return
	}

	new := extractTable(codesp)
	if new == nil {
		return
	}

	old := getLastData()
	if old == nil {
		return
	}

	for key := range new {
		if getHeader(key) == "" {
			insertHeader(key)
			fmt.Println("Novo header:", key)
		}
	}

	compareData(new, old)
	// extractObservation(codesp)

	setLastChange(date)
}

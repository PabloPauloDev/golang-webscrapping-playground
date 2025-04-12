package db

import (
	"context"
	"fmt"

	"codesp_data/pkg/postgres"
)

func GetLastChange() string {
	conn := postgres.GetPool()

	rows, err := conn.Query(context.Background(), "SELECT date FROM last_changes")
	if err != nil {
		fmt.Println("Erro ao executar consulta:", err)
		return ""
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

	return lastDate
}

func SetLastChange(actualDate string) bool {

	conn := postgres.GetPool()

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

func GetHeader(header string) string {

	conn := postgres.GetPool()

	var result string

	err := conn.QueryRow(context.Background(), "SELECT id_header FROM headers WHERE nm_header = $1", header).Scan(&result)

	if err != nil {
		fmt.Println("Erro ao executar consulta(Header): ", err)
		return ""
	}

	return result
}
func InsertHeader(header string) string {

	conn := postgres.GetPool()

	idHeader := ""

	err := conn.QueryRow(context.Background(), "INSERT INTO headers (nm_header) VALUES ($1) RETURNING id_header", header).Scan(&idHeader)

	if err != nil {
		fmt.Println("Erro ao executar consulta:", err)
	}

	return idHeader
}

func UpdateData(berth, header, value string) {

	conn := postgres.GetPool()

	_, err := conn.Exec(context.Background(), "UPDATE data SET value = $1 WHERE cod_header = $2 AND cod_berth = $3", value, GetHeader(header), GetHeader(berth))

	if err != nil {
		fmt.Println("Erro ao executar update:", err)
		return
	}

	fmt.Println("Dado atualizado:", berth, header, value)

}

func InsertData(berth, header, value string) {
	conn := postgres.GetPool()

	idBerth, idHeader := "", ""

	err := conn.QueryRow(context.Background(), "SELECT id_header FROM headers WHERE nm_header = $1", berth).Scan(&idBerth)
	if err != nil {
		fmt.Println("Erro ao buscar idBerth:", err)
		idBerth = InsertHeader(berth)
	}

	err = conn.QueryRow(context.Background(), "SELECT id_header FROM headers WHERE nm_header = $1", header).Scan(&idHeader)
	if err != nil {
		fmt.Println("Erro ao buscar idHeader:", err)
		idHeader = InsertHeader(header)
	}

	_, err = conn.Exec(context.Background(), "INSERT INTO data (cod_berth, cod_header, value) VALUES ($1, $2, $3)", idBerth, idHeader, value)

	if err != nil {
		fmt.Println("Erro ao executar consulta:", err)
	}
}
func GetLastData() map[string]map[string]string {
	conn := postgres.GetPool()

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

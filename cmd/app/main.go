package main

import (
	"codesp_data/internal/repository/codesp"
	"codesp_data/internal/repository/db"
	"codesp_data/internal/repository/notification"
	"codesp_data/pkg/postgres"
	"fmt"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Erro ao carregar o arquivo .env")
	}
}

func compareData(newData, oldData map[string]map[string]string) {
	changesText := ""
	for key, value := range newData {
		if _, exists := oldData[key]; !exists {
			fmt.Println("Novo dado:", key, value)
			for k, v := range value {
				db.InsertData(key, k, v)
			}
		} else {
			for k, v := range value {
				if oldData[key][k] != v {
					changesText += fmt.Sprintf("<strong>%s</strong>: %s -> De %s para %s<br>", key, k, oldData[key][k], v)
					db.UpdateData(key, k, v)
				}
			}
		}
	}

	if changesText != "" {
		notification.SendEmail(changesText)
	}
}

func main() {

	defer postgres.ClosePool()

	codesp := codesp.NewCodesp()
	if codesp == nil {
		return
	}

	lastChange := db.GetLastChange()

	if codesp.GetDate() == lastChange {
		fmt.Println("Nenhuma alteração encontrada")
		return
	}

	newData := codesp.ExtractTable()
	oldData := db.GetLastData()

	if newData == nil || oldData == nil {
		fmt.Println("Erro ao extrair dados")
		return
	}

	for key := range newData {
		if db.GetHeader(key) == "" {
			db.InsertHeader(key)
			fmt.Println("Novo header:", key)
		}
	}

	compareData(newData, oldData)

	db.SetLastChange(codesp.GetDate())
}

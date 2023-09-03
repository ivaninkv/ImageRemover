package output

import (
	"fmt"
	"os"
)

func WriteToTXT(fileName string, data map[string]bool) {
	// Открываем файл для записи
	file, err := os.Create(fileName)
	defer file.Close()

	// Записываем ключи в файл
	if err == nil {
		for key := range data {
			fmt.Fprintln(file, key)
		}
	}
}

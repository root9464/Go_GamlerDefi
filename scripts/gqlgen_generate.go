package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Использование: go run scripts/generate.go [модуль1] [модуль2] ...")
		fmt.Println("Пример: go run scripts/generate.go test auth user")
		os.Exit(1)
	}

	modules := os.Args[1:]
	fmt.Printf("Генерация кода для модулей: %s\n", strings.Join(modules, ", "))

	for _, module := range modules {
		fmt.Printf("Генерация для модуля %s...\n", module)
		tempConfig := fmt.Sprintf("gqlgen_%s.yml", module)

		configData, err := os.ReadFile("gqlgen.yml")
		if err != nil {
			fmt.Printf("Ошибка чтения файла конфигурации: %v\n", err)
			os.Exit(1)
		}

		configData = []byte(strings.ReplaceAll(string(configData), "{module}", module))
		err = os.WriteFile(tempConfig, configData, 0644)
		if err != nil {
			fmt.Printf("Ошибка записи временного файла конфигурации: %v\n", err)
			os.Exit(1)
		}

		cmd := exec.Command("go", "run", "github.com/99designs/gqlgen", "generate", "--config", tempConfig)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()

		os.Remove(tempConfig)

		if err != nil {
			fmt.Printf("Ошибка при генерации кода для модуля %s: %v\n", module, err)
			os.Exit(1)
		}

		fmt.Printf("Генерация для модуля %s завершена успешно\n", module)
	}

	fmt.Println("Генерация кода завершена для всех указанных модулей")
}

package main

import (
	"fmt"
	"sso/internal/config"
)

func main() {
	config := config.Parse()
	fmt.Println(config)

	//TODO: setup logger

	//TODO: подключиться к бд

	//TODO: инициализировать приложение
}

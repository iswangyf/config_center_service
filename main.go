package main

import (
	"fmt"

	"github.com/iswangyf/config_center_service/internal/dbinit"
)

func main() {
	fmt.Println("Config Center Service started.")
	dbinit.InitConfig()
	dbinit.InitDB()
}

package config

import (
	"fmt"
	"os"
)

func RunHttp() string {
	addr := os.Getenv("SERVER_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	fmt.Println("Server running at http://localhost" + addr)
	return addr
}

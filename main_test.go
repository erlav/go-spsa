package main

import (
	"log"

	pg "github.com/erlav/go-spsa/monitor"
	dotenv "github.com/joho/godotenv"
)

func main() {
	err := dotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	s1 := make([]string, 0)
	d1 := make(map[string]string)
	m, _ := pg.New("myjob", s1, d1)
	d2 := make(map[string]string)
	d2["fn"] = "ayuda.txt"
	if err := m.PushMetric("test_pending_file", 1, d2); err != nil {
		log.Fatal("No se pudo enviar:", err)
	}
}

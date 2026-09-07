package main

import (
	"log"
	"os"
	"os/exec"

	"github.com/irpanzy/Task-Forge/internal/config"
)

func main() {
	config.LoadEnv()

	if len(os.Args) < 2 {
		log.Fatal("Perintah migrasi dibutuhkan. Contoh: up, down, version, force")
	}

	if config.AppConfig.DatabaseURL == "" {
		log.Fatal("DATABASE_URL belum diatur di .env")
	}

	args := []string{
		"-path", "database/migrations",
		"-database", config.AppConfig.DatabaseURL,
	}
	args = append(args, os.Args[1:]...)

	cmd := exec.Command("migrate", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("Eksekusi migrate gagal: %v", err)
	}
}

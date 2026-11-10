package migrationutil

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"api-template/configs"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const MigrationPath = "file://migrations"

func HandleCLI() {
	if len(os.Args) > 1 {
		args := os.Args[1:]
		cmd := args[0]

		switch cmd {
		case "migrate", "migrate:refresh", "make:migration":
			HandleMigrationCommand(args)

		case "db:seed", "create:seeder":
			HandleSeederCommand(args)

		case "fix:dirty":
			FixDirtyMigration()

		default:
			log.Println("⚠️  Command tidak dikenali:", cmd)
			log.Println("Usage:")
			log.Println("  go run cmd/main.go migrate")
			log.Println("  go run cmd/main.go migrate:refresh")
			log.Println("  go run cmd/main.go make:migration <name>")
			log.Println("  go run cmd/main.go db:seed")
			log.Println("  go run cmd/main.go create:seeder <name>")
			log.Println("  go run cmd/main.go fix:dirty")
		}
		os.Exit(0)
	}
}

func RunMigration() {
	cfg, err := configs.Load()
	if err != nil {
		log.Fatal("Gagal load config:", err)
	}

	db, err := configs.NewDB(cfg.DB)
	if err != nil {
		log.Fatal("Gagal koneksi database:", err)
	}
	defer db.Close()

	sqlDB := db.DB

	driver, err := mysql.WithInstance(sqlDB, &mysql.Config{})
	if err != nil {
		log.Fatal("Gagal buat driver MySQL:", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		MigrationPath,
		"mysql",
		driver,
	)
	if err != nil {
		log.Fatal("Gagal inisialisasi migrate:", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Migration error:", err)
	}

	log.Println("✅ Migration sukses dijalankan!")
}

func RefreshMigration() {
	cfg, err := configs.Load()
	if err != nil {
		log.Fatal("Gagal load config:", err)
	}

	db, err := configs.NewDB(cfg.DB)
	if err != nil {
		log.Fatal("Gagal koneksi database:", err)
	}
	defer db.Close()

	sqlDB := db.DB

	if _, err := sqlDB.Exec("SET FOREIGN KEY_CHECKS=0"); err != nil {
		log.Fatal("Gagal disable foreign key checks:", err)
	}

	driver, err := mysql.WithInstance(sqlDB, &mysql.Config{})
	if err != nil {
		log.Fatal("Gagal membuat driver migration:", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		MigrationPath,
		"mysql",
		driver,
	)
	if err != nil {
		log.Fatal("Gagal inisialisasi migration:", err)
	}

	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		log.Println("Warning saat rollback:", err)
	}

	if _, err := sqlDB.Exec("SET FOREIGN KEY_CHECKS=1"); err != nil {
		log.Fatal("Gagal enable foreign key checks:", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Gagal menjalankan migration:", err)
	}

	log.Println("✅ Migration refresh sukses!")
}

func extractTableName(migrationName string) string {
	name := migrationName
	prefixes := []string{"create_", "drop_", "add_", "remove_", "alter_"}
	for _, prefix := range prefixes {
		name = strings.TrimPrefix(name, prefix)
	}

	name = strings.TrimSuffix(name, "_table")

	if strings.Contains(name, "_to_") {
		parts := strings.Split(name, "_to_")
		if len(parts) > 1 {
			return parts[len(parts)-1]
		}
	}

	if strings.Contains(name, "_from_") {
		parts := strings.Split(name, "_from_")
		if len(parts) > 1 {
			return parts[len(parts)-1]
		}
	}

	return name
}

func MakeMigration(name string) {
	if name == "" {
		log.Fatal("Nama migration wajib diisi, contoh: go run cmd/main.go make:migration create_users_table")
	}

	timestamp := time.Now().Format("20060102150405")
	base := fmt.Sprintf("%s_%s", timestamp, name)

	upFile := filepath.Join("migrations", base+".up.sql")
	downFile := filepath.Join("migrations", base+".down.sql")

	os.MkdirAll("migrations", os.ModePerm)

	tableName := extractTableName(name)

	if err := os.WriteFile(upFile, []byte("-- "+name+" UP SQL\n"), 0644); err != nil {
		log.Fatal(err)
	}

	downContent := fmt.Sprintf("-- %s DOWN SQL\nDROP TABLE IF EXISTS %s;\n", name, tableName)
	if err := os.WriteFile(downFile, []byte(downContent), 0644); err != nil {
		log.Fatal(err)
	}

	log.Println("Migration file dibuat:")
	log.Println(" - " + upFile)
	log.Println(" - " + downFile)
}

func HandleMigrationCommand(args []string) {
	if len(args) < 1 {
		log.Fatal("Gunakan: migrate | migrate:refresh | make:migration <name>")
	}

	cmd := args[0]
	switch cmd {
	case "migrate":
		RunMigration()
	case "migrate:refresh":
		RefreshMigration()
	case "make:migration":
		if len(args) < 2 {
			log.Fatal("Gunakan: make:migration <nama_migration>")
		}
		MakeMigration(strings.Join(args[1:], "_"))
	default:
		log.Fatalf("Command %s tidak dikenali", cmd)
	}
}

func FixDirtyMigration() {
	cfg, err := configs.Load()
	if err != nil {
		log.Fatal("Gagal load config:", err)
	}

	db, err := configs.NewDB(cfg.DB)
	if err != nil {
		log.Fatal("Gagal koneksi database:", err)
	}
	defer db.Close()

	sqlDB := db.DB

	var version uint
	var dirty bool

	driver, err := mysql.WithInstance(sqlDB, &mysql.Config{})
	if err != nil {
		log.Fatal("Gagal membuat driver migration:", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		MigrationPath,
		"mysql",
		driver,
	)
	if err != nil {
		log.Fatal("Gagal inisialisasi migration:", err)
	}

	version, dirty, err = m.Version()
	if err != nil {
		if err == migrate.ErrNilVersion {
			log.Println("Belum ada migration yang dijalankan.")
			return
		}
		log.Fatal("Gagal membaca versi migration:", err)
	}

	if !dirty {
		log.Printf("Migration versi %d tidak dalam status dirty.\n", version)
		return
	}

	if err := m.Force(int(version)); err != nil {
		log.Fatal("Gagal memperbaiki dirty migration:", err)
	}

	log.Printf("✅ Dirty migration pada versi %d berhasil dibersihkan.\n", version)
}


package utils

import (
	"database/sql"
	"fmt"
	"log"
	"modify-framework-mod/database"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

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
			log.Println("⚠️ Command tidak dikenali:", cmd)
			log.Println("Usage:")
			log.Println("  go run main.go migrate")
			log.Println("  go run main.go migrate:refresh")
			log.Println("  go run main.go make:migration <name>")
			log.Println("  go run main.go db:seed")
			log.Println("  go run main.go create:seeder <name>")
			log.Println("  go run main.go fix:dirty")
		}
		os.Exit(0)
	}
}

func RunMigration() {
	sqlDB, err := database.DB.DB()
	if err != nil {
		log.Fatal("Gagal ambil instance sql.DB:", err)
	}

	driver, err := mysql.WithInstance(sqlDB, &mysql.Config{})
	if err != nil {
		log.Fatal("Gagal buat driver MySQL:", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "mysql", driver)
	if err != nil {
		log.Fatal("Gagal inisialisasi migrate:", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Migration error:", err)
	}
	log.Println("Migration sukses dijalankan!")
}

func RefreshMigration() {
	sqlDB, err := database.DB.DB()
	if err != nil {
		log.Fatal("Gagal ambil instance sql.DB:", err)
	}

	if _, err := sqlDB.Exec("SET FOREIGN_KEY_CHECKS=0;"); err != nil {
		log.Fatal("Gagal disable foreign key checks:", err)
	}

	driver, err := mysql.WithInstance(sqlDB, &mysql.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "mysql", driver)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		log.Println("Warning saat drop table:", err)
	}

	if _, err := sqlDB.Exec("SET FOREIGN_KEY_CHECKS=1;"); err != nil {
		log.Fatal("Gagal enable foreign key checks:", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	log.Println("Migration refresh sukses!")
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
		log.Fatal("Nama migration wajib diisi, contoh: go run main.go make:migration create_users_table")
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

	downContent := fmt.Sprintf("-- %s DOWN SQL\nDROP TABLE IF EXISTS %s;", name, tableName)
	if err := os.WriteFile(downFile, []byte(downContent), 0644); err != nil {
		log.Fatal(err)
	}

	log.Println("Migration file dibuat:")
	log.Println(" - " + upFile)
	log.Println(" - " + downFile)
}

func HandleMigrationCommand(args []string) {
	if len(args) < 1 {
		log.Fatal("Gunakan: migrate | migrate:rollback | migrate:refresh | make:migration <name>")
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
	sqlDB, err := database.DB.DB()
	if err != nil {
		log.Fatal("Gagal ambil instance sql.DB:", err)
	}

	var version int
	var dirty bool

	err = sqlDB.QueryRow("SELECT version, dirty FROM schema_migrations LIMIT 1").Scan(&version, &dirty)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Tidak ada record di schema_migrations — database bersih.")
			return
		}
		log.Fatal("Gagal membaca schema_migrations:", err)
	}

	if !dirty {
		log.Println("Database tidak dalam status dirty. Tidak perlu diperbaiki.")
		return
	}

	// Tandai dirty jadi false tanpa ubah versi
	_, err = sqlDB.Exec("UPDATE schema_migrations SET dirty = 0")
	if err != nil {
		log.Fatal("Gagal memperbaiki status dirty:", err)
	}

	log.Printf("✅ Status dirty pada migration versi %d berhasil dibersihkan!\n", version)
}

// SEEDER ===================

func replaceBlobPlaceholders(query string) string {
	re := regexp.MustCompile(`\{\{BLOB:([^\}]+)\}\}`)

	return re.ReplaceAllStringFunc(query, func(match string) string {
		matches := re.FindStringSubmatch(match)
		if len(matches) < 2 {
			return match
		}

		fileName := matches[1]
		filePath := filepath.Join("seeders/assets", fileName)

		data, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("⚠️  Gagal baca file blob %s: %v", filePath, err)
			return "NULL"
		}

		return fmt.Sprintf("x'%x'", data)
	})
}

func RunSeeder() {
	db := database.DB
	sPath := "seeders"

	files, err := os.ReadDir(sPath)
	if err != nil {
		log.Fatalf("Gagal membaca folder seeders: %v", err)
	}

	var sqlFiles []string
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".sql") {
			sqlFiles = append(sqlFiles, filepath.Join(sPath, f.Name()))
		}
	}
	sort.Strings(sqlFiles)

	if len(sqlFiles) == 0 {
		log.Println("⚠️  Tidak ada file seeder ditemukan.")
		return
	}

	for _, file := range sqlFiles {
		content, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("Gagal baca file %s: %v", file, err)
		}

		queries := strings.Split(string(content), ";")
		for _, q := range queries {
			q = strings.TrimSpace(q)
			if q == "" {
				continue
			}

			q = replaceBlobPlaceholders(q)

			// pakai GORM
			if err := db.Exec(q).Error; err != nil {
				log.Fatalf("❌ Error eksekusi %s: %v", file, err)
			}
		}
		log.Printf("✅ Seeder dijalankan: %s\n", file)
	}

	log.Println("✅ Semua seeder sukses dijalankan!")
}

func MakeSeeder(name string) {
	if name == "" {
		log.Fatal("Nama seeder wajib diisi, contoh: go run main.go create:seeder users")
	}

	sPath := "seeders"

	timestamp := time.Now().Format("20060102150405")
	base := fmt.Sprintf("%s_%s.sql", timestamp, name)

	file := filepath.Join(sPath, base)

	os.MkdirAll(sPath, os.ModePerm)

	if err := os.WriteFile(file, []byte("-- "+name+" SEED SQL\n"), 0644); err != nil {
		log.Fatal(err)
	}

	log.Println("📄 Seeder file dibuat:", file)
}

func HandleSeederCommand(args []string) {
	if len(args) < 1 {
		log.Fatal("Gunakan: db:seed | create:seeder <name>")
	}

	cmd := args[0]
	switch cmd {
	case "db:seed":
		RunSeeder()
	case "create:seeder":
		if len(args) < 2 {
			log.Fatal("Gunakan: create:seeder <nama_seeder>")
		}
		MakeSeeder(strings.Join(args[1:], "_"))
	default:
		log.Fatalf("Command %s tidak dikenali", cmd)
	}
}

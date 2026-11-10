package migrationutil

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"api-template/configs"
	"api-template/pkg/hash"

	"github.com/jmoiron/sqlx"
)

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

func updateHashIDs(db *sqlx.DB, table string) error {
	type Row struct {
		ID int64 `db:"id"`
	}

	var rows []Row

	query := fmt.Sprintf(`
		SELECT id
		FROM %s
		WHERE hash_id IS NULL
		   OR hash_id = ''
	`, table)

	if err := db.Select(&rows, query); err != nil {
		return err
	}

	for _, row := range rows {
		hashID, err := hash.EncodeID(row.ID)
		if err != nil {
			return err
		}

		updateQuery := fmt.Sprintf(`
			UPDATE %s
			SET hash_id = ?
			WHERE id = ?
		`, table)

		if _, err := db.Exec(updateQuery, hashID, row.ID); err != nil {
			return err
		}
	}

	log.Printf("✅ %d hash_id berhasil diupdate pada table %s", len(rows), table)

	return nil
}

func RunSeeder() {
	cfg, err := configs.Load()
	if err != nil {
		log.Fatalf("Gagal load config: %v", err)
	}

	db, err := configs.NewDB(cfg.DB)
	if err != nil {
		log.Fatalf("Gagal koneksi database: %v", err)
	}
	defer db.Close()

	sPath := "seeders"
	if _, err := os.Stat(sPath); os.IsNotExist(err) {
		log.Println("⚠️  Folder seeder tidak ditemukan. Tidak ada data yang di-seed.")
		return
	}

	files, err := os.ReadDir(sPath)
	if err != nil {
		log.Fatalf("Gagal membaca folder seeder: %v", err)
	}

	var sqlFiles []string
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".sql") {
			sqlFiles = append(sqlFiles, filepath.Join(sPath, f.Name()))
		}
	}

	sort.Strings(sqlFiles)

	if len(sqlFiles) == 0 {
		log.Println("⚠️  Tidak ada file seeder ditemukan di folder seeders/.")
		return
	}

	for _, file := range sqlFiles {
		content, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("Gagal membaca file %s: %v", file, err)
		}

		queries := strings.Split(string(content), ";")

		for _, q := range queries {
			q = strings.TrimSpace(q)
			if q == "" {
				continue
			}

			q = replaceBlobPlaceholders(q)

			if _, err := db.Exec(q); err != nil {
				log.Fatalf("❌ Error eksekusi %s: %v", file, err)
			}
		}

		log.Printf("✅ Seeder dijalankan: %s", file)
	}

	log.Println("✅ Semua seeder sukses dijalankan!")
}

func MakeSeeder(name string) {
	if name == "" {
		log.Fatal("Nama seeder wajib diisi, contoh: go run cmd/main.go create:seeder users")
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

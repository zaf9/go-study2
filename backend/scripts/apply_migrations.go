package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/os/gctx"
)

func main() {
	ctx := gctx.New()
	dbPath := "./data/gostudy.db"
	migrationDir := "./internal/infra/migrations"

	// Ensure db dir exists
	os.MkdirAll(filepath.Dir(dbPath), 0755)

	absPath, _ := filepath.Abs(dbPath)
	link := fmt.Sprintf("sqlite::@file(%s)", filepath.ToSlash(absPath))

	db, err := gdb.New(gdb.ConfigNode{
		Type: "sqlite",
		Link: link,
	})
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	files, err := ioutil.ReadDir(migrationDir)
	if err != nil {
		log.Fatalf("Failed to read migrations dir: %v", err)
	}

	var sqlFiles []string
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".sql") {
			// We only want to run 012 and 013 for this phase if they are specific to T001/T002
			// But usually we run all. Given 011 is already there, we should be careful.
			// For this task, let's just run 012 and 013.
			if strings.HasPrefix(f.Name(), "012") || strings.HasPrefix(f.Name(), "013") {
				sqlFiles = append(sqlFiles, f.Name())
			}
		}
	}
	sort.Strings(sqlFiles)

	for _, f := range sqlFiles {
		fmt.Printf("Applying migration: %s\n", f)
		content, err := ioutil.ReadFile(filepath.Join(migrationDir, f))
		if err != nil {
			log.Fatalf("Failed to read %s: %v", f, err)
		}

		// Quick split by semicolon for execution
		queries := strings.Split(string(content), ";")
		for _, q := range queries {
			q = strings.TrimSpace(q)
			if q == "" || strings.HasPrefix(q, "--") || strings.HasPrefix(q, "BEGIN") || strings.HasPrefix(q, "COMMIT") {
				continue
			}
			if _, err := db.Exec(ctx, q); err != nil {
				log.Fatalf("Failed to execute query in %s: %v\nQuery: %s", f, err, q)
			}
		}
	}

	fmt.Println("Migrations applied successfully.")

	// Verify tables
	tables, err := db.Tables(ctx)
	if err != nil {
		log.Fatalf("Failed to list tables: %v", err)
	}
	fmt.Printf("Current tables: %v\n", tables)

	for _, t := range []string{"quiz_sessions", "quiz_attempts"} {
		fields, err := db.TableFields(ctx, t)
		if err != nil {
			log.Fatalf("Failed to get fields for %s: %v", t, err)
		}
		fmt.Printf("Fields for %s:\n", t)
		for name, f := range fields {
			fmt.Printf("  - %s: %s\n", name, f.Type)
		}
	}
}

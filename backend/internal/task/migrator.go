package task

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// RunMigrations executes all SQL files found inside dir alphabetically.
func RunMigrations(ctx context.Context, pool *pgxpool.Pool, dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read migrations: %w", err)
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		if err := applyFile(ctx, pool, path); err != nil {
			return fmt.Errorf("apply %s: %w", entry.Name(), err)
		}
	}
	return nil
}

func applyFile(ctx context.Context, pool *pgxpool.Pool, path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	statements := extractUpStatements(string(content))
	if statements == "" {
		return nil
	}
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, statements); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

// RunMigrationsFromFS allows executing migrations from an embed FS; used by tests.
func RunMigrationsFromFS(ctx context.Context, pool *pgxpool.Pool, filesystem fs.FS, dir string) error {
	entries, err := fs.ReadDir(filesystem, dir)
	if err != nil {
		return err
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		content, err := fs.ReadFile(filesystem, filepath.Join(dir, entry.Name()))
		if err != nil {
			return err
		}
		statements := extractUpStatements(string(content))
		if statements == "" {
			continue
		}
		if _, err := pool.Exec(ctx, statements); err != nil {
			return err
		}
	}
	return nil
}

// extractUpStatements trims away goose-style Down sections so we don't accidentally run rollbacks on startup.
func extractUpStatements(content string) string {
	lines := strings.Split(content, "\n")
	inDown := false
	inUp := false
	hasDirective := false
	var filtered []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(trimmed, "-- +goose Up"):
			inUp = true
			hasDirective = true
			continue
		case strings.HasPrefix(trimmed, "-- +goose Down"):
			if hasDirective {
				inDown = true
			}
			continue
		}
		if inDown {
			continue
		}
		if hasDirective && !inUp {
			continue
		}
		filtered = append(filtered, line)
	}
	return strings.TrimSpace(strings.Join(filtered, "\n"))
}

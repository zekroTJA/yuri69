package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type migrationFunc func(tx *sql.Tx) error

func (t *Postgres) Migrate() error {
	currentMigration := -1

	err := t.db.QueryRow(`
		SELECT "id" FROM migrations
		ORDER BY "id" DESC
		LIMIT 1
	`).Scan(&currentMigration)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if currentMigration >= len(migrationFuncs) {
		return fmt.Errorf(
			"Database is on an invalid migration level. "+
				"One reason for that could be that you are using "+
				"an older version of yuri than the database has "+
				"been initialized with. "+
				"(currentMigration: %d, latestMigration: %d)",
			currentMigration, len(migrationFuncs)-1)
	}

	if currentMigration == len(migrationFuncs)-1 {
		logrus.Info("Database is up to date")
		return nil
	}

	for i := currentMigration + 1; i < len(migrationFuncs); i++ {
		err = t.tx(func(tx *sql.Tx) error {
			logrus.WithField("migration", i).Info("Applying migration ...")
			mf := migrationFuncs[i]
			if err := mf(tx); err != nil {
				return err
			}
			_, err := tx.Exec(`
				INSERT INTO migrations ("id", "timestamp")
				VALUES ($1, $2)
			`, i, time.Now())
			return err
		})
		if err != nil {
			return err
		}
	}

	logrus.Info("All migrations successfully applied")

	return nil
}

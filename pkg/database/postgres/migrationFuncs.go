package postgres

import "database/sql"

var migrationFuncs = []migrationFunc{
	migration_01,
}

func migration_01(tx *sql.Tx) error {
	_, err := tx.Exec(`
		ALTER TABLE "twitchsettings"
		ADD COLUMN "blocklist" TEXT NOT NULL DEFAULT '';
	`)
	return err
}

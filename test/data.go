package test

import (
	"database/sql"
	"os"
)

// SetupTestDatabase creates the test
func SetupTestDatabase(dbName string) error {
	os.Remove(dbName)

	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStmt := `CREATE TABLE "projects" ("id" integer primary key autoincrement,"created_at" datetime,"updated_at" datetime,"deleted_at" datetime,"name" varchar(255),"description" varchar(255) );`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	sqlStmt = `CREATE INDEX idx_projects_deleted_at ON "projects"(deleted_at);`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	sqlStmt = `CREATE UNIQUE INDEX uix_projects_email ON "projects"("name");`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	sqlStmt = `INSERT INTO "projects" ("created_at","updated_at","deleted_at","name","description") VALUES ('2018-05-06 20:42:37','2018-05-06 20:42:37',NULL,'testproject123','test project description goes here');`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	return nil
}

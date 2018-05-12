package dao

import (
	"database/sql"
	"os"
	"testing"

	"github.com/bojand/ghz-web/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stretchr/testify/assert"
)

const dbName = "../test/project_test.db"

func createTestData() error {
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

func TestProjectService_FindByID(t *testing.T) {
	defer os.Remove(dbName)

	err := createTestData()
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	db, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	defer db.Close()

	dao := ProjectService{DB: db}

	t.Run("test existing", func(t *testing.T) {
		p := model.Project{}
		err := dao.FindByID(1, &p)

		assert.NoError(t, err)
		assert.Equal(t, uint(1), p.ID)
		assert.Equal(t, "testproject123", p.Name)
		assert.Equal(t, "test project description goes here", p.Description)
		assert.NotNil(t, p.CreatedAt)
		assert.NotNil(t, p.UpdatedAt)
		assert.Nil(t, p.DeletedAt)
	})

	t.Run("test not found", func(t *testing.T) {
		p := model.Project{}
		err := dao.FindByID(2, &p)

		assert.Error(t, err)
		assert.Equal(t, uint(0), p.ID)
		assert.Equal(t, "", p.Name)
		assert.Equal(t, "", p.Description)
	})
}

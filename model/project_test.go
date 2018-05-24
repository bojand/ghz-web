package model

import (
	"os"
	"testing"

	"github.com/bojand/ghz-web/test"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stretchr/testify/assert"
)

const dbName = "../test/project_test.db"

func TestProjectService_FindByID(t *testing.T) {
	defer os.Remove(dbName)

	err := test.SetupTestProjectDatabase(dbName)
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
		p, err := dao.FindByID(1)

		assert.NoError(t, err)
		assert.Equal(t, uint(1), p.ID)
		assert.Equal(t, "testproject123", p.Name)
		assert.Equal(t, "test project description goes here", p.Description)
		assert.NotNil(t, p.CreatedAt)
		assert.NotNil(t, p.UpdatedAt)
		assert.Nil(t, p.DeletedAt)
	})

	t.Run("test not found", func(t *testing.T) {
		p, err := dao.FindByID(2)

		assert.Error(t, err)
		assert.Nil(t, p)
	})
}

func TestProjectService_FindByName(t *testing.T) {
	defer os.Remove(dbName)

	err := test.SetupTestProjectDatabase(dbName)
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
		p, err := dao.FindByName("testproject123")

		assert.NoError(t, err)
		assert.Equal(t, uint(1), p.ID)
		assert.Equal(t, "testproject123", p.Name)
		assert.Equal(t, "test project description goes here", p.Description)
		assert.NotNil(t, p.CreatedAt)
		assert.NotNil(t, p.UpdatedAt)
		assert.Nil(t, p.DeletedAt)
	})

	t.Run("test not found", func(t *testing.T) {
		p, err := dao.FindByName("testproject999")

		assert.Error(t, err)
		assert.Nil(t, p)
	})
}

func TestProjectService_Create(t *testing.T) {
	defer os.Remove(dbName)

	err := test.SetupTestProjectDatabase(dbName)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	db, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	defer db.Close()

	dao := ProjectService{DB: db}

	t.Run("test new", func(t *testing.T) {
		p := Project{
			Name:        "TestProj111 ",
			Description: "Test Description Asdf ",
		}
		err := dao.Create(&p)

		assert.NoError(t, err)
		assert.NotZero(t, p.ID)
		assert.Equal(t, "testproj111", p.Name)
		assert.Equal(t, "Test Description Asdf", p.Description)
		assert.NotNil(t, p.CreatedAt)
		assert.NotNil(t, p.UpdatedAt)
		assert.Nil(t, p.DeletedAt)

		p2, err := dao.FindByID(p.ID)

		assert.NoError(t, err)
		assert.Equal(t, "testproj111", p2.Name)
		assert.Equal(t, "Test Description Asdf", p2.Description)
		assert.NotNil(t, p2.CreatedAt)
		assert.NotNil(t, p2.UpdatedAt)
		assert.Nil(t, p2.DeletedAt)
	})

	t.Run("test new with empty name", func(t *testing.T) {
		p := Project{
			Description: "Test Description Asdf 2",
		}
		err := dao.Create(&p)

		assert.NoError(t, err)
		assert.NotZero(t, p.ID)
		assert.NotEmpty(t, p.Name)
		assert.Equal(t, "Test Description Asdf 2", p.Description)
		assert.NotNil(t, p.CreatedAt)
		assert.NotNil(t, p.UpdatedAt)
		assert.Nil(t, p.DeletedAt)

		p2, err := dao.FindByID(p.ID)

		assert.NoError(t, err)
		assert.Equal(t, p.Name, p2.Name)
		assert.Equal(t, "Test Description Asdf 2", p2.Description)
		assert.NotNil(t, p2.CreatedAt)
		assert.NotNil(t, p2.UpdatedAt)
		assert.Nil(t, p2.DeletedAt)
	})

	t.Run("test new with ID", func(t *testing.T) {
		p := Project{
			Name:        " FooProject ",
			Description: " Bar Desc ",
		}
		p.ID = 123

		err := dao.Create(&p)

		assert.NoError(t, err)
		assert.Equal(t, uint(123), p.ID)
		assert.Equal(t, "fooproject", p.Name)
		assert.Equal(t, "Bar Desc", p.Description)
		assert.NotNil(t, p.CreatedAt)
		assert.NotNil(t, p.UpdatedAt)
		assert.Nil(t, p.DeletedAt)

		p2, err := dao.FindByID(p.ID)

		assert.NoError(t, err)
		assert.Equal(t, uint(123), p2.ID)
		assert.Equal(t, "fooproject", p2.Name)
		assert.Equal(t, "Bar Desc", p2.Description)
		assert.NotNil(t, p.CreatedAt)
		assert.NotNil(t, p.UpdatedAt)
		assert.Nil(t, p.DeletedAt)
	})

	t.Run("should fail with same ID", func(t *testing.T) {
		p := Project{
			Name:        "ACME",
			Description: "Lorem Ipsum",
		}
		p.ID = 123

		err := dao.Create(&p)

		assert.Error(t, err)
	})

	t.Run("should fail with same name", func(t *testing.T) {
		p := Project{
			Name:        "FooProject",
			Description: "Lorem Ipsum",
		}
		err := dao.Create(&p)

		assert.Error(t, err)
	})
}

func TestProjectService_Update(t *testing.T) {
	defer os.Remove(dbName)

	err := test.SetupTestProjectDatabase(dbName)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	db, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	defer db.Close()

	dao := ProjectService{DB: db}

	t.Run("fail with new", func(t *testing.T) {
		p := Project{
			Name:        "testproject124",
			Description: "asdf",
		}
		p.ID = 4321

		err := dao.Update(&p)

		assert.Error(t, err)
	})

	t.Run("test update existing", func(t *testing.T) {
		p := Project{
			Name:        " New Name ",
			Description: "Baz",
		}
		p.ID = uint(1)

		err := dao.Update(&p)

		assert.NoError(t, err)

		assert.NotZero(t, p.ID)
		assert.Equal(t, "newname", p.Name)
		assert.Equal(t, "Baz", p.Description)
		assert.NotNil(t, p.CreatedAt)
		assert.NotNil(t, p.UpdatedAt)
		assert.Nil(t, p.DeletedAt)
	})

	t.Run("test update existing just name", func(t *testing.T) {
		p := Project{
			Name: " New Name 2",
		}
		p.ID = uint(1)

		err := dao.Update(&p)

		assert.NoError(t, err)

		assert.NotZero(t, p.ID)
		assert.Equal(t, "newname2", p.Name)
		assert.Equal(t, "", p.Description)
		assert.NotNil(t, p.CreatedAt)
		assert.NotNil(t, p.UpdatedAt)
		assert.Nil(t, p.DeletedAt)
	})

	t.Run("test update existing no name", func(t *testing.T) {
		p := Project{
			Description: "Foo Test Bar",
		}
		p.ID = uint(1)

		err := dao.Update(&p)

		assert.NoError(t, err)

		assert.NotZero(t, p.ID)
		assert.Equal(t, "newname2", p.Name)
		assert.Equal(t, "Foo Test Bar", p.Description)
		assert.NotNil(t, p.CreatedAt)
		assert.NotNil(t, p.UpdatedAt)
		assert.Nil(t, p.DeletedAt)
	})
}

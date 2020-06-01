package database

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stretchr/testify/require"
)

// Test Model
type TestModel struct {
	ID   int64  `gorm:"column:id;primary_key" json:"id"`
	Data string `gorm:"column:data" json:"data"`
}

func (m *TestModel) TableName() string {
	return "test_model"
}

func setupTestCase(t *testing.T) (func(t *testing.T), *gorm.DB) {
	dbFile := filepath.Base("test_db.db")
	db, err := gorm.Open("sqlite3", dbFile)
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&TestModel{})

	return func(t *testing.T) {
		_ = os.Remove(dbFile)
		db.Close()
	}, db
}

func TestRepo_CRUD(t *testing.T) {
	teardownTestCase, db := setupTestCase(t)
	defer teardownTestCase(t)

	model := &TestModel{Data: "123"}
	require.Empty(t, model.ID)

	r := NewGormRepo(db)
	err := r.Create(model)
	require.NoError(t, err)
	require.NotNil(t, model.ID)

	m := &TestModel{}
	err = r.FindOne(TestModel{ID: model.ID}, m)
	require.NoError(t, err)
	require.Equal(t, model.ID, m.ID)
	require.Equal(t, model.Data, m.Data)

	model.Data = "555"
	err = r.Update(model)
	require.NoError(t, err)
	m = &TestModel{}
	err = r.FindOne(TestModel{ID: model.ID}, m)
	require.NoError(t, err)
	require.Equal(t, model.ID, m.ID)
	require.Equal(t, model.Data, m.Data)

	err = r.Delete(model)
	require.NoError(t, err)
	require.Equal(t, model.ID, m.ID)
	require.Equal(t, model.Data, m.Data)

	getdb := r.GetDB()
	require.NotEmpty(t, getdb)

	r.InitDB(db)
	getNewDB := r.GetDB()
	require.NotEmpty(t, getNewDB)
}

package apiserver

import (
	"testing"

	"github.com/VitalyCone/websocket-messenger/internal/app/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestDatabase(t *testing.T, databaseURL string) (*gorm.DB, func(...interface{})){
	t.Helper()

	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	
	sqlDB, err := db.DB()
	if err != nil{
		t.Fatal(err)
	}
	
	if err := sqlDB.Ping(); err != nil {
		t.Fatal(err)
	}
	
	
	if err := db.AutoMigrate(model.User{}); err != nil{
		t.Fatal(err)
	}
	
	return db, func(tables ...interface{}){
		// if len(tables) > 0{
		// 	db.Exec(fmt.Sprintf("TRUNCATE %s CASCADE", strings.Join(tables, ", ")))
		for _, table := range tables {
			db.Where("1 = 1").Delete(table)
		}
	}
}
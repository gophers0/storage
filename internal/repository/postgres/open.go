package postgres

import (
	"fmt"
	"github.com/jinzhu/gorm"
	gaarx "github.com/zergu1ar/Gaarx"
)

type Repo struct {
	DB *gorm.DB
}

func GetConnString(user, pass, host, port, dbName string) string {
	return fmt.Sprintf(`
host=%s   
port=%s
dbname=%s
user=%s
password=%s
sslmode=%s
`,
		host,
		port,
		dbName,
		user,
		pass,
		"disable",
	)
}

// Include Database
func WithDatabase(dsn string, dbType string, entities ...interface{}) gaarx.Option {
	return func(app *gaarx.App) error {
		db, err := gorm.Open(
			dbType,
			dsn,
		)
		if err != nil {
			app.GetLog().Fatal(err)
		}
		db.SetLogger(app.GetLog())
		db.Set("gorm:table_options", "CHARSET=utf8")
		for _, e := range entities {
			db.AutoMigrate(e)
		}

		db.DB().SetMaxOpenConns(90)
		db.DB().SetMaxIdleConns(90)

		repo, err := newPostgresRepo(db)
		if err != nil {
			panic(err)
		}

		app.SetDatabase(repo)
		return nil
	}
}

// gorm open
func newPostgresRepo(db *gorm.DB) (*Repo, error) {
	return &Repo{
		DB: db,
	}, nil
}

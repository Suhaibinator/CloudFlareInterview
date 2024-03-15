package main

import (
	"ci/config"
	db "ci/dbadapters"
)

func main() {

	dbConfig := db.DBConfig{
		DBType:   db.SQLite,
		Host:     config.Configuration.Database.PostgresHost,
		Port:     config.Configuration.Database.PostgresPort,
		Username: config.Configuration.Database.PostgresUser,
		Password: config.Configuration.Database.PostgresPassword,
		DBName:   config.Configuration.Database.PostgresDBName,
		FilePath: config.Configuration.Database.SqlitePath,
	}

	pg := db.NewDBAdapter(dbConfig.DBType)
	err := pg.Connect(dbConfig)
	if err != nil {
		panic(err)
	}
	defer pg.Close()

	err = pg.CreateTablesAndStatements()
	if err != nil {
		panic(err)
	}

	pg.InsertNewShortUrl("test", "http://www.google.com", nil)

	result, expires_at, err := pg.GetFullUrl("test")

	if err != nil {
		panic(err)
	}
	println(result, expires_at)
}

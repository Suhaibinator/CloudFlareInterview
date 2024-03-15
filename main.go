package main

import (
	"ci/config"
	db "ci/dbadapters"
)

func main() {

	pg := db.NewPostgresConnection()
	pg.Connect(config.Configuration.Database.PostgresHost, config.Configuration.Database.PostgresPort, config.Configuration.Database.PostgresUser, config.Configuration.Database.PostgresPassword, config.Configuration.Database.PostgresDBName)
	pg.CreateTablesAndStatements()
	//pg.InsertNewShortUrl("test", "http://www.google.com", nil)

	pg.PrintAllTableContents()
	result, expires_at, err := pg.GetFullUrl("test")

	if err != nil {
		panic(err)
	}
	println(result, expires_at)
}

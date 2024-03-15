package main

import (
	"ci/api"
	"ci/application"
	"ci/config"
	db "ci/dbadapters"
	"net/http"
	"os"
	"os/signal"
)

func main() {

	dbConfig := db.DBConfig{
		DBType:   db.ConvertDBType(config.Configuration.Database.Type),
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

	app := application.NewApp(pg, config.Configuration.Worker.Identifier)
	r := api.SetupRoutes(app)

	sysChannel := make(chan os.Signal, 1)
	signal.Notify(sysChannel, os.Interrupt)

	go func() {
		http.ListenAndServe(":"+config.Configuration.App.Port, r)
	}()

	<-sysChannel

}

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/SergeyTyurin/banner_rotation/configs"
	"github.com/SergeyTyurin/banner_rotation/database"
	"github.com/SergeyTyurin/banner_rotation/message_broker"
	"github.com/SergeyTyurin/banner_rotation/router"
)

func main() {
	dbConfig, err := configs.GetDBConnectionConfig("config/connection_config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	// Подключение в БД
	db := database.NewDatabase()
	closeFunc, err := db.Connect(dbConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := closeFunc(); err != nil {
			log.Fatal(err)
		}
	}()

	msgConfig, err := configs.GetMessageBrokerConfig("config/connection_config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	broker := message_broker.NewBroker()
	closeBroker, err := broker.Connect(msgConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer closeBroker()

	appConfig, err := configs.GetAppSettings("config/connection_config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	// Создание сервера с мультиплексором запросов
	muxRouter := router.NewRouter(db, broker)
	server := http.Server{
		Addr:    fmt.Sprintf("%s:%d", appConfig.Host(), appConfig.Port()),
		Handler: muxRouter.Mux(),
	}
	defer server.Close()
	log.Println("listening...")
	// Прослушивание сервера
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

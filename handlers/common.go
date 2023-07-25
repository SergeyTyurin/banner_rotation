package handlers

import (
	"github.com/SergeyTyurin/banner-rotation/database"
	"github.com/SergeyTyurin/banner-rotation/messagebroker"
)

type Handlers struct {
	db     database.Database
	broker messagebroker.MessageBroker
}

func NewHandlers(db database.Database, broker messagebroker.MessageBroker) Handlers {
	return Handlers{db, broker}
}

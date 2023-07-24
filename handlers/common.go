package handlers

import (
	"github.com/SergeyTyurin/banner_rotation/database"
	"github.com/SergeyTyurin/banner_rotation/message_broker"
)

type Handlers struct {
	db     database.Database
	broker message_broker.MessageBroker
}

func NewHandlers(db database.Database, broker message_broker.MessageBroker) Handlers {
	return Handlers{db, broker}
}

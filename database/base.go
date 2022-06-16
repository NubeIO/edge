package dbase

import (
	"fmt"
	"github.com/NubeIO/edge/service/apps"
	"gorm.io/gorm"
)

type DB struct {
	DB *gorm.DB
}

type DeleteMessage struct {
	Message string `json:"message"`
}

func deleteResponse(query *gorm.DB) (*DeleteMessage, error) {
	msg := &DeleteMessage{
		Message: fmt.Sprintf("no record found, deleted count:%d", 0),
	}
	if query.Error != nil {
		return msg, query.Error
	}
	r := query.RowsAffected
	if r == 0 {
		return msg, query.Error
	}
	msg.Message = fmt.Sprintf("deleted count:%d", query.RowsAffected)
	return msg, nil
}

func initAppService(serviceName string) (*apps.Apps, error) {
	inst := &apps.Apps{
		App: &apps.Store{
			ServiceName: serviceName,
		},
	}
	app, err := apps.New(inst)
	return app, err
}

func handelNotFound(name string) error {
	return fmt.Errorf("no %s with that id was found", name)
}

package app

import (
	"lambda-func/api"
	"lambda-func/database"
)

type App struct {
	ApiHandler api.ApiHandler
}

func NewApp() App {
	// we init the db store
	// gest passed down to api handler
	db := database.NewDynamoDBClient()
	apiHandler := api.NewApiHandler(db)

	return App {
		ApiHandler: apiHandler,
	}
}
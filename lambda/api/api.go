package api

import (
	"encoding/json"
	"fmt"
	"lambda-func/database"
	"lambda-func/types"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type ApiHandler struct {
	dbStore database.UserStore 
}

func NewApiHandler(dbStore database.UserStore) ApiHandler {
	return ApiHandler{
		dbStore: dbStore,
	}
}

func (api ApiHandler) RegisterUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse,error) {
	var registerUser types.RegisterUser
	err := json.Unmarshal([]byte(request.Body), &registerUser)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body: "error unmarshalling request",
			StatusCode: 500,
		}, err
	}

	if registerUser.Username == "" || registerUser.Password == "" {
		return events.APIGatewayProxyResponse{
			Body: "username or password is empty",
			StatusCode: 500,
		},fmt.Errorf("registerUser is empty")
	}

	// does a user with this username exists ?
	userExists,err := api.dbStore.DoesUserExist(registerUser.Username)
	if err != nil {
		// !notes
		// error gets surfaced up to the app layer
		return events.APIGatewayProxyResponse{
			Body: "error checking if user exists",
			StatusCode: http.StatusInternalServerError,
		},fmt.Errorf("error checking if user exists %w",err)
	}

	if userExists {
		return events.APIGatewayProxyResponse{
			Body: "user already exists",
			StatusCode: http.StatusConflict,
		},fmt.Errorf("user already exists")
	}

	user , err := types.NewUser(registerUser)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body: "Internal server error",
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("could not create new user %w",err)
	}

	err = api.dbStore.InsertUser(user)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body: "error inserting user",
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("error inserting user %w",err)
	}

	return events.APIGatewayProxyResponse{
		Body: "user registered",
		StatusCode: http.StatusOK,
	}, nil
}

func (api ApiHandler) LoginUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var loginRequest LoginRequest

	err := json.Unmarshal([]byte(request.Body), &loginRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "invalid request",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	user, err := api.dbStore.GetUser(loginRequest.Username)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal server error",
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	if !types.ValidatePassword(user.Password, loginRequest.Password) {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid login credentials",
			StatusCode: http.StatusUnauthorized,
		}, nil
	}

	accessToken := types.CreateToken(user)
	successMsg := fmt.Sprintf(`{"access_token": "%s"}`, accessToken)

	return events.APIGatewayProxyResponse{
		Body:       successMsg,
		StatusCode: http.StatusOK,
	}, nil

}
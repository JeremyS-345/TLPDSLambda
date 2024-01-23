package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/TLPDSLambda/ingestData/dao"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/zerolog"
    "github.com/rs/zerolog/log"

)

func router(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "GET":
		//return show(req)
	case "POST":
		return create(req)
	default:
		return clientError(http.StatusMethodNotAllowed)
	}
	return events.APIGatewayProxyResponse{}, nil
}

func create(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if req.Headers["content-type"] != "application/json" && req.Headers["Content-Type"] != "application/json" {
		return clientError(http.StatusNotAcceptable)
	}

	item := new(dao.Item)
	err := json.Unmarshal([]byte(req.Body), item)
	if err != nil {
		return clientError(http.StatusUnprocessableEntity)
	}

	fmt.Printf("%+v", item)
	err = dao.PutItem(item)
	if err != nil {
		log.Print(err)
		return serverError(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusAccepted,
	}, nil
}

func clientError(status int) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       http.StatusText(status),
	}, nil
}

func serverError(err error) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       http.StatusText(http.StatusInternalServerError),
	}, nil
}
func main() {
	lambda.Start(router)
}

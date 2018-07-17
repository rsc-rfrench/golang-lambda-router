package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"router"
)

func Help(request router.Request) (router.Response, error) {
	return router.Response{
		Body:       "Hello my friend",
		StatusCode: 200,
	}, nil
}

func main() {
	r := router.Router{}
	r.GET("/help", Help)
	lambda.Start(r.DelegateRequest)
}

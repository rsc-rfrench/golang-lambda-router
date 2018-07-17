package router

import (
	"github.com/aws/aws-lambda-go/events"
)

type Request events.APIGatewayProxyRequest
type Response events.APIGatewayProxyResponse

type Router struct {
	routes map[string]func(Request) (Response, error)
}

func (*Router) DelegateRequest(request Request) (Response, error) {
	return Response{
		Body:       request.Body,
		StatusCode: 200,
	}, nil
}

func (r *Router) GET(path string, f func(Request) (Response, error)) {
	if r.routes == nil {
		r.routes = make(map[string]func(Request) (Response, error))
	}
	r.routes[path] = f
}

func (r *Router) dumpRoutes() map[string]func(Request) (Response, error) {
	return r.routes
}

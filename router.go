package router

import (
	"github.com/aws/aws-lambda-go/events"
	"regexp"
)

type Request events.APIGatewayProxyRequest
type Response events.APIGatewayProxyResponse

type Router struct {
	routes map[string]func(Request) (Response, error)
}

func (r *Router) DelegateRequest(request Request) (Response, error) {
	route, ok := r.routes[request.Path]
	if ok {
		return route(request)
	}
	return Response{StatusCode: 404}, nil
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

func matchRoute(route string, path string) (map[string]string, bool) {
	pattern := regexp.MustCompile(route)
	matched := pattern.MatchString(path)
	template := "$key"
	result := []byte{}
	for _, submatches := range pattern.FindAllStringSubmatchIndex(path, -1) {
		result = pattern.ExpandString(result, template, path, submatches)
	}
	return map[string]string{"key": string(result)}, matched
}

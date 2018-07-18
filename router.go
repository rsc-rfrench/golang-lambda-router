package router

import (
	"github.com/aws/aws-lambda-go/events"
	"regexp"
	"strings"
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

func createPatternFromRoute(route string) string {
	components := strings.Split(route, "/")
	for i, component := range components {
		if strings.HasPrefix(component, ":") {
			name := component[1:]
			regex := `(?P<` + name + `>\w+)`
			components[i] = regex
		}
	}
	return strings.Join(components, "/")
}

func createTemplateFromRoute(route string) string {
	components := strings.Split(route, "/")
	var template_items []string
	for _, component := range components {
		if strings.HasPrefix(component, ":") {
			name := component[1:]
			template_item := name + `=$` + name
			template_items = append(template_items, template_item)
		}
	}
	return strings.Join(template_items, "\n")
}

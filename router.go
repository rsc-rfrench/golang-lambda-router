package router

import (
	"github.com/aws/aws-lambda-go/events"
	"regexp"
	"strings"
)

type Request events.APIGatewayProxyRequest
type Response events.APIGatewayProxyResponse

type Route struct {
	Pattern string
	Handler func(Request, map[string]string) (Response, error)
}

type Router struct {
	GETs    []Route
	POSTs   []Route
	DELETEs []Route
}

func (r *Router) DelegateRequest(request Request) (Response, error) {
	var routes []Route
	switch request.HTTPMethod {
	case "GET":
		routes = r.GETs
	case "POST":
		routes = r.POSTs
	case "DELETE":
		routes = r.DELETEs
	default:
		routes = r.GETs
	}
	for _, route := range routes {
		params, ok := matchRoute(route.Pattern, request.Path)
		if ok {
			return route.Handler(request, params)
		}
	}
	return Response{StatusCode: 404}, nil
}

func (r *Router) GET(path string, f func(Request, map[string]string) (Response, error)) {
	r.GETs = append(r.GETs, Route{
		Pattern: path,
		Handler: f,
	})
}

func (r *Router) POST(path string, f func(Request, map[string]string) (Response, error)) {
	r.POSTs = append(r.POSTs, Route{
		Pattern: path,
		Handler: f,
	})
}

func (r *Router) DELETE(path string, f func(Request, map[string]string) (Response, error)) {
	r.DELETEs = append(r.DELETEs, Route{
		Pattern: path,
		Handler: f,
	})
}

func (r *Router) dumpRoutes() map[string][]Route {
	return map[string][]Route{
		"GET":  r.GETs,
		"POST": r.POSTs,
	}
}

func matchRoute(route string, path string) (map[string]string, bool) {
	pattern := regexp.MustCompile(createPatternFromRoute(route))
	matched := pattern.MatchString(path)
	results := make(map[string]string)
	if matched {
		template := createTemplateFromRoute(route)
		result := []byte{}
		for _, submatches := range pattern.FindAllStringSubmatchIndex(path, -1) {
			result = pattern.ExpandString(result, template, path, submatches)
		}
		for _, pair := range strings.Split(string(result), "\n") {
			kv := strings.Split(pair, "=")
			if len(kv) == 2 {
				results[kv[0]] = kv[1]
			}
		}
	}
	return results, matched
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

package router

import "testing"

func TestDelegateRequest(t *testing.T) {
	router := Router{}
	request := Request{}
	response, err := router.DelegateRequest(request)
	if response.StatusCode != 200 {
		t.Fail()
	}
	if err != nil {
		t.Fail()
	}
}

func hello(Request) (Response, error) {
	return Response{}, nil
}

func TestGETInstallsRoute(t *testing.T) {
	installed := false
	router := Router{}
	router.GET("/hello", hello)
	routes := router.dumpRoutes()
	if _, ok := routes["/hello"]; ok {
		installed = true
	}
	if !installed {
		t.Fail()
	}
}

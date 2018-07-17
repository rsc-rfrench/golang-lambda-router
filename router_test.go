package router

import "testing"

func hello(Request) (Response, error) {
	return Response{
		Body:       "hello",
		StatusCode: 200,
	}, nil
}

func goodbye(Request) (Response, error) {
	return Response{
		Body:       "goodbye",
		StatusCode: 200,
	}, nil
}

func TestGETInstallsRoute(t *testing.T) {
	router := Router{}
	router.GET("/hello", hello)
	routes := router.dumpRoutes()
	_, ok := routes["/hello"]
	if !ok {
		t.Fail()
	}
}

func TestGETExecutesDesignatedHandler(t *testing.T) {
	router := Router{}
	router.GET("/hello", hello)

	_, err := router.DelegateRequest(Request{
		Path: "/hello",
	})
	if err != nil {
		t.Fail()
	}
}

func TestHandlerReturnsDesiredBody(t *testing.T) {
	router := Router{}
	router.GET("/hello", hello)

	resp, _ := router.DelegateRequest(Request{
		Path: "/hello",
	})
	if resp.Body != "hello" {
		t.Fail()
	}
}

func TestGETCanDistinguishHandlers(t *testing.T) {
	router := Router{}
	router.GET("/hello", hello)
	router.GET("/goodbye", goodbye)

	resp, _ := router.DelegateRequest(Request{
		Path: "/goodbye",
	})
	if resp.Body != "goodbye" {
		t.Fail()
	}
}

func TestHandlerOrderIsCommutative(t *testing.T) {
	router := Router{}
	router.GET("/goodbye", goodbye)
	router.GET("/hello", hello)

	resp, _ := router.DelegateRequest(Request{
		Path: "/goodbye",
	})
	if resp.Body != "goodbye" {
		t.Fail()
	}
}

func TestHandlersDontShadow(t *testing.T) {
	router := Router{}
	router.GET("/goodbye", goodbye)
	router.GET("/hello", hello)

	resp, _ := router.DelegateRequest(Request{
		Path: "/hello",
	})
	if resp.Body != "hello" {
		t.Fail()
	}
}

func TestMissingPathGets404(t *testing.T) {
	router := Router{}

	resp, _ := router.DelegateRequest(Request{
		Path: "/hello",
	})
	if resp.StatusCode != 404 {
		t.Fail()
	}
}

func TestStaticRouteMatches(t *testing.T) {
	route := `/path`
	_, matches := matchRoute(route, "/path")
	if !matches {
		t.Fail()
	}
}

func TestBadStaticRouteDoesntMatch(t *testing.T) {
	route := `/path`
	_, matches := matchRoute(route, "/junk")
	if matches {
		t.Fail()
	}
}

func TestRouteParamMatches(t *testing.T) {
	route := `/path/(?P<key>\w+)`
	_, matches := matchRoute(route, "/path/value")
	if !matches {
		t.Fail()
	}
}

func TestBadRouteParamDoesntMatch(t *testing.T) {
	route := `/path/:key`
	_, matches := matchRoute(route, "/junk/value")
	if matches {
		t.Fail()
	}
}

func TestRouteParamCapturesKey(t *testing.T) {
	route := `/path/:key`
	results, _ := matchRoute(route, "/path/value")
	_, ok := results["key"]
	if !ok {
		t.Fail()
	}
}

func TestRouteParamCapturesValue(t *testing.T) {
	route := `/path/(?P<key>\w+)`
	results, _ := matchRoute(route, "/path/value")
	value, _ := results["key"]
	if value != "value" {
		t.Fail()
	}
}

func TestCreatePatternFromRoute(t *testing.T) {
	pattern := createPatternFromRoute(`/path/:key`)
	if pattern != `/path/(?P<key>\w+)` {
		t.Fail()
	}
}

func TestCreatePatternFromMultiparameterRoute(t *testing.T) {
	pattern := createPatternFromRoute(`/path/:a/fixed/:b`)
	if pattern != `/path/(?P<a>\w+)/fixed/(?P<b>\w+)` {
		t.Fail()
	}
}

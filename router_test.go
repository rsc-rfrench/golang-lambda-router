package router

import "testing"

func hello(Request, map[string]string) (Response, error) {
	return Response{
		Body:       "hello",
		StatusCode: 200,
	}, nil
}

func goodbye(Request, map[string]string) (Response, error) {
	return Response{
		Body:       "goodbye",
		StatusCode: 200,
	}, nil
}

func TestGETInstallsRoute(t *testing.T) {
	router := Router{}
	if len(router.dumpRoutes()) != 0 {
		t.Fail()
	}
	router.GET("/hello", hello)
	if len(router.dumpRoutes()) != 1 {
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
	route := `/path/:key`
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
	route := `/path/:key`
	results, _ := matchRoute(route, "/path/value")
	value, _ := results["key"]
	if value != "value" {
		t.Errorf("[%s]", value)
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

func TestCreateTemplateFromMultiparameterRoute(t *testing.T) {
	template := createTemplateFromRoute(`/path/:a/fixed/:b`)
	if template != "a=$a\nb=$b" {
		t.Error(template)
	}
}

func TestCreateTemplateFromRoute(t *testing.T) {
	template := createTemplateFromRoute(`/path/:a`)
	if template != "a=$a" {
		t.Error(template)
	}
}

func TestCreateTemplateFromFixedRoute(t *testing.T) {
	template := createTemplateFromRoute(`/path/fixed`)
	if template != "" {
		t.Error(template)
	}
}

func TestParseFixedPathHasZeroParams(t *testing.T) {
	route := `/path/fixed`
	path := `/path/fixed`
	path_params, _ := matchRoute(route, path)
	if len(path_params) != 0 {
		t.Error(len(path_params))
	}
}

func TestParseRouteYieldsCorrectNumberOfParams(t *testing.T) {
	route := `/path/:a/fixed/:b`
	path := `/path/aye/fixed/bee`
	path_params, _ := matchRoute(route, path)
	if len(path_params) != 2 {
		t.Fail()
	}
}

func TestParseRouteYieldsCorrectParamValues(t *testing.T) {
	route := `/path/:a/fixed/:b`
	path := `/path/aye/fixed/bee`
	path_params, _ := matchRoute(route, path)
	if path_params["a"] != "aye" {
		t.Fail()
	}
	if path_params["b"] != "bee" {
		t.Fail()
	}
}

func greet(_ Request, x map[string]string) (Response, error) {
	return Response{
		Body:       "Hello, " + x["name"],
		StatusCode: 200,
	}, nil
}

func TestGETGreetBruce(t *testing.T) {
	router := Router{}
	router.GET("/greet/:name", greet)

	resp, _ := router.DelegateRequest(Request{
		Path: "/greet/bruce",
	})
	if resp.Body != "Hello, bruce" {
		t.Error(resp.Body)
	}
}

func TestGETGreetLucy(t *testing.T) {
	router := Router{}
	router.GET("/greet/:name", greet)

	resp, _ := router.DelegateRequest(Request{
		Path: "/greet/lucy",
	})
	if resp.Body != "Hello, lucy" {
		t.Error(resp.Body)
	}
}

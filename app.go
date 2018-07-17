package main

func help(request Request) (Response, error) {
	return Response{
		Body:       "Hello my friend",
		StatusCode: 200,
	}, nil
}

func main() {
	r := Router{}
	r.GET("/help", help)
	r.Start()
}

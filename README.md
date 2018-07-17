# Golang Lambda Router
*Example of request routing for Golang on Lambda*

### Goals
* Write something **specifically for Lambda**
  * Small collection of routes
  * No need for caching routes since the lambda will not live long
* Applications should look like libraries, be thoroughly testable without http

### Prior Art
* [Julien Schmidt's HTTP Router][schmidt] which expects HTTP requests
* [Golang's net/http ServeMux][golang] which also expects HTTP requests

[schmidt]: https://github.com/julienschmidt/httprouter
[golang]: https://golang.org/pkg/net/http/#ServeMux

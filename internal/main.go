package main

import (
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/scraly/learning-go-by-examples/go-rest-api/pkg/swagger/server/restapi"
	"github.com/scraly/learning-go-by-examples/go-rest-api/pkg/swagger/server/restapi/operations"
)

func main() {
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
    if err != nil {
        log.Fatalln(err)
    }

	api := operations.NewAPI(swaggerSpec)
    server := restapi.NewServer(api)

    defer func() {
        if err := server.Shutdown(); err != nil {
            log.Fatalln(err)
        }
    }()

	
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Its alive, %q", html.EscapeString(r.URL.Path))
	})

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func Health(operations.CheckHealthParams) middleware.Responder {
    return operations.NewCheckHealthOK().WithPayload("OK")
}

func GetHelloUser(user operations.GetHelloUserParams) middleware.Responder {
    return operations.NewGetHelloUserOK().WithPayload("Hello " + user.User + "!")
}

func GetGopherByName(gopher operations.GetGopherNameParams) middleware.Responder {

    var URL string
    if gopher.Name != "" {
        URL = "https://raw.githubusercontent.com/vitortrimer/gophers/master/" + gopher.Name + ".png"
    } else {
        URL = "https://raw.githubusercontent.com/vitortrimer/gophers/master/BELGIUM.png"
    }

    response, err := http.Get(URL)
    if err != nil {
        fmt.Println("error")
    }

    return operations.NewGetGopherNameOK().WithPayload(response.Body)
}

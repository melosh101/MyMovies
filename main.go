package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/angelofallars/htmx-go"
	"github.com/joho/godotenv"
)

//go:embed public
var public embed.FS

func main() {

	envErr := godotenv.Load()
	if envErr != nil {
		log.Fatalln("error loading .env file")
	}
	fmt.Println("Hello World")
	fmt.Println(htmx.SwapAfterEnd.
		Scroll(htmx.Bottom).
		SettleAfter(time.Millisecond * 500),
	)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		templ := template.Must(template.ParseFS(public, "index.html"))
		err := templ.Execute(w, nil)
		if err != nil {
			log.Println(err)
			return
		}
	})

	fmt.Println("server running...")
	httpErr := http.ListenAndServe(":8080", nil)
	if httpErr != nil {
		fmt.Println("failed to listen on port 8080")
		return
	}
}

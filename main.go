package main

import (
	"log"
	"net/http"

	c "github.com/ianthomasict/go-templ-pdf/components"
)




func main() {
    srv := http.NewServeMux()


    invoices := []c.Invoice{
        {Name: "Roy's Solutions", Total: 1022.53},
        {Name: "Ribbiq IT", Total: 627.03},
        {Name: "Value-brand Uniques", Total: 5562.91},
    }

    srv.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
        c.Page(invoices).Render(r.Context(), w);
    })

    log.Println("Starting server on port 8080!")
    log.Fatalln(http.ListenAndServe(":8080", srv))
}

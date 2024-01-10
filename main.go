package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-rod/rod"
	"github.com/ianthomasict/go-templ-pdf/src/reports"
	"github.com/ianthomasict/go-templ-pdf/src/view"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	srv := echo.New();

	invoices := []view.Invoice{
		{Name: "Roy's Solutions", Total: 1022.53},
		{Name: "Ribbiq IT", Total: 627.03},
		{Name: "Value-brand Uniques", Total: 5562.91},
	}

	srv.Use(middleware.CORS())
	srv.Use(middleware.Logger())

	srv.GET("/", func(ctx echo.Context) error {
		return view.Page(invoices).Render(ctx.Request().Context(), ctx.Response().Writer)
	}) 

	api := srv.Group("/api")
	api.GET("/download/:id", func(c echo.Context) error {
		id, err := strconv.ParseInt(c.Param("id"), 10, 0); 
		if err != nil {
			c.Logger().Error(err)
			return c.JSON(http.StatusBadRequest, "Missing valid 'id' param")
		}
		name := c.QueryParam("name")
		if id == 0 || name == "" {
			c.Logger().Error(err)
			return c.JSON(400, "missing required parameter");
		}

		// html, err := templ.ToGoHTML(c.Request().Context(), reports.MainReport(name, id))
		// if err != nil {
		// 	return c.JSON(500, "An error occurred while generating report contents")
		// }

		htmlBuf := new(bytes.Buffer)
		reports.MainReport(name,id).Render(c.Request().Context(), htmlBuf)
		fmt.Println(htmlBuf.String()[:100])

		pdfBytes, err := generatePDF(c, htmlBuf.String())
		if err != nil {
			c.Logger().Error(err)
			return c.JSON(500, "An error occurred while generating PDF")
		}


		return c.Blob(200, "application/pdf", pdfBytes)
	})

	log.Println("Starting server on port 8080!")
	log.Fatalln(http.ListenAndServe(":8080", srv))
}



func generatePDF(c echo.Context, html string) ([]byte, error) {
	c.Logger().Info("Initializing page")
	page := rod.New().MustConnect().MustPage("http://google.com")
	c.Logger().Info("Setting content")
	err := page.SetDocumentContent(html)
	if err != nil {
		return nil, fmt.Errorf("Failed to set page html content: %w", err)
	}
	pdfBytes := page.MustPDF("") ;

	return pdfBytes, nil;
}


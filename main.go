package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/go-rod/rod"
	"github.com/ianthomasict/go-templ-pdf/src/reports"
	"github.com/ianthomasict/go-templ-pdf/src/view"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type TemplateLoader struct {
	ctx context.Context
	buf *bytes.Buffer
}

func (t *TemplateLoader) load(tmpl templ.Component) {
	tmpl.Render(t.ctx, t.buf)
}

func (t *TemplateLoader) toPDF() ([]byte, error) {
	html := t.buf.String()

	page := rod.New().MustConnect().MustPage()
	err := page.SetDocumentContent(html)
	if err != nil {
		return nil, fmt.Errorf("Failed to set page html content: %w", err)
	}
	pdfBytes := page.MustPDF()

	return pdfBytes, nil
}

func main() {
	srv := echo.New()

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
		id, err := strconv.ParseInt(c.Param("id"), 10, 0)
		if err != nil || id == 0 {
			c.Logger().Error(err)
			return c.JSON(http.StatusBadRequest, fmt.Sprintf("Missing valid 'id' param: %v", id))
		}
		name := c.QueryParam("name")
		if name == "" {
			c.Logger().Error("Missing required parameter: 'name'")
			return c.JSON(400, "Missing required parameter: 'name'")
		}

		templater := TemplateLoader{c.Request().Context(), new(bytes.Buffer)}
		templater.load(reports.MainReport(name, id))

		pdfBytes, err := templater.toPDF()
		if err != nil {
			c.Logger().Error(err)
			return c.JSON(500, "An error occurred while generating PDF")
		}

		return c.Blob(200, "application/pdf", pdfBytes)
	})

	log.Println("Starting server on port 8080!")
	log.Fatalln(http.ListenAndServe(":8080", srv))
}

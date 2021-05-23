package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

var uri string

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	uri = "https://raw.githubusercontent.com/tacklehop/csvsearch/main/sample.csv"
	t := &Template{
		templates: template.Must(template.ParseGlob("*.html")),
	}
	e := echo.New()
	e.Renderer = t
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to CSV Search Cloud!")
	})
	e.GET("/search", searchHandler)
	e.POST("/save", saveHandler)
	e.GET("/result", resultHandler)

	e.Logger.Fatal(e.Start(":8080"))
}

func searchHandler(c echo.Context) error {
	fmt.Println("In : searchHandler()")
	defer fmt.Println("Out: searchHandler()")
	return c.Render(http.StatusOK, "search", "Enter")
}

func saveHandler(c echo.Context) error {
	fmt.Println("In : saveHandler()")
	defer fmt.Println("Out: saveHandler()")
	key := c.FormValue("body")
	result, err := searchCsvFromHttp(uri, key)
	if err != nil {
		fmt.Println("CSV search error:", err)
		return err
	}
	err = ioutil.WriteFile("result.txt", []byte(result), 0600)
	if err != nil {
		fmt.Println("File Not Written:", err)
		return err
	}
	return c.Redirect(http.StatusFound, "/result")
}

func resultHandler(c echo.Context) error {
	fmt.Println("In : resultHandler()")
	defer fmt.Println("Out: resultHandler()")
	body, err := ioutil.ReadFile("result.txt")
	if err != nil {
		fmt.Println("File Not Read:", err)
		return err
	}
	data := struct {
		Body string
	}{
		Body: string(body),
	}
	return c.Render(http.StatusOK, "result", data)
}

func searchCsvFromHttp(uri, searchWord string) (string, error) {
	// Handle input file
	response, err := http.Get(uri)
	if err != nil {
		fmt.Println("http.Get() error:", err)
		return "", err
	}
	defer response.Body.Close()

	result := ""

	// Search CSV and output hit lines to stdout
	reader := csv.NewReader(response.Body)
	reader.FieldsPerRecord = -1 // Accepts reading irregular csv format
	for {
		words, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("File read error:", err)
			return "", err
		}
		for _, w := range words {
			// In case dir/file in URI does not exist
			if w == "404: Not Found" {
				fmt.Println(w)
				return "", errors.New(w)
			} else if strings.Contains(w, searchWord) {
				fmt.Println(words)
				result += strings.Join(words, ", ")
				result += "\n"
				break
			}
		}
	}
	return result, nil
}

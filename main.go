package main

import (
	"fmt"
	"strings"

	"github.com/LEEBONGHAK/job_scrapper/scrapper"
	"github.com/labstack/echo"
)

func handleHome(c echo.Context) error {
	return c.File("home.html")
}

func handleSrape(c echo.Context) error {
	term := strings.ToLower(scrapper.CleanString(c.FormValue("term")))
	fmt.Println(term)
	return nil
}

func main() {
	e := echo.New()
	e.GET("/", handleHome)
	e.POST("/scrape", handleSrape)
	e.Logger.Fatal(e.Start(":1323"))
}

package main

import (
	"os"
	"strings"

	"github.com/LEEBONGHAK/job_scrapper/scrapper"
	"github.com/labstack/echo"
)

const FILE_NAME string = "jobs.csv"

func handleHome(c echo.Context) error {
	return c.File("home.html")
}

func handleSrape(c echo.Context) error {
	defer os.Remove(FILE_NAME)
	term := strings.ToLower(scrapper.CleanString(c.FormValue("term")))
	scrapper.Scrape(term)
	// 첨부파일을 리턴하는 기능
	return c.Attachment(FILE_NAME, FILE_NAME)
}

func main() {
	e := echo.New()
	e.GET("/", handleHome)
	e.POST("/scrape", handleSrape)
	e.Logger.Fatal(e.Start(":1323"))
}

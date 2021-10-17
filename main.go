package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type extractedJob struct {
	id       string
	title    string
	location string
	salary   string
	summary  string
}

var baseURL string = "https://kr.indeed.com/jobs?q=python"

func main() {
	totalPages := getPages()
	for i := 0; i < totalPages; i++ {
		getPage(i)
	}
}

// 각 page에서 정보를 가저오는 함수
func getPage(page int) {
	pageURL := baseURL + "&start=" + strconv.Itoa(page*10)
	fmt.Println(pageURL)
	response, err := http.Get(pageURL)
	checkErr(err)
	checkCode(response)

	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	checkErr(err)

	searchCards := document.Find(".tapItem ")
	searchCards.Each(func(i int, card *goquery.Selection) {
		id, _ := card.Attr("data-jk")
		fmt.Println(id)
		title := card.Find(".jobTitle").Text()
		fmt.Println(title)
		location := card.Find(".companyLocation").Text()
		fmt.Println(location)
	})
}

func cleanString(str string) string {

}

// page의 개수를 가져오는 함수
func getPages() int {
	pages := 0
	response, err := http.Get(baseURL)
	checkErr(err)
	checkCode(response)

	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	checkErr(err)

	document.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		pages = s.Find("a").Length()
	})

	return pages
}

// error를 check하는 함수 (error가 발생하면 프로그램 종료)
func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

// 제대로된 Response를 가져오는지 check하는 함수 (200이 아니면 프로그램 종료)
func checkCode(response *http.Response) {
	if response.StatusCode != 200 {
		log.Fatalln("Request failed with Statues:", response.StatusCode)
	}
}

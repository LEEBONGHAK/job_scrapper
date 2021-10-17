package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

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
	var jobs []extractedJob
	totalPages := getPages()

	for i := 0; i < totalPages; i++ {
		extractedJobs := getPage(i)
		// 각각의 array을 하나의 array로 만드는 방법 -> not [[x1], [x2], [x3]] but [x1, x2, x3]
		jobs = append(jobs, extractedJobs...)
	}

	writeJobs(jobs)
	fmt.Println("Done, extracted", len(jobs))
}

// csv 파일에 작성하는 함수
func writeJobs(jobs []extractedJob) {
	// csv file 만들기
	file, err := os.Create("jobs.csv")
	checkErr(err)

	write := csv.NewWriter(file)
	// .Flush(): 함수가 끝나는 시점에 파일에 데이터를 입력하는 함수
	defer write.Flush()

	headers := []string{"ID", "Title", "Location", "Salary", "Summary"}
	writeErr := write.Write(headers)
	checkErr(writeErr)

	for _, job := range jobs {
		jobSlice := []string{"https://kr.indeed.com/viewjob?jk=" + job.id, job.title, job.location, job.salary, job.summary}
		writeErr = write.Write(jobSlice)
		checkErr(writeErr)
	}
}

// 각 page에서 각 card의 정보를 추출해서 array로 return
func getPage(page int) []extractedJob {
	var jobs []extractedJob
	pageURL := baseURL + "&start=" + strconv.Itoa(page*10)
	response, err := http.Get(pageURL)
	checkErr(err)
	checkCode(response)

	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	checkErr(err)

	searchCards := document.Find(".tapItem")
	searchCards.Each(func(i int, card *goquery.Selection) {
		job := extractJob(card)
		jobs = append(jobs, job)
	})
	return jobs
}

// card 내 정보를 추출하는 함수
func extractJob(card *goquery.Selection) extractedJob {
	id, _ := card.Attr("data-jk")
	title := cleanString(card.Find(".jobTitle").Text())
	location := cleanString(card.Find(".companyLocation").Text())
	salary := cleanString(card.Find(".salary-snippet").Text())
	summary := cleanString(card.Find(".job-snippet").Text())

	return extractedJob{
		id:       id,
		title:    title,
		location: location,
		salary:   salary,
		summary:  summary,
	}
}

// string의 앞뒤 공백을 없애고, 모든 공백을 없앤 후 배열로 만들고 strings.Join()으로 재구성
func cleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
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

package sporclient

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strconv"
	"strings"
	"time"
)

const URL = "https://online.spor.istanbul/uyeseanssecim"

func Run() string {
	timeMap := map[time.Weekday]string{
		time.Sunday:    "6",
		time.Saturday:  "5",
		time.Friday:    "4",
		time.Thursday:  "3",
		time.Wednesday: "2",
		time.Tuesday:   "1",
		time.Monday:    "0",
	}
	cookieJar, _ := cookiejar.New(nil)

	client := &http.Client{
		Jar: cookieJar,
	}
	req, err := http.NewRequest(http.MethodGet, URL, nil)

	if err != nil {
		log.Fatalf("Error while getting the RESPONSE: %s\n", err)
	}

	req.AddCookie(&http.Cookie{Name: "ASP.NET_SessionId", Value: ""})
	req.AddCookie(&http.Cookie{Name: "_hjSessionUser_1859734", Value: ""})

	resp, err := client.Do(req)

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error while getting the RESPONSE: status code is different\n", err)
	}

	bodyBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatalf("Couldn't read the body")
	}

	//body := string(bodyBytes)

	newfile, _ := os.Create("newFile.html")

	defer newfile.Close()

	newfile.Write(bodyBytes)

	if err != nil {
		log.Fatal("error while parsing the doc\n")
	}
	doc, _ := goquery.NewDocumentFromReader(resp.Body)

	initialText := ""
	doc.Find("#pageContent_dvScheduler .panel").Each(func(i int, selection *goquery.Selection) {
		item := selection.Find(".panel-title")
		text := strings.TrimSpace(item.Text())
		date := strings.Fields(text)[1]
		splitted := strings.Split(date, ".")
		normalDate := splitted[2] + "-" + splitted[1] + "-" + splitted[0]
		parsedDate, _ := time.Parse("2006-01-02", normalDate)

		// For each .form-group inside the .panel-body
		selection.Find(".panel-body .form-group").Each(func(index int, s *goquery.Selection) {
			sessionFinder := "#pageContent_rptList_ChildRepeater_" + timeMap[parsedDate.Weekday()] + "_lblSeansSaat_" + strconv.Itoa(index)
			// Extracting and printing details for demonstration
			sessionTime := s.Find(sessionFinder).Text()
			sessionCapacityLabel := strings.TrimSpace(s.Find(".well label.label-default").Text())
			sessionGender := s.Find(".well label.cinsiyet").Text()
			activeMembers, _ := strconv.Atoi(s.Find(".well span.label-danger").Text())
			quota, _ := strconv.Atoi(s.Find(".well span.label-primary").Text())

			initialText += fmt.Sprintf("Session Time: %s\n", sessionTime)
			initialText += fmt.Sprintf("Capacity Status: %s\n", sessionCapacityLabel)
			initialText += fmt.Sprintf("Gender: %s\n", sessionGender)
			initialText += fmt.Sprintf("Active Members: %v\n", activeMembers)
			initialText += fmt.Sprintf("Quota: %v\n", quota)
			initialText += fmt.Sprintf("---")
		})
	})
	return initialText
}

package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type Article struct {
	Publisher   string
	TimeAgo     string
	TimeMinutes int
	Title       string
	URL         string
	CommentsURL string
}

func timeAgoToMinutes(timeAgo string) int {
	timeAgo = strings.TrimSpace(timeAgo)

	if strings.HasSuffix(timeAgo, "mo") {
		months, _ := strconv.Atoi(strings.TrimSuffix(timeAgo, "mo"))
		return months * 30 * 24 * 60
	} else if strings.HasSuffix(timeAgo, "w") {
		weeks, _ := strconv.Atoi(strings.TrimSuffix(timeAgo, "w"))
		return weeks * 7 * 24 * 60
	} else if strings.HasSuffix(timeAgo, "d") {
		days, _ := strconv.Atoi(strings.TrimSuffix(timeAgo, "d"))
		return days * 24 * 60
	} else if strings.HasSuffix(timeAgo, "h") {
		hours, _ := strconv.Atoi(strings.TrimSuffix(timeAgo, "h"))
		return hours * 60
	} else if strings.HasSuffix(timeAgo, "m") {
		minutes, _ := strconv.Atoi(strings.TrimSuffix(timeAgo, "m"))
		return minutes
	} else if strings.HasSuffix(timeAgo, "y") {
		years, _ := strconv.Atoi(strings.TrimSuffix(timeAgo, "y"))
		return years * 365 * 24 * 60 // Approximate a year as 365 days
	}

	return 0 // fallback if format is unknown
}

func main() {
	URL := "https://devurls.com/"
	resp, err := http.Get(URL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Load the HTML document directly from the response body
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var articles []Article

	// Find each publisher block
	doc.Find(".publisher-block").Each(func(i int, s *goquery.Selection) {
		publisherName := s.Find(".publisher-text .primary").Text()

		// Find all article links within this publisher block
		s.Find(".publisher-link").Each(func(j int, linkS *goquery.Selection) {
			timeAgo := linkS.Find(".aside .text").Text()
			title := linkS.Find(".article-link").Text()
			url, _ := linkS.Find(".article-link").Attr("href")
			commentsURL, _ := linkS.Attr("data-comment-url")

			article := Article{
				Publisher:   publisherName,
				TimeAgo:     strings.TrimSpace(timeAgo),
				TimeMinutes: timeAgoToMinutes(timeAgo),
				Title:       strings.TrimSpace(title),
				URL:         strings.TrimSpace(url),
				CommentsURL: strings.TrimSpace(commentsURL),
			}
			if article.TimeMinutes < 10080 {
				articles = append(articles, article)
			}
		})
	})

	// Print extracted articles
	sort.Slice(articles, func(i, j int) bool {
		return timeAgoToMinutes(articles[i].TimeAgo) < timeAgoToMinutes(articles[j].TimeAgo)
	})
	for _, article := range articles {
		fmt.Println(article)
	}
	fmt.Println(len(articles))

}

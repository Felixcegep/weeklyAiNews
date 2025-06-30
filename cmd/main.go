package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	read "github.com/go-shiori/go-readability"
	"github.com/gocolly/colly/v2"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Article struct {
	Publisher   string
	TimeAgo     string
	TimeMinutes int
	Title       string
	URL         string
	ParsedURL   string
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
func extractLinks() []Article {
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
			link, _ := linkS.Find(".article-link").Attr("href")
			commentsURL, _ := linkS.Attr("data-comment-url")
			parsedURL, err := url.Parse(link)
			if err != nil {
				fmt.Println(err)
			}
			baseURL := fmt.Sprintf("%s://%s/", parsedURL.Scheme, parsedURL.Host)
			article := Article{
				Publisher:   publisherName,
				TimeAgo:     strings.TrimSpace(timeAgo),
				TimeMinutes: timeAgoToMinutes(timeAgo),
				Title:       strings.TrimSpace(title),
				URL:         strings.TrimSpace(link),
				ParsedURL:   strings.TrimSpace(baseURL),
				CommentsURL: strings.TrimSpace(commentsURL),
			}
			if article.TimeMinutes < 1440 {
				articles = append(articles, article)
			}
		})
	})
	return articles
}
func clean(s string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(s)), " ")
}
func Extract(url string) (title, body string, err error) {
	// 1 — create the collector with ordinary options
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) " +
			"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125 Safari/537.36"),
	)

	// 2 — attach your custom transport *after* creation
	c.WithTransport(&http.Transport{IdleConnTimeout: 20 * time.Second}) // ← fixed

	// 3 — polite limits
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		Delay:       2 * time.Second,
		RandomDelay: 800 * time.Millisecond,
	})

	// ---------- Readability pass ----------
	var readableOK bool
	c.OnResponse(func(r *colly.Response) {
		art, e := read.FromReader(strings.NewReader(string(r.Body)), r.Request.URL)
		if e == nil && strings.TrimSpace(art.TextContent) != "" {
			readableOK = true
			title = art.Title
			body = clean(art.TextContent)
		}
	})

	// ---------- CSS fallback --------------
	var sb strings.Builder
	for _, sel := range []string{
		"article", "main", `div[class*='content']`, `section[itemprop='articleBody']`,
	} {
		c.OnHTML(sel, func(e *colly.HTMLElement) {
			if !readableOK {
				sb.WriteString(e.Text + " ")
			}
		})
	}
	c.OnScraped(func(_ *colly.Response) {
		if !readableOK {
			body = clean(sb.String())
		}
	})

	// ---------- one retry on 5xx ----------
	c.OnError(func(r *colly.Response, e error) {
		if r.StatusCode >= 500 && r.Request.Ctx.Get("retried") == "" {
			r.Request.Ctx.Put("retried", "true")
			_ = r.Request.Retry()
			return
		}
		err = e
	})

	// run
	err = c.Visit(url)
	c.Wait()
	return
}

func main() {
	articles := extractLinks()

	sort.Slice(articles, func(i, j int) bool {
		return timeAgoToMinutes(articles[i].TimeAgo) < timeAgoToMinutes(articles[j].TimeAgo)
	})
	valideUrl := []string{}
	notValideUrl := []string{}
	var sb strings.Builder

	for _, value := range articles {
		fmt.Println(len(valideUrl))
		fmt.Println(len(notValideUrl))
		_, body, err := Extract(value.URL)
		charcount := len([]rune(sb.String()))
		fmt.Println("actual len ", charcount)
		if err != nil {
			fmt.Println(err)
		}
		if len(body) > 50 {
			valideUrl = append(valideUrl, value.ParsedURL)
			title := fmt.Sprintf("%s - %s", value.Title)
			theUrl := fmt.Sprintf("%s", value.URL)
			theContent := fmt.Sprintf("%s", body)
			sb.WriteString(title + theUrl + theContent + "\n")
		} else {
			notValideUrl = append(notValideUrl, value.ParsedURL)
		}
	}

	fmt.Println("valide", valideUrl)
	fmt.Println("not valide", notValideUrl)
	fmt.Println("valide number", len(valideUrl))
	fmt.Println("not valide number", len(notValideUrl))
	err := os.WriteFile("output.txt", []byte(sb.String()), 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("job finished")
}

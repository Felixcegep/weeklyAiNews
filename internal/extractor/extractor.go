package extractor

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	read "github.com/go-shiori/go-readability"
	"github.com/gocolly/colly/v2"
	"log"
	"net/http"
	"net/url"
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
func ExtractLinks() []Article {
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
				TimeMinutes: TimeAgoToMinutes(timeAgo),
				Title:       strings.TrimSpace(title),
				URL:         strings.TrimSpace(link),
				ParsedURL:   strings.TrimSpace(baseURL),
				CommentsURL: strings.TrimSpace(commentsURL),
			}
			if article.TimeMinutes < 360 {
				articles = append(articles, article)
			}
		})
	})
	return articles
}

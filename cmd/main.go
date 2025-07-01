package main

import (
	extract "awesomeProject/internal/extractor"
	"fmt"
	"os"
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

func main() {

	articles := extract.ExtractLinks()

	valideUrl := []string{}
	notValideUrl := []string{}
	var sb strings.Builder
	var allArticles strings.Builder
	articleLlm := make(map[string]Article)

	for i, value := range articles {
		fmt.Println(len(valideUrl))
		fmt.Println(len(notValideUrl))
		_, body, err := extract.Extract(value.URL)
		charcount := len([]rune(sb.String()))
		fmt.Println("actual len ", charcount)
		if err != nil {
			fmt.Println(err)
		}
		if len(body) > 50 && len(body) < 50000 {
			valideUrl = append(valideUrl, value.URL)
			title := fmt.Sprintf("%s - %s", value.Title, value.Publisher)
			theUrl := fmt.Sprintf("%s", value.URL)
			theContent := fmt.Sprintf("%s", body)
			allArticles.WriteString(fmt.Sprintf("%d - %s %s \n", i, title, theUrl))
			sb.WriteString(title + theUrl + theContent + "\n")
			articleLlm[value.URL] = Article{
				Publisher:   value.Publisher,
				TimeAgo:     value.TimeAgo,
				TimeMinutes: value.TimeMinutes,
				Title:       value.Title,
				URL:         value.URL,
				ParsedURL:   value.ParsedURL,
				CommentsURL: value.CommentsURL,
			}
		} else {
			notValideUrl = append(notValideUrl, value.URL)
		}
	}

	fmt.Println("valide", valideUrl)
	fmt.Println("not valide", notValideUrl)
	fmt.Println("valide number", len(valideUrl))
	fmt.Println("not valide number", len(notValideUrl))
	currentTime := time.Now()
	dateToday := fmt.Sprintf("new-%d-%d-%d.txt", currentTime.Day(), currentTime.Month(), currentTime.Year())
	err := os.WriteFile("output/"+dateToday, []byte(sb.String()), 0644)

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(allArticles.String())
	fmt.Println("job finished")
	fmt.Println("------------------", articleLlm["https://github.com/codeddarkness/taco_pardons"])
}

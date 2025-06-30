package main

import (
	extract "awesomeProject/internal/extractor"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

func main() {

	articles := extract.ExtractLinks()
	sort.Slice(articles, func(i, j int) bool {
		return extract.TimeAgoToMinutes(articles[i].TimeAgo) < extract.TimeAgoToMinutes(articles[j].TimeAgo)
	})

	valideUrl := []string{}
	notValideUrl := []string{}
	var sb strings.Builder

	for _, value := range articles {
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
			sb.WriteString(title + theUrl + theContent + "\n")
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
	fmt.Println("job finished")
}

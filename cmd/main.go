package main

import (
	extract "awesomeProject/internal/extractor"
	"os"
	"strings"

	//extract "awesomeProject/internal/extractor"
	llm "awesomeProject/internal/llm"
	"fmt"
	//"fmt"
	//"log"
	//"os"
	//"strings"
	//"time"
)

func main() {
	articles := extract.ExtractLinks()
	for i, article := range articles {
		fmt.Println(i, article)
	}
	var prompt strings.Builder

	for i, article := range articles {
		format := fmt.Sprintf("%d %s %s \n", i, article.Title, article.URL)
		prompt.WriteString(format)

	}
	fmt.Println(prompt.String())

	result, err := llm.Llm_choices(prompt.String())
	if err != nil {
		fmt.Println(err)
	}
	// les liens que l'ia pense important
	articleValide := result["articles"].([]interface{})

	var allcontent strings.Builder
	for i, singleLink := range articleValide {
		urlStr, _ := singleLink.(string)
		_, body, err := extract.Extract(urlStr)
		if err != nil {
			fmt.Println(err)
		}
		allcontent.WriteString(body)
		if len(body) > 50 && len(body) < 50000 {
			fmt.Println(i, urlStr)
			linkformate := fmt.Sprintf("%d, %s, %s \n", i, urlStr, body)
			allcontent.WriteString(linkformate)
		} else {
			fmt.Printf("%d %v is not valid", i, articles)
		}
	}
	os.WriteFile("rawtext.txt", []byte(allcontent.String()), 0644)
	content := llm.LlmSummarization(allcontent.String())
	fmt.Println(content)
	os.WriteFile("formatedtext.md", []byte(content), 0644)
}

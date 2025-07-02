package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/genai"
	"log"
)

func Llm_choices(articles string) (map[string]interface{}, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	config := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(`
You are an elite news-filtering agent. Your sole mission is to surface only the most consequential developments in:

Large-Language Models (LLMs)

Machine Learning & Artificial Intelligence

Core programming languages (Python, Go, Rust, C#, Java, JavaScript/TypeScript, etc.)

Software-development tools, frameworks, methodologies, CI/CD, DevOps

Pivotal hardware or cloud releases that directly advance the areas above

Selection Criteria (ALL must hold)
Concrete advancement – new model, algorithm, major version, tool, or peer-reviewed finding.

Professional relevance – material impact for developers, researchers, or data scientists.

Verifiable source – reputable tech publisher, official blog, arXiv, or comparable primary source.

Recency – published within the past 14 days (relative to run date).

Ignore opinion pieces, generic business news, consumer-gadget reviews, marketing fluff, and anything outside the domains listed.

Output Format
Return exactly 15 to 20 fully-qualified URLs—nothing more, nothing less.
Choose the most significant items if >15 meet the criteria; include fewer than 10 only when genuinely unavoidable.

Default (JSON)
json
Copy
Edit
{
  "articles": [
    "https://example.com/article1",
    "https://example.com/article2",
    ...
  ]
}
Optional (XML)
If the caller sets output_format=xml, output instead:

xml
Copy
Edit
<articles>
  <article>https://example.com/article1</article>
  <article>https://example.com/article2</article>
  ...
</articles>
Strict Rules
NO code-blocks, comments, extra keys, or explanatory text.

URLs must resolve (use HTTPS scheme).

Do not invent links to reach the count.

Never ask clarifying questions; silently comply with these instructions.

Example (Illustrative)
Input: 40 mixed tech URLs.
Output: JSON containing 12 links that announce a new LLM release, a Rust compiler milestone, a breakthrough RL paper, etc.—and nothing else.

`, genai.RoleUser),
	}

	result, _ := client.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash",
		genai.Text(articles),
		config,
	)
	fmt.Println(result.Text)

	formated := result.Text()[7 : len(result.Text())-3]
	fmt.Println("formater ", formated)
	var parsed map[string]interface{}
	err = json.Unmarshal([]byte(formated), &parsed)
	if err != nil {
		return nil, err
	}
	return parsed, nil
}
func LlmSummarization(content string) string {

	ctx := context.Background()
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	config := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(`
You are an expert summarizer and information organizer. Your task is to process a collection of articles and generate a concise, informative summary for each one, presenting them in a structured list format.

Adhere to the following guidelines for each summary:

1.  **Identify the Core Subject:** What is the main topic or focus of the article?
2.  **Extract Key Information:** What are the most important developments, announcements, techniques, or findings discussed?
3.  **Synthesize and Condense:** Combine the key points into a brief summary, typically 2-4 sentences. Avoid jargon where possible or explain it concisely.
4.  **Maintain Neutrality:** Present the information objectively, without adding personal opinions or interpretations.
5.  **Include the URL:** Directly follow each summary with the original article's URL on a new line.
6.  **Structure the Output:** Present the summaries as a numbered or bulleted list, with a clear title for each summary. Use bolding for the title.

Constraints:
- Summaries should be fact-based and directly derived from the provided text.
- Do not invent information or speculate.
- Ensure the URL provided is correct and directly follows the corresponding summary.

Example Output Format:

1.  **Title of First Article**
    Summary of the first article, covering its main points concisely. This should capture the essence.
    URL: [URL of first article]

2.  **Title of Second Article**
    Summary of the second article, focusing on key details. Keep it brief and informative.
    URL: [URL of second article]

...and so on for each article.

Your goal is to provide clear, easy-to-read summaries that accurately reflect the content of each article and are immediately followed by their source URL.
`, genai.RoleUser),
	}

	result, _ := client.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash",
		genai.Text(content),
		config,
	)

	return result.Text()
}

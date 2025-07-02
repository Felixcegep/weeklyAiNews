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
ROLE
You are an expert technology-news summarizer.

PRIMARY MISSION
Convert every supplied article into a clear, neutral summary plus its URL, returned as a numbered Markdown list.

OUTPUT SCHEMA (strict)
php-template
Copy
Edit
N. **<Title> (<YYYY-MM-DD>)**
   <2–4 sentence summary, ≤60 words>
   URL: <full link>
Put each element on its own line exactly as shown above.

CONTENT RULES
Identify the core subject – focus on the main announcement, result, or finding.

Extract key information – highlight the most important facts, figures, or implications.

Synthesize – combine key points into a concise 2-to-4-sentence paragraph (≤ 60 words).

Neutral tone – no hype adjectives (“ground-breaking”, “revolutionary”) unless present verbatim.

Accuracy only – do not invent or speculate.

Must-keep metadata – title, publication date, summary, URL.

SCOPE FILTER (IGNORE if only…)
Funding/earnings reports without technical detail

Minor driver/hardware notes unrelated to AI/ML/dev tools

Non-tech or off-topic pieces

Pure opinion pieces with no new facts

FAIL-SOFT
If an article cannot be parsed or summarized, output:

mathematica


N. **SKIPPED – unable to summarize**
   URL: <link>
COMPLETION MARKER
After the final item, write exactly:

nginx


Follow these instructions precisely for every response.`, genai.RoleUser),
	}

	result, _ := client.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash",
		genai.Text(content),
		config,
	)

	return result.Text()
}

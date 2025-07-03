package extractor

import "testing"

func TestExtract(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{"Reuters – Microsoft to cut about 4% of jobs amid hefty AI bets", "https://www.reuters.com/business/world-at-work/microsoft-lay-off-many-9000-employees-seattle-times-reports-2025-07-02/"},
		{"BBC – Microsoft to lay off up to 9,000 employees", "https://www.bbc.com/news/articles/cdxl0w1w394o"},
		{"The Daily Beast – Microsoft to Slash 9,000 Workers Amid AI Investment", "https://www.thedailybeast.com/microsoft-to-slash-9000-workers-amid-ai-investment"},
		{"The Guardian – Microsoft to lay off 6,000 workers despite streak of profitable quarters", "https://www.theguardian.com/technology/2025/may/13/microsoft-layoffs"},
		{"KUOW – Latest Microsoft layoffs could hit 9,000 employees", "https://www.kuow.org/stories/latest-microsoft-layoffs-could-hit-9-000-employees"},
		{"Yahoo Finance – Microsoft said it would cut 9,000 jobs as tech giant bets on AI", "https://finance.yahoo.com/news/microsoft-lays-off-9-000-170029460.html"},
		{"CNBC via San Antonio News4 – Microsoft to cut 9,000 jobs in largest layoff since 2023", "https://news4sanantonio.com/news/instagram/microsoft-to-cut-9000-jobs-in-largest-layoff-since-2023-uncertainty-looms-over-teams-employees-purge-tech-companies-workforce"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, body, err := Extract(tt.url)
			if err != nil {
				t.Errorf("Extract() error = %v", err)
			}
			if len(body) == 0 {
				t.Errorf("Extract() nothing  empty")
			}
		})
	}
}
func TestExtractLinks(t *testing.T) {
	articles := ExtractLinks()
	if len(articles) == 0 {
		t.Errorf("ExtractLinks() no articles where extracted")
	}
	for _, article := range articles {
		if article.TimeMinutes > 10080 {
			t.Errorf("ExtractLinks() article = %d, want <= 10080", article.TimeMinutes)
		}

	}
}

package mcp

type BuiltinMCP struct{}

func NewBuiltinMCP() *BuiltinMCP {
	return &BuiltinMCP{}
}

type ExaSearchResult struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Snippet string `json:"snippet"`
}

func (m *BuiltinMCP) ExaSearch(query string, numResults int) ([]ExaSearchResult, error) {
	return []ExaSearchResult{
		{
			Title:   "Example Result",
			URL:     "https://example.com",
			Snippet: "This is an example search result",
		},
	}, nil
}

func (m *BuiltinMCP) Context7Docs(query string) (string, error) {
	return "Context7 documentation for: " + query, nil
}

func (m *BuiltinMCP) GrepApp(query string) ([]string, error) {
	return []string{
		"https://github.com/example/repo1",
		"https://github.com/example/repo2",
	}, nil
}

package tools_test

import (
	"context"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"searxng-mcp/internal/mcp/tools"
	"searxng-mcp/pkg/searxng"
)

var _ = Describe("MCP Tools", func() {
	var (
		server *httptest.Server
		client searxng.Client
		ctx    context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		
		// Create mock SearXNG server
		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"query": "machine learning tutorial",
				"number_of_results": 2500,
				"results": [
					{
						"url": "https://example.com/ml-tutorial",
						"title": "Complete Machine Learning Tutorial",
						"content": "Learn machine learning from scratch with this comprehensive tutorial covering algorithms, data preprocessing, and model evaluation.",
						"thumbnail": "",
						"engine": "google",
						"template": "default.html",
						"parsed_url": ["https", "example.com", "/ml-tutorial", "", "", ""],
						"img_src": "",
						"priority": "",
						"engines": ["google"],
						"positions": [1],
						"score": 18.5,
						"category": "science",
						"publishedDate": "2024-01-15"
					},
					{
						"url": "https://tech.example.com/ai-guide",
						"title": "AI and Machine Learning Development Guide",
						"content": "A practical guide for developers looking to implement machine learning solutions in their applications.",
						"thumbnail": "",
						"engine": "bing",
						"template": "default.html", 
						"parsed_url": ["https", "tech.example.com", "/ai-guide", "", "", ""],
						"img_src": "",
						"priority": "",
						"engines": ["bing"],
						"positions": [2],
						"score": 16.2,
						"category": "it",
						"publishedDate": null
					}
				]
			}`))
		}))
		
		client = searxng.NewClient(server.URL)
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("Category Search", func() {
		Context("with valid categories science and it", func() {
			It("should validate categories correctly", func() {
				// Test category validation directly
				scienceCategory, err := searxng.ValidateCategory("science")
				Expect(err).NotTo(HaveOccurred())
				Expect(scienceCategory).To(Equal(searxng.CategoryScience))
				
				itCategory, err := searxng.ValidateCategory("it")
				Expect(err).NotTo(HaveOccurred())
				Expect(itCategory).To(Equal(searxng.CategoryIT))
			})

			It("should perform search with multiple categories", func() {
				// Test the SearchWithCategory function
				result, err := searxng.SearchWithCategory(ctx, client, "machine learning tutorial", 
					searxng.CategoryScience, searxng.CategoryIT)

				Expect(err).NotTo(HaveOccurred())
				Expect(result).NotTo(BeNil())
				Expect(result.Query).To(Equal("machine learning tutorial"))
				Expect(result.NumberOfResults).To(Equal(2500))
				Expect(result.Results).To(HaveLen(2))
				
				// Check that we have results from both categories
				categories := make(map[string]bool)
				for _, res := range result.Results {
					categories[res.Category] = true
				}
				Expect(categories["science"]).To(BeTrue())
				Expect(categories["it"]).To(BeTrue())
			})

			It("should format results as JSON correctly", func() {
				response := &searxng.SearchResponse{
					Query:           "machine learning tutorial",
					NumberOfResults: 2500,
					Results: []searxng.SearchResult{
						{
							Title:         "Complete Machine Learning Tutorial",
							URL:           "https://example.com/ml-tutorial",
							Content:       "Learn machine learning from scratch...",
							PublishedDate: "2024-01-15",
							Category:      "science",
						},
						{
							Title:         "AI Development Guide",
							URL:           "https://tech.example.com/ai-guide",
							Content:       "A practical guide for developers...",
							PublishedDate: nil,
							Category:      "it",
						},
					},
				}

				jsonResult := tools.FormatSearchResultsJSON(response)
				
				Expect(jsonResult).To(ContainSubstring("machine learning tutorial"))
				Expect(jsonResult).To(ContainSubstring("Complete Machine Learning Tutorial"))
				Expect(jsonResult).To(ContainSubstring("AI Development Guide"))
				Expect(jsonResult).To(ContainSubstring("2024-01-15"))
				Expect(jsonResult).To(ContainSubstring("\"rank\": 1"))
				Expect(jsonResult).To(ContainSubstring("\"rank\": 2"))
				Expect(jsonResult).To(ContainSubstring("\"total\": 2500"))
			})
		})
		
		Context("with invalid categories", func() {
			It("should reject invalid category names", func() {
				_, err := searxng.ValidateCategory("invalid_category")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid category"))
			})
		})
	})

	Describe("Helper Functions", func() {
		Context("getAllCategoryNames", func() {
			It("should return all valid categories including science and it", func() {
				categories := searxng.GetAllCategories()
				categoryNames := make([]string, len(categories))
				for i, cat := range categories {
					categoryNames[i] = string(cat)
				}
				
				Expect(categoryNames).To(ContainElement("science"))
				Expect(categoryNames).To(ContainElement("it"))
				Expect(categoryNames).To(ContainElement("general"))
				Expect(categoryNames).To(ContainElement("images"))
				Expect(categoryNames).To(ContainElement("videos"))
			})
		})
	})
})
package searxng_test

import (
	"context"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"searxng-mcp/pkg/searxng"
)

var _ = Describe("SearXNG Client", func() {
	var (
		client searxng.Client
		ctx    context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
	})

	AfterEach(func() {
		// Cleanup if needed
	})

	Describe("HTTPClient", func() {
		Context("with valid configuration", func() {
			It("should create a new client with default URL", func() {
				httpClient := searxng.NewClient("")
				Expect(httpClient).NotTo(BeNil())
			})

			It("should create a new client with custom URL", func() {
				httpClient := searxng.NewClient("http://example.com:8888")
				Expect(httpClient).NotTo(BeNil())
			})
		})

		Context("with mock server", func() {
			var server *httptest.Server

			BeforeEach(func() {
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{
						"query": "test",
						"number_of_results": 1000,
						"results": [
							{
								"url": "https://example.com",
								"title": "Test Result",
								"content": "Test content",
								"thumbnail": "",
								"engine": "google",
								"template": "default.html",
								"parsed_url": ["https", "example.com", "/", "", "", ""],
								"img_src": "",
								"priority": "",
								"engines": ["google"],
								"positions": [1],
								"score": 10.0,
								"category": "general",
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

			It("should perform a successful search", func() {
				req := searxng.SearchRequest{
					Query:    "test",
					Category: []searxng.Category{searxng.CategoryGeneral},
				}

				result, err := client.Search(ctx, req)

				Expect(err).NotTo(HaveOccurred())
				Expect(result).NotTo(BeNil())
				Expect(result.Query).To(Equal("test"))
				Expect(result.NumberOfResults).To(Equal(1000))
				Expect(result.Results).To(HaveLen(1))
				Expect(result.Results[0].Title).To(Equal("Test Result"))
				Expect(result.Results[0].URL).To(Equal("https://example.com"))
			})

			It("should handle empty query", func() {
				req := searxng.SearchRequest{
					Query: "",
				}

				result, err := client.Search(ctx, req)

				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(err.Error()).To(ContainSubstring("query cannot be empty"))
			})
		})
	})

	Describe("Search Functions", func() {
		var server *httptest.Server

		BeforeEach(func() {
			server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{
					"query": "golang",
					"number_of_results": 5000,
					"results": [
						{
							"url": "https://go.dev",
							"title": "Go Programming Language",
							"content": "Go is an open source programming language",
							"thumbnail": "",
							"engine": "google",
							"template": "default.html",
							"parsed_url": ["https", "go.dev", "/", "", "", ""],
							"img_src": "",
							"priority": "",
							"engines": ["google"],
							"positions": [1],
							"score": 15.0,
							"category": "general",
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

		Context("SimpleSearch", func() {
			It("should perform a simple search", func() {
				result, err := searxng.SimpleSearch(ctx, client, "golang")

				Expect(err).NotTo(HaveOccurred())
				Expect(result).NotTo(BeNil())
				Expect(result.Query).To(Equal("golang"))
				Expect(result.NumberOfResults).To(Equal(5000))
				Expect(result.Results).To(HaveLen(1))
			})
		})

		Context("SearchWithCategory", func() {
			It("should perform a search with categories", func() {
				result, err := searxng.SearchWithCategory(ctx, client, "golang",
					searxng.CategoryGeneral, searxng.CategoryIT)

				Expect(err).NotTo(HaveOccurred())
				Expect(result).NotTo(BeNil())
				Expect(result.Query).To(Equal("golang"))
			})
		})

		Context("SearchWithOptions", func() {
			It("should perform a search with options", func() {
				opts := searxng.SearchOptions{
					Language:  "en",
					TimeRange: searxng.TimeRangeMonth,
					PageNo:    1,
				}

				result, err := searxng.SearchWithOptions(ctx, client, "golang", opts)

				Expect(err).NotTo(HaveOccurred())
				Expect(result).NotTo(BeNil())
				Expect(result.Query).To(Equal("golang"))
			})
		})
	})

	Describe("Validation Functions", func() {
		Context("ValidateCategory", func() {
			It("should validate valid categories", func() {
				category, err := searxng.ValidateCategory("general")

				Expect(err).NotTo(HaveOccurred())
				Expect(category).To(Equal(searxng.CategoryGeneral))
			})

			It("should reject invalid categories", func() {
				category, err := searxng.ValidateCategory("invalid")

				Expect(err).To(HaveOccurred())
				Expect(category).To(Equal(searxng.Category("")))
				Expect(err.Error()).To(ContainSubstring("invalid category"))
			})
		})

		Context("ValidateTimeRange", func() {
			It("should validate valid time ranges", func() {
				timeRange, err := searxng.ValidateTimeRange("month")

				Expect(err).NotTo(HaveOccurred())
				Expect(timeRange).To(Equal(searxng.TimeRangeMonth))
			})

			It("should reject invalid time ranges", func() {
				timeRange, err := searxng.ValidateTimeRange("invalid")

				Expect(err).To(HaveOccurred())
				Expect(timeRange).To(Equal(searxng.TimeRange("")))
				Expect(err.Error()).To(ContainSubstring("invalid time range"))
			})
		})
	})

	Describe("Helper Functions", func() {
		Context("GetAllCategories", func() {
			It("should return all available categories", func() {
				categories := searxng.GetAllCategories()

				Expect(categories).To(HaveLen(10))
				Expect(categories).To(ContainElement(searxng.CategoryGeneral))
				Expect(categories).To(ContainElement(searxng.CategoryImages))
				Expect(categories).To(ContainElement(searxng.CategoryVideos))
			})
		})

		Context("GetAllTimeRanges", func() {
			It("should return all available time ranges", func() {
				timeRanges := searxng.GetAllTimeRanges()

				Expect(timeRanges).To(HaveLen(3))
				Expect(timeRanges).To(ContainElement(searxng.TimeRangeDay))
				Expect(timeRanges).To(ContainElement(searxng.TimeRangeMonth))
				Expect(timeRanges).To(ContainElement(searxng.TimeRangeYear))
			})
		})
	})
})
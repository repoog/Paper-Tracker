package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"Paper-Tracker/db"
	"Paper-Tracker/trans"
)

type Feed struct {
	XMLName xml.Name `xml:"feed"`
	XMLNS   string   `xml:"xmlns,attr"`
	Entries []Entry  `xml:"entry"`
}

type Entry struct {
	Title     string `xml:"title"`
	Published string `xml:"published"`
	Updated   string `xml:"updated"`
	Summary   string `xml:"summary"`
	Link      []Link `xml:"link"`
}

type Link struct {
	Href string `xml:"href,attr"`
	Type string `xml:"type,attr"`
}

type TranslatedPaper struct {
	OriginalTitle     string
	TranslatedTitle   string
	OriginalSummary   string
	TranslatedSummary string
}

func main() {
	fmt.Println("[*] Start crawling the latest papers...")

	database, err := initDatabase()
	if err != nil {
		fmt.Printf("[!] Error initializing database: %v\n", err)
		return
	}
	defer database.Close()

	feed, err := fetchAndParseFeed()
	if err != nil {
		fmt.Printf("[!] Error fetching and parsing feed: %v\n", err)
		return
	}

	processFeed(database, feed)
}

func formatDate(dateStr string) string {
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return dateStr
	}
	return t.Format("2006-01-02 15:04:05")
}

// Clean string with no wrap and no extra spaces.
func cleanString(s string) string {
	re, err := regexp.Compile(`\s+`)
	if err != nil {
		fmt.Println("[!] Error cleaning strings.")
	}

	s = re.ReplaceAllString(s, " ")

	return strings.TrimSpace(s)
}

func initDatabase() (*db.Database, error) {
	database, err := db.ConnectDatabase("./db/arxiv_papers.db")
	if err != nil {
		return nil, fmt.Errorf("[!] Error opening database: %v", err)
	}

	if err := database.CreateTable(); err != nil {
		return nil, fmt.Errorf("[!] Error creating table: %v", err)
	}

	return database, nil
}

func fetchAndParseFeed() (*Feed, error) {
	paperUrl := "http://export.arxiv.org/api/query?search_query=cat:cs.CR&sortBy=lastUpdatedDate&sortOrder=descending&max_results=10"

	// Access and get paper url.
	resp, err := http.Get(paperUrl)
	if err != nil {
		return nil, fmt.Errorf("[!] Error fetching URL: %v", paperUrl)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[!] Error reading response body %v", err)
	}

	var feed Feed
	err = xml.Unmarshal(body, &feed)
	if err != nil {
		return nil, fmt.Errorf("[!] Error unmarshaling XML: %v", err)
	}

	if feed.Entries == nil {
		return nil, fmt.Errorf("[!] Entries is empty")
	}

	return &feed, nil
}

func processFeed(database *db.Database, feed *Feed) {
	for _, entry := range feed.Entries {
		paper := preparePaper(entry)

		exists, err := checkPaperExists(database, paper.Title)
		if err != nil {
			fmt.Printf("[!] Error checking if paper exists: %v\n", err)
			continue
		}
		if exists {
			fmt.Printf("[*] Paper already exists: %s\n", paper.Title)
			continue
		}

		if err := insertPaper(database, paper); err != nil {
			fmt.Printf("[!] Error inserting paper: %v\n", err)
			continue
		}

		printPaperInfo(paper)
	}
}

func preparePaper(entry Entry) *db.Paper {
	var pdfLink string
	for _, link := range entry.Link {
		if link.Type == "application/pdf" {
			pdfLink = link.Href
			break
		}
	}

	originalTitle := cleanString(entry.Title)
	originalSummay := cleanString(entry.Summary)
	translatedTitle, err := trans.Trans(originalTitle)
	if err != nil {
		fmt.Printf("[!] Error translating title: %v", err)
	}
	translatedSummary, err := trans.Trans(originalSummay)
	if err != nil {
		fmt.Printf("[!] Error translating summary: %v", err)
	}

	return &db.Paper{
		Title:      originalTitle,
		Title_CN:   translatedTitle,
		Published:  formatDate(entry.Published),
		Updated:    formatDate(entry.Updated),
		Link:       pdfLink,
		Summary:    originalSummay,
		Summary_CN: translatedSummary,
	}
}

func checkPaperExists(database *db.Database, title string) (bool, error) {
	return database.PaperExists(title)
}

func insertPaper(database *db.Database, paper *db.Paper) error {
	return database.InsertPaper(paper)
}

func printPaperInfo(paper *db.Paper) {
	fmt.Printf("标题: %s\n", paper.Title_CN)
	fmt.Printf("发布时间: %s\n", paper.Published)
	fmt.Printf("更新时间: %s\n", paper.Updated)
	fmt.Printf("论文链接: %v\n", paper.Link)
	fmt.Printf("摘要: %s\n", paper.Summary_CN)
	fmt.Println("----------------------------")
}

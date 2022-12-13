package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

// Book stores information
type Book struct {
	Title          string  `json:"title"`
	URL            string  `json:"detail"`
	Author         string  `json:"author"`
	ISBN           string  `json:"isbn"`
	Price          int64   `json:"price"`
	TotalPage      int64   `json:"number_of_page"`
	CollectionType string  `json:"category"`
	Name           string  `json:"book_nane"`
	Group          string  `json:"book_shelf"`
	Description    string  `json:"description"`
	Rating         float64 `json:"rate"`
}

func main() {
	site := flag.String("site", "nxbkimdong.com.vn", "input site to scrap")
	flag.Parse()

	config, err := LoadConfig(*site)
	if err != nil {
		log.Fatalf(err.Error())
	}

	fName := fmt.Sprintf("%s_books.json", *site)
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fName, err)
		return
	}
	defer file.Close()

	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: nxbkimdong.com.vn, www.nxbkimdong.com.vn
		colly.AllowedDomains(config.Site, fmt.Sprintf("www%s", config.Site)),

		// Cache responses to prevent multiple download of pages
		// even if the collector is restarted
		colly.CacheDir("./nxbkimdong_book_cache"),
	)

	// Create another collector to scrape book details
	detailCollector := c.Clone()

	books := make([]Book, 0, 200)

	// On every <a> element which has "href" attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {

		link := e.Attr("href")
		// If link is not all page then return
		if !strings.Contains(link, config.Collection) {
			return
		}

		// start scaping the page under the link found
		e.Request.Visit(link)
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
	})

	// On every <a> element with collection-product-card class call callback
	c.OnHTML(config.Product, func(e *colly.HTMLElement) {
		// Activate detailCollector if the link contains "nxbkimdong.com.vn/products"
		bookDetailUrl := e.Request.AbsoluteURL(e.Attr("href"))
		if strings.Contains(bookDetailUrl, config.UrlDetail) {
			log.Println("Load child: ", bookDetailUrl)
			detailCollector.Visit(bookDetailUrl)
		}
	})

	// Extract details of the book
	detailCollector.OnHTML(config.Section, func(e *colly.HTMLElement) {
		log.Println("Book found", e.Request.URL)

		book := Book{
			Title:          title(config, e),
			URL:            e.Request.URL.String(),
			Author:         author(config, e),
			ISBN:           isbn(config, e),
			Price:          price(config, e),
			TotalPage:      page(config, e),
			CollectionType: category(config, e),
			Name:           name(config, e),
			Group:          group(config, e),
			Description:    description(config, e),
			Rating:         rate(config, e),
		}

		log.Println("Book", book)
		books = append(books, book)
	})

	c.Visit(config.ScrapSite)

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")

	// Dump json to the standard output

	enc.Encode(books)

	log.Printf("Total books extracted %d", len(books))
}

var (
	nonAlphanumericRegex = regexp.MustCompile(`[^0-9]+`)
	nonisbn              = regexp.MustCompile(`[^0-9\\-]+`)
)

func ClearString(str string) string {
	return nonAlphanumericRegex.ReplaceAllString(str, "")
}

func toISBN(str string) string {
	return nonisbn.ReplaceAllString(str, "")
}

func title(config *StrategyParse, e *colly.HTMLElement) string {
	return extract(config.Title, e)
}

func author(config *StrategyParse, e *colly.HTMLElement) string {
	return extract(config.Author, e)
}

func isbn(config *StrategyParse, e *colly.HTMLElement) string {

	return toISBN(extract(config.ISBN, e))
}

func price(config *StrategyParse, e *colly.HTMLElement) int64 {
	return extractInt(config.Price, e)
}

func page(config *StrategyParse, e *colly.HTMLElement) int64 {
	return extractInt(config.Page, e)
}

func category(config *StrategyParse, e *colly.HTMLElement) string {
	return extract(config.Category, e)
}

func name(config *StrategyParse, e *colly.HTMLElement) string {
	return extract(config.Name, e)
}

func group(config *StrategyParse, e *colly.HTMLElement) string {
	return extract(config.Group, e)
}

func description(config *StrategyParse, e *colly.HTMLElement) string {
	return extract(config.Description, e)
}

func rate(config *StrategyParse, e *colly.HTMLElement) float64 {
	for _, k := range config.Rating {
		if text := e.ChildText(k); text != "" {
			rate, err := strconv.ParseFloat(ClearString(e.ChildText(k)), 64)
			if err != nil {
				rate = -1
			} else {
				log.Printf("extracted %v from %s\n", rate, k)
				return rate
			}

		}
	}
	return -1
}

func extract(keys []string, e *colly.HTMLElement) string {
	for _, k := range keys {
		if text := e.ChildText(k); text != "" {
			log.Printf("Extracted %s from %s\n", text, k)
			return text
		}
	}
	return ""
}

func extractInt(keys []string, e *colly.HTMLElement) int64 {
	for _, k := range keys {
		if text := e.ChildText(k); text != "" {
			page, err := strconv.ParseInt(ClearString(text), 10, 0)
			if err != nil {
				page = -1
			} else {
				log.Printf("Extracted %d from %s\n", page, k)
				return page
			}

		}
	}
	return -1
}

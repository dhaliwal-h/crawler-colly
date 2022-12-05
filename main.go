package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
)

type star struct {
	Name      string
	Photo     string
	JobTitle  string
	BirthDate string
	Bio       string
	TopMovies []movie
}

type movie struct {
	Title string
	Year  string
}

func main() {
	fmt.Println("Hello World")
	month := flag.Int("month", 1, "Month to fetch birthdays for")
	day := flag.Int("day", 1, "Day to fetch birthdays for")
	flag.Parse()
	crawl(*month, *day)
}

func crawl(month int, day int) {
	c := colly.NewCollector(
		colly.AllowedDomains("imdb.com", "www.imdb.com"),
	)

	infoC := c.Clone()

	c.OnHTML(".mode-detail", func(h *colly.HTMLElement) {
		profileUrl := h.ChildAttr("div.lister-item-image > a", "href")
		profileUrl = h.Request.AbsoluteURL(profileUrl)
		infoC.Visit(profileUrl)
	})

	c.OnHTML("a.lister-page-next", func(h *colly.HTMLElement) {
		nextPage := h.Request.AbsoluteURL(h.Attr("href"))
		c.Visit(nextPage)
	})

	infoC.OnHTML("#content-2-wide", func(h *colly.HTMLElement) {
		tmpProfile := star{}
		tmpProfile.Name = h.ChildText("h1.header > span.itemprop")
		tmpProfile.Photo = h.ChildAttr("#name-poster", "src")
		tmpProfile.JobTitle = h.ChildText("#name-job-categories > a > span.itemprop")
		tmpProfile.BirthDate = h.ChildAttr("#name-born-info time", "datetime")
		tmpProfile.Bio = strings.TrimSpace(h.ChildText("#name-bio-text > div.name-trivia-bio-text > div.inline"))
		h.ForEach("div.knownfor-title", func(_ int, kf *colly.HTMLElement) {
			tmpMovie := movie{}
			tmpMovie.Title = kf.ChildText("div.knownfor-title-role > a.knownfor-ellipsis")
			tmpMovie.Year = kf.ChildText("div.knownfor-year > span.knownfor-ellipsis")
			tmpProfile.TopMovies = append(tmpProfile.TopMovies, tmpMovie)
		})

		js, err := json.MarshalIndent(tmpProfile, "", "    ")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(js))
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("visiting:", r.URL.String())
	})

	infoC.OnRequest(func(r *colly.Request) {
		fmt.Println("visiting profile URL:", r.URL.String())
	})

	startUrl := fmt.Sprintf("https://www.imdb.com/search/name/?birth_monthday=%d-%d", month, day)
	c.Visit(startUrl)
}

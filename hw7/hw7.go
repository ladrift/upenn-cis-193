// Homework 7: Web Scraping
package main

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	// for _, slice := range ScrapeHackerNews(5) {
	// 	fmt.Println(slice)
	// }

	// for _, email := range GetEmails() {
	// 	fmt.Println(email)
	// }

	gdp, err := GetCountryGDP("Japan")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(gdp)
}

// News is a Hacker News article listing
type News struct {
	Points   int
	Title    string
	Username string
	URL      string
}

// NewsSlice is a slice of News pointers
type NewsSlice []*News

// ScrapeHackerNews scrapes the website "https://news.ycombinator.com/" using goquery and returns
// information on the first n posts.
//
// For each post, the attributes to be extracted are: points, title, username and url.
// This data should be returned as a NewsSlice, where NewsSlice is a custom slice of News structs.
//
// For example, for the sample image located at `https://www.cis.upenn.edu/~cis193/homeworks/hn.png`,
// the struct would look like:
// News{24, "QEMU(TCG): user-to-root privesc inside VM via bad translation caching",
// "webaholic", "https://bugs.chromium.org/p/project-zero/issues/detail?id=1122"}.
//
// If n is greater than the number of total posts available (which should be 30), return data from
// the all of the available posts (all thirty).
func ScrapeHackerNews(n int) NewsSlice {
	site := "https://news.ycombinator.com/"
	doc, err := goquery.NewDocument(site)
	if err != nil {
		log.Fatal(err)
	}

	news := make([]*News, n)
	sel := doc.Find(".athing")
	for i := range sel.Nodes[:n] {
		s := sel.Eq(i)

		link := s.Find("td.title").Find("a.storylink")
		title := link.Text()
		url, _ := link.Attr("href")

		sub := s.Next().Find("td.subtext")
		user := sub.Find("a.hnuser").Text()

		rawScore := sub.Find("span.score").Text()
		score, err := strconv.Atoi(strings.Split(rawScore, " ")[0])
		if err != nil {
			log.Fatal(err)
		}

		news[i] = &News{score, title, user, url}
	}

	return NewsSlice(news)
}

// GetEmails returns a string slice of the emails found on the given URL.
//
// Scenario: you are a student enthusiastic about spreading awareness about Go. To effectively
// market Go, you decide to email Penn CIS professors about the wonders of the Go programming
// language. In this function, use goquery to extract the email addresses from the URL
// "http://www.cis.upenn.edu/about-people/" and return them as a string slice. This will involve you
// having to investigate where and how emails are located on the webpage.
// Note: you should have 47 total emails returned.
func GetEmails() []string {
	site := "http://www.cis.upenn.edu/about-people/"
	doc, err := goquery.NewDocument(site)
	if err != nil {
		log.Fatal(err)
	}

	row := doc.Find("tbody").Children().First().Next()
	if err != nil {
		log.Fatal(err)
	}

	emails := make([]string, 0)
	for ; row != nil && goquery.NodeName(row) == "tr"; row = row.Next().Next() {
		html, err := row.Html()
		if err != nil {
			log.Fatal(err)
		}

		emailRe := regexp.MustCompile(`[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,6}`)
		emails = append(emails, emailRe.FindString(html))
	}

	return emails
}

// CountryData has GDP information on a country
type CountryData struct {
	Country string
	GDP     string
}

// GetCountryGDP takes in a string country name and returns the GDP (in millions) as
// an integer. Information on the country is found by concurrently scraping a hidden website with
// data on countries scattered on many pages.
//
// Scenario: imagine you are a spy and you have discovered a URL with top secret GDP information:
// "https://www.cis.upenn.edu/~cis193/scraping/9828772efc2bd314a277c8880695dea2.html". This webpage
// has a country name and the GDP (in millions of US Dollars). It also has links to two other
// country's webpages. Based on intelligence you've received, every country has a webpage on this
// website with information about it, but you do not know the URL for each page. You can assume that
// none of the page links lead you to a cycle and every country can be reached from a path from the
// initial URL that you are given. So, for this function, you will need to traverse from the initial
// url to every webpage link you encounter in order to find information on the target `country`
// string. Since time is of the essence, you want to use concurrency to scrape webpages
// simultaneously. Note that for this function, we only care about getting the GDP for the input
// `country` string. You may find it useful to use the CountryData struct to send country
// information between goroutines.
//
// To prevent the function from getting stuck if an invalid `country` string is entered,
// you should also implement a timeout that will automatically return an error after 10 seconds
// if the program hasn't already finished terminating.
//
// Feel free to make and use helper functions for this function. To help with testing this
// function, we know from intelligence reports that the GDP for "Canada" is 1532343 and
// the GDP for "Colombia" is 274135.
func GetCountryGDP(country string) (int, error) {
	root := "https://www.cis.upenn.edu/~cis193/scraping/9828772efc2bd314a277c8880695dea2.html"

	data, errc, cancelFunc := scrapeCountryData(root)
	defer cancelFunc()

	for {
		select {
		case countryData, ok := <-data:
			if !ok {
				return 0, fmt.Errorf("country not found")
			}

			if country == countryData.Country {
				reg := regexp.MustCompile("[^a-zA-Z0-9]+")
				gdp, err := strconv.Atoi(reg.ReplaceAllString(countryData.GDP, ""))
				if err != nil {
					return 0, err
				}
				return gdp, nil
			}
		case err := <-errc:
			return 0, err
		}
	}
}

func scrapeCountryData(root string) (<-chan CountryData, <-chan error, context.CancelFunc) {
	NumScraper := runtime.NumCPU()

	data := make(chan CountryData)
	wg := new(sync.WaitGroup)
	wg.Add(NumScraper)
	go func() {
		wg.Wait()
		close(data)
	}()

	errc := make(chan error)
	ctx, cancelFunc := context.WithCancel(context.Background())

	urls := make(chan string, 1)
	urls <- root
	for i := 0; i < NumScraper; i++ {
		go func() {
			countryDataScraper(ctx, urls, data, errc)
			wg.Done()
		}()
	}

	return data, errc, cancelFunc
}

func countryDataScraper(ctx context.Context, urls chan string, data chan<- CountryData, errc chan<- error) {
	for url := range urls {
		doc, err := goquery.NewDocument(url)
		if err != nil {
			select {
			case errc <- err:
				return
			case <-ctx.Done():
				return
			}
		}

		country := doc.Find("h3.country").Text()
		gdp := doc.Find("h3.gdp").Text()
		select {
		case data <- CountryData{country, gdp}:
		case <-ctx.Done():
			return
		}

		links := doc.Find("a")
		for i := range links.Nodes {
			sel := links.Eq(i)

			url, exists := sel.Attr("href")
			if !exists {
				select {
				case errc <- fmt.Errorf("invalid <a> tag"):
					return
				case <-ctx.Done():
					return
				}
			}

			select {
			case urls <- url:
			case <-ctx.Done():
				return
			}
		}
	}
}

package search

import (
	"coffeeintocode/search-engine/db"
	"fmt"
	"time"
)

func RunEngine() {
	fmt.Println("started search engine crawl...")
	defer fmt.Println("search engine crawl has finished")
	// Get crawl settings from DB
	settings := &db.SearchSettings{}
	err := settings.Get()
	if err != nil {
		fmt.Println("something went wrong getting the settings")
		return
	}
	// Check if search is turned on by checking settings
	if !settings.SearchOn {
		fmt.Println("search is turned off")
		return
	}
	crawl := &db.CrawledUrl{}
	// Get next X urls to be tested
	nextUrls, err := crawl.GetNextCrawlUrls(int(settings.Amount))
	if err != nil {
		fmt.Println("something went wrong getting the url list")
		return
	}
	newUrls := []db.CrawledUrl{}
	testedTime := time.Now()
	// Loop over the slice and run crawl on each url
	for _, next := range nextUrls {
		result := runCrawl(next.Url)
		// Check if the crawl was not successul
		if !result.Success {
			// Update row in database with the failed crawl
			err := next.UpdateUrl(db.CrawledUrl{
				ID:              next.ID,
				Url:             next.Url,
				Success:         false,
				CrawlDuration:   result.CrawlData.CrawlTime,
				ResponseCode:    result.ResponseCode,
				PageTitle:       result.CrawlData.PageTitle,
				PageDescription: result.CrawlData.PageDescription,
				Headings:        result.CrawlData.Headings,
				LastTested:      &testedTime,
			})
			if err != nil {
				fmt.Println("something went wrong updating a failed url")
			}
			continue
		}
		// Update a successful row in database
		err := next.UpdateUrl(db.CrawledUrl{
			ID:              next.ID,
			Url:             next.Url,
			Success:         result.Success,
			CrawlDuration:   result.CrawlData.CrawlTime,
			ResponseCode:    result.ResponseCode,
			PageTitle:       result.CrawlData.PageTitle,
			PageDescription: result.CrawlData.PageDescription,
			Headings:        result.CrawlData.Headings,
			LastTested:      &testedTime,
		})
		if err != nil {
			fmt.Printf("something went wrong updating %v /n", next.Url)
		}
		// Push the newly found external urls to an array
		for _, newUrl := range result.CrawlData.Links.External {
			newUrls = append(newUrls, db.CrawledUrl{Url: newUrl})
		}
	} // End of range
	// Check if we should add the newly found urls to the database
	if !settings.AddNew {
		fmt.Printf("Adding new urls to database is disabled")
		return
	}
	// Insert newly found urls into database
	for _, newUrl := range newUrls {
		err := newUrl.Save()
		if err != nil {
			fmt.Printf("something went wrong adding new url to database: %v", newUrl.Url)
		}
	}
	fmt.Printf("\nAdded %d new urls to database \n", len(newUrls))

}

func RunIndex() {
	fmt.Println("started search indexing...")
	defer fmt.Println("search indexing has finished")
	// Get index settings from DB
	crawled := &db.CrawledUrl{}
	// Get all urls that are not indexed
	notIndexed, err := crawled.GetNotIndexed()
	fmt.Println("not indexed urls: ", len(notIndexed))
	if err != nil {
		fmt.Println("something went wrong getting the not indexed urls")
		return
	}
	// Create a new index
	idx := make(Index)
	// Add the not indexed urls to the index
	idx.Add(notIndexed)
	// Save the index to the database
	searchIndex := &db.SearchIndex{}
	err = searchIndex.Save(idx, notIndexed)
	if err != nil {
		fmt.Println(err)
		fmt.Println("something went wrong saving the index")
		return
	}
	// Update the urls to be indexed=true
	err = crawled.SetIndexedTrue(notIndexed)
	if err != nil {
		fmt.Println("something went wrong updating the indexed urls")
		return
	}

}

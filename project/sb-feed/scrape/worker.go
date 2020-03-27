package scrape

import (
	"github.com/MartinNikolovMarinov/sb-feed/scrape/persistance"
	mysqlpers "github.com/MartinNikolovMarinov/sb-feed/scrape/persistance/mysql"
	"github.com/MartinNikolovMarinov/sb-infra/entities"
	"github.com/MartinNikolovMarinov/sb-infra/infra"
)

// OddsFeed comment
type OddsFeed interface {
	/*
		Fetch should do all networking related calls to a specific data source API.
		This might include combining all needed data into an internal structure,
		specific for different sources.
		The end result should be an array of Game entities.
	*/
	Fetch() ([]*entities.Game, error)
}

// FeedWorker comment
type FeedWorker struct {
	feeds            []OddsFeed
	connectionString string
}

// NewFeedWorker comment
func NewFeedWorker(connectionString string) *FeedWorker {
	return &FeedWorker{
		feeds:            make([]OddsFeed, 0),
		connectionString: connectionString,
	}
}

// AddFeed commnet
func (fw *FeedWorker) AddFeed(feed OddsFeed) {
	fw.feeds = append(fw.feeds, feed)
}

// Work commnet
func (fw *FeedWorker) Work() {
	if fw.feeds == nil || len(fw.feeds) == 0 {
		infra.Warn("No feeds!")
		return
	}

	var feedLen = len(fw.feeds)

	db := &mysqlpers.DatabaseMySQL{}
	db.Init(fw.connectionString)
	dl := persistance.NewDataLayer(db)
	defer db.Close()

	for i := 0; i < feedLen; i++ {
		feed := fw.feeds[i]

		games, err := feed.Fetch()
		if err != nil {
			infra.Error(err.Error())
			return
		}

		persistGames(dl, games)
	}
}

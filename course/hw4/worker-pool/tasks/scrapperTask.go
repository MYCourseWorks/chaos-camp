package tasks

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	rindexdb "github.com/homework/hw4/reverse-index-db"
	"github.com/homework/hw4/util"
	workpool "github.com/homework/hw4/worker-pool"
)

// ScrapperTask comment
type ScrapperTask struct {
	TimeToWaitRequest time.Duration
	visitedURLS       map[string]bool
	mux               sync.Mutex
}

// NewScrapperTask comment
func NewScrapperTask(timeToWaitRequest time.Duration) *ScrapperTask {
	return &ScrapperTask{
		TimeToWaitRequest: timeToWaitRequest,
		visitedURLS:       make(map[string]bool),
	}
}

var tr *http.Transport = &http.Transport{
	Dial: (&net.Dialer{
		KeepAlive: 1 * time.Second,
		Timeout:   1 * time.Second,
		Deadline:  time.Now().Add(5 * time.Second),
	}).Dial,
	MaxIdleConns:        1,
	MaxIdleConnsPerHost: 1,
	DisableKeepAlives:   true,
	TLSHandshakeTimeout: 1 * time.Second,
	IdleConnTimeout:     1 * time.Second,
}
var client http.Client = http.Client{
	Transport: tr,
	Timeout:   1 * time.Second,
}

// Execute comment
func (st *ScrapperTask) Execute(workItem *workpool.WorkerItem) ([]*workpool.WorkerItem, error) {
	url := workItem.Value

	st.mux.Lock()
	if _, ok := st.visitedURLS[url]; ok {
		// already visited this site
		return nil, nil
	}

	fmt.Printf("Visiting: %s\n", url)
	st.visitedURLS[url] = true
	// FIXME: the defer here makes the code worse than single threaded
	st.mux.Unlock()

	// Receive response :
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "Scrapper")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code %d for URL %s", resp.StatusCode, url)
	}

	h := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(h, "text/html") {
		// Response not html
		return nil, nil
	}

	// Traverse Tokens in response body :
	textInTags, urlsInHTML, err := util.ExtractTextAndUrls(resp.Body)
	if err != nil {
		return nil, err
	}

	// Map urls to worker items :
	workerItems := make([]*workpool.WorkerItem, len(urlsInHTML))
	for i, u := range urlsInHTML {
		workerItems[i] = &workpool.WorkerItem{
			Value:    u,
			Priority: (workItem.Priority + 1),
		}
	}

	// Store scrapped document text and url :
	err = rindexdb.GetInstance().StoreDocument(url, textInTags)
	if err != nil {
		return nil, err
	}

	return workerItems, nil
}

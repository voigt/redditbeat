package beater

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"
	"github.com/jzelinskie/geddit"

	"github.com/voigt/redditbeat/config"
)

type Redditbeat struct {
	done       chan struct{}
	config     config.Config
	client     publisher.Client
	session    *geddit.LoginSession
	collecting bool
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	config := config.DefaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	session, err := geddit.NewLoginSession(
		*config.Reddit.Username,
		*config.Reddit.Password,
		*config.Reddit.Useragent,
	)
	if err != nil {
		return nil, err
	}

	bt := &Redditbeat{
		done:    make(chan struct{}),
		config:  config,
		session: session,
	}

	return bt, nil
}

func (bt *Redditbeat) Run(b *beat.Beat) error {
	logp.Info("redditbeat is running! Hit CTRL-C to stop it.")

	bt.client = b.Publisher.Connect()
	ticker := time.NewTicker(bt.config.Period)
	counter := 1
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
			if !bt.collecting {
				err := bt.collectPosts()
				if err != nil {
					return err
				}
			}
		}

		// event := common.MapStr{
		// 	"@timestamp": common.Time(time.Now()),
		// 	"type":       b.Name,
		// 	"counter":    counter,
		// }
		// bt.client.PublishEvent(event)
		//logp.Info("Event sent")
		counter++
	}
}

func (bt *Redditbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}

func (bt *Redditbeat) collectPosts() error {
	logp.Info("Collecting posts")

	bt.collecting = true
	defer func() {
		bt.collecting = false
	}()

	sync, err, processed := make(chan byte), make(chan error), 0

	for _, sub := range *bt.config.Reddit.Subs {
		fmt.Println(sub)
		go bt.processSubs(sub, sync, err)
	}

	for {
		select {
		case <-sync:
			processed++
			if processed == len(*bt.config.Reddit.Subs) {
				return nil
			}
		case e := <-err:
			return e
		}
	}

	return nil

}

func (bt *Redditbeat) processSubs(name string, sync chan byte, err chan error) {
	logp.Info("Collecting Submissions for '%v'", name)

	result, e := bt.fetchSubmissions(name)

	if e != nil {
		logp.Critical("Error while fetching Submissions happend: %v", e)
	}

	logp.Info("Got %v Submissions for '%v'", len(result), name)

	for _, submission := range result {

		b, err := json.Marshal(submission)
		if err != nil {
			logp.Critical("Could not read JSON: %v", err)
			return
		}

		fmt.Println(string(b))
		event := common.MapStr{
			"@timestamp": common.Time(time.Now()),
			"type":       "submission",
			"subName":    name,
			"submission": submission,
		}

		bt.client.PublishEvent(event)
		logp.Info("Event sent")
	}

	// if len(result) >= 1 {
	// 	bt.redditMap.Set(name, strconv.FormatInt(result[0].Id, 10))
	// }

	sync <- 1
}

func (bt *Redditbeat) fetchSubmissions(name string) ([]*geddit.Submission, error) {

	subOpts := geddit.ListingOptions{
		Limit: *bt.config.Reddit.Limit,
	}

	logp.Info("Fetching... '%v'", name)
	submissions, err := bt.session.SubredditSubmissions(name, geddit.NewSubmissions, subOpts)

	return submissions, err
}

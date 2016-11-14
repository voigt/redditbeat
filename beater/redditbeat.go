package beater

import (
	"flag"
	"fmt"
	"strconv"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"
	"github.com/jzelinskie/geddit"

	"github.com/voigt/redditbeat/config"
	"github.com/voigt/redditbeat/persistency"
)

// Redditbeat is a beater object. Contains all objects needed to run the beat.
type Redditbeat struct {
	done        chan struct{}
	config      config.Config
	CmdLineArgs CmdLineArgs
	client      publisher.Client
	session     *geddit.LoginSession
	collecting  bool
	redditMap   *persistency.StringMap
}

// CmdLineArgs is a helper struct for adding custom command line flags
type CmdLineArgs struct {
	PersistencyMap *string
}

var cmdLineArgs CmdLineArgs

func init() {
	cmdLineArgs = CmdLineArgs{
		PersistencyMap: flag.String("p", "redditmap.json", "Path to the persistency map json file"),
	}
}

// New Creates a new beater
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

	redditMap := persistency.NewStringMap()
	redditMap.Load(*cmdLineArgs.PersistencyMap)

	bt := &Redditbeat{
		done:      make(chan struct{}),
		config:    config,
		session:   session,
		redditMap: redditMap,
	}

	return bt, nil
}

// Run defines whats happening when the beat is running
func (bt *Redditbeat) Run(b *beat.Beat) error {
	logp.Info("redditbeat is running! Hit CTRL-C to stop it.")

	bt.client = b.Publisher.Connect()
	ticker := time.NewTicker(bt.config.Period)

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
	}
}

// Stop makes sure the beat is closed properly
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
	logp.Info("[%v] Collecting Submissions for '%v'", name, name)

	var latest float32 // latest contains the timestamp of the last event published
	var newCount int   // counts how much new Submissions have been published

	latest, e := bt.getLatestSubmissionTimestamp(name)
	if e != nil {
		logp.Critical("[%v] Error while getting latest submission timestamp: %v", name, e)
	}

	result, e := bt.fetchSubmissions(name)
	if e != nil {
		logp.Critical("[%v] Error while fetching Submissions: %v", name, e)
	}

	logp.Critical("Resultlength: %v", len(result))
	// Iterating resultset backwards, to check oldest submission first
	for i := len(result) - 1; i >= 0; i-- {
		submission := result[i]

		// If the created Timestamp of the submission is newer than latest, publish this submission as event
		if submission.DateCreated > latest {

			logp.Info("[%v] %v > %v", name, submission.DateCreated, latest)

			event := common.MapStr{
				"@timestamp": common.Time(time.Now()),
				"type":       "submission",
				"subName":    name,
				"submission": submission,
			}

			bt.client.PublishEvent(event)
			logp.Info("[%v] Event sent", name)
			logp.Info("[%v] Persisting new Timestamp of %v", name, submission.FullID)

			newCount++
			latest = submission.DateCreated
			dateCreated := float64(submission.DateCreated)
			bt.redditMap.Set(name, strconv.FormatFloat(dateCreated, 'E', -1, 32))
		}
	}

	logp.Info("[%v] Got %v new Submissions for %v", name, newCount, name)

	sync <- 1
}

// Check whether a timestamp for the Subreddit `name` has already been saved
// If so, use this as reference in order to publish only new submissions
// else create mapping for the new Subreddit
func (bt *Redditbeat) getLatestSubmissionTimestamp(name string) (float32, error) {

	var latest float32

	if bt.redditMap.Contains(name) {
		latestString := bt.redditMap.Get(name)

		latestFloat, e := strconv.ParseFloat(latestString, 32)
		if e != nil {
			logp.Critical("[%v] Could not convert to Float 'latest' Timestamp %v", name, e)
			return 0, e
		}

		latest = float32(latestFloat)
		// logp.Info("[%v] latest Timestamp: %v", name, latest)
	} else {
		latest = 0.0
		bt.redditMap.Set(name, strconv.FormatFloat(float64(latest), 'E', -1, 32))
	}

	return latest, nil
}

// Fetches the X last submissions of a given Subreddit.
// X is defined as "Limit" in beat.yml
func (bt *Redditbeat) fetchSubmissions(name string) ([]*geddit.Submission, error) {

	subOpts := geddit.ListingOptions{
		Limit: *bt.config.Reddit.Limit,
	}

	logp.Info("[%v] Fetching... '%v'", name, name)
	submissions, err := bt.session.SubredditSubmissions(name, geddit.NewSubmissions, subOpts)

	return submissions, err
}

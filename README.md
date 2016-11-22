# Redditbeat

Redditbeat is an [Elastic Beat](https://github.com/elastic/beats) to index new Reddit Submissions of one or multiple Subreddits.

## Getting Started with Redditbeat

### Requirements

Ensure that this folder is at the following location: `${GOPATH}/github.com/voigt`

* [Golang](https://golang.org/dl/) 1.7
* [Glide](https://github.com/Masterminds/glide) (it's used during `make scaffold`)

### Clone Project

```
git clone https://github.com/voigt/redditbeat.git $GOPATH/src/github.com/voigt/redditbeat
cd $GOPATH/src/github.com/voigt/redditbeat
```

### Init Project
To get running with Redditbeat and also install the
dependencies, run the following command:

```
make scaffold
```

### Build

To build the binary for Redditbeat run the command below. This will generate a binary
in the same directory with the name redditbeat.

```
make
```


### Run

To run Redditbeat with debugging output enabled, run:

```
./redditbeat -c redditbeat.yml -p data/redditmap.json -e -d "*"
```

**Hint:** If you want to reindex already indexed Subreddits (resets data/redditmap.json):

```
make clear-cache
```

### Configuration

You'll want to configure which Subreddits to index. You will do this in `redditbeat.yml`.

```
redditbeat:
  # Defines how often an event is sent to the output
  period: 60s                       # how often to check for new Submissions

  reddit:
    username: "username"
    password: "password"
    useragent: "Redditbeat v0.1"
    subs: ["kitten", "news"]        # a list of Subreddits to index
    limit: 10                       # curret limit is 100
```

---

# Todo

* [x] index new Submissions of one or multiple given Subreddits
* [x] add persistency, so already indexed submissions will not be indexed again
* [ ] add dockerfile `make package`
* [ ] index new Submissions of one or multiple Users

---

# Known issues

* **Redditbeat misses some new Submissions**  
Redditbeat is making use of [geddit](https://github.com/jzelinskie/geddit). Unfortunately geddit saves the timestamp of a submission in `float32`, which means we lose up to 99 seconds of the timestamp. Ultimately this leads to the fact, that Redditbeat does not recognise new Submissions of which created date is closer than 99 secs. geddit is [already informed](https://github.com/jzelinskie/geddit/issues/25).  

---

# Thanks to

* [@buehler](https://github.com/buehler), I used his [twitterbeat](https://github.com/buehler/twitterbeat) as a pattern. I copied the persistency approach.

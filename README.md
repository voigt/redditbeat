# Redditbeat

Welcome to Redditbeat.

Ensure that this folder is at the following location:
`${GOPATH}/github.com/voigt`

## Getting Started with Redditbeat

### Requirements

* [Golang](https://golang.org/dl/) 1.7

### Init Project
To get running with Redditbeat and also install the
dependencies, run the following command:

```
make setup
```

It will create a clean git history for each major step. Note that you can always rewrite the history if you wish before pushing your changes.

To push Redditbeat in the git repository, run the following commands:

```
git remote set-url origin https://github.com/voigt/redditbeat
git push origin master
```

For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).

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


### Test

To test Redditbeat, run the following command:

```
make testsuite
```

alternatively:
```
make unit-tests
make system-tests
make integration-tests
make coverage-report
```

The test coverage is reported in the folder `./build/coverage/`

### Update

Each beat has a template for the mapping in elasticsearch and a documentation for the fields
which is automatically generated based on `etc/fields.yml`.
To generate etc/redditbeat.template.json and etc/redditbeat.asciidoc

```
make update
```


### Cleanup

To clean  Redditbeat source code, run the following commands:

```
make fmt
make simplify
```

To clean up the build directory and generated artifacts, run:

```
make clean
```


### Clone

To clone Redditbeat from the git repository, run the following commands:

```
mkdir -p ${GOPATH}/github.com/voigt
cd ${GOPATH}/github.com/voigt
git clone https://github.com/voigt/redditbeat
```


For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).


## Packaging

The beat frameworks provides tools to crosscompile and package your beat for different platforms. This requires [docker](https://www.docker.com/) and vendoring as described above. To build packages of your beat, run the following command:

```
make package
```

This will fetch and create all images required for the build process. The hole process to finish can take several minutes.

---

# Known issues

* **Redditbeat misses some new Submissions**  
Redditbeat is making use of [geddit](https://github.com/jzelinskie/geddit). Unfortunately geddit saves the timestamp of a submission in `float32`, which means we lose up to 99 seconds of the timestamp. Ultimately this leads to the fact, that Redditbeat does not recognise new Submissions of which created date is closer than 99 secs. geddit is [already informed](https://github.com/jzelinskie/geddit/issues/25).  

---

# Thanks to

* [@buehler](https://github.com/buehler), I used his [twitterbeat](https://github.com/buehler/twitterbeat) as a pattern. I copied the persistency approach.
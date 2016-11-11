package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/voigt/redditbeat/beater"
)

func main() {
	err := beat.Run("redditbeat", "", beater.New)
	if err != nil {
		os.Exit(1)
	}
}

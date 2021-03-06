package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/nifibeat/beater"
)

func main() {
	err := beat.Run("nifibeat", "", beater.New)
	if err != nil {
		os.Exit(1)
	}
}

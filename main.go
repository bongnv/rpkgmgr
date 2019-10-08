package main

import (
	"github.com/bongnv/rpkgmgr/logic"
	"github.com/bongnv/rpkgmgr/repository"
	"github.com/robfig/cron/v3"
)

const (
	// define the number of concurrent jobs to download package and parse DESCRIPTION
	maxConcurrency = 5
	cranURL        = "http://cran.r-project.org/src/contrib/"
)

func main() {
	repo := repository.NewRepository()
	descIndexer := logic.NewDescriptionIndexer(cranURL, maxConcurrency, repo)
	indexer := logic.NewIndexer(cranURL, descIndexer, repo)
	indexer.Run()

	c := cron.New()
	c.AddJob("0 12 * * *", indexer)
	c.Start()
}

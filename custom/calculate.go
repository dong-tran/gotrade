package custom

import (
	"log"
	"time"

	"github.com/dong-tran/gotrade/asset"
	"github.com/dong-tran/gotrade/helper"
	"github.com/dong-tran/gotrade/strategy"
	"github.com/dong-tran/gotrade/strategy/momentum"
	"github.com/dong-tran/gotrade/strategy/trend"
	"github.com/dong-tran/gotrade/strategy/volume"
)

func Calculate(symbols []string) {
	var csvPath = "csv/"
	var strategies = []strategy.Strategy{trend.NewAroonStrategy(), momentum.NewStochasticRsiStrategy(), volume.NewEaseOfMovementStrategyWith(4)}
	repository := asset.NewFileSystemRepository(csvPath)
	since := time.Now().AddDate(0, 0, -30)
	for _, name := range symbols {
		snapshots, err := repository.GetSince(name, since)
		if err != nil {
			log.Panicf("Unable to retrieve snapshots. %v", err)
			continue
		}
		for _, currentStrategy := range strategies {
			actions := currentStrategy.Compute(snapshots)
			lastAction := helper.Last(actions, 1)
			l := <-lastAction
			log.Printf("%s: %s", name, l.Annotation())
		}
	}
}

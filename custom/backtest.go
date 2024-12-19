package custom

import (
	"log"
	"os"
	"time"

	"github.com/dong-tran/gotrade/asset"
	"github.com/dong-tran/gotrade/backtest"
	"github.com/dong-tran/gotrade/strategy"
	"github.com/dong-tran/gotrade/strategy/decorator"
	"github.com/dong-tran/gotrade/strategy/momentum"
	"github.com/dong-tran/gotrade/strategy/trend"
	"github.com/dong-tran/gotrade/strategy/volume"
)

type DailyData struct {
	Date     string
	Open     float64
	High     float64
	Low      float64
	Close    float64
	AdjClose float64
	Volume   float64
}

func Backtest(days int, symbols []string) {
	// report := backtest.NewDataReport()
	var reportPath = "/tmp/report"
	err := os.RemoveAll(reportPath + "/")
	if err != nil {
		log.Fatalf("Error when clean report output folder: %v", err)
	}

	var buyStrategy = strategy.NewAndStrategy("BoP-KDJ-FI", trend.NewBopStrategy(), trend.NewKdjStrategy(), volume.NewForceIndexStrategy())
	var sellStrategy = strategy.NewAndStrategy("AOsi-FI", momentum.NewAwesomeOscillatorStrategy(), volume.NewForceIndexStrategy())
	var combined = stopLoss(strategy.NewSplitStrategy(buyStrategy, sellStrategy))
	var secondStrategy = stopLoss(strategy.NewAndStrategy("BoP-RSI (Sideway - Down)", trend.NewBopStrategy(), momentum.NewRsiStrategy()))

	// backtest.Strategies = append(backtest.Strategies, trend.NewBopStrategy())
	// backtest.Strategies = append(backtest.Strategies, trend.NewKamaStrategy())
	// backtest.Strategies = append(backtest.Strategies, trend.NewKdjStrategy())

	// backtest.Strategies = append(backtest.Strategies, trend.NewQstickStrategy())
	// backtest.Strategies = append(backtest.Strategies, trend.NewCciStrategy())
	// backtest.Strategies = append(backtest.Strategies, trend.NewAroonStrategy())
	// backtest.Strategies = append(backtest.Strategies, trend.NewVwmaStrategy())
	// backtest.Strategies = append(backtest.Strategies, trend.NewEnvelopeStrategy())
	// backtest.Strategies = append(backtest.Strategies, trend.NewMacdStrategy())
	// backtest.Strategies = append(backtest.Strategies, trend.NewTrimaStrategy())
	// backtest.Strategies = append(backtest.Strategies, trend.NewTrixStrategy())
	// backtest.Strategies = append(backtest.Strategies, trend.NewDemaStrategy())
	// backtest.Strategies = append(backtest.Strategies, trend.NewTsiStrategy())
	// backtest.Strategies = append(backtest.Strategies, trend.NewTripleMovingAverageCrossoverStrategy())
	// backtest.Strategies = append(backtest.Strategies, trend.NewGoldenCrossStrategy())

	// backtest.Strategies = append(backtest.Strategies, momentum.NewAwesomeOscillatorStrategy())
	// backtest.Strategies = append(backtest.Strategies, momentum.NewStochasticRsiStrategy())
	// backtest.Strategies = append(backtest.Strategies, momentum.NewRsiStrategy())
	// backtest.Strategies = append(backtest.Strategies, momentum.NewTripleRsiStrategy())

	// backtest.Strategies = append(backtest.Strategies, volume.NewForceIndexStrategy())
	// backtest.Strategies = append(backtest.Strategies, volume.NewEaseOfMovementStrategy())
	// backtest.Strategies = append(backtest.Strategies, volume.NewChaikinMoneyFlowStrategy())
	// backtest.Strategies = append(backtest.Strategies, volume.NewMoneyFlowIndexStrategy())
	// backtest.Strategies = append(backtest.Strategies, volume.NewNegativeVolumeIndexStrategy())
	// backtest.Strategies = append(backtest.Strategies, volume.NewVolumeWeightedAveragePriceStrategy())
	doReport(reportPath, "swd", symbols, days, secondStrategy, false)
	doReport(reportPath, "", symbols, days, combined, false)
	if time.Now().Weekday() == time.Friday {
		wd := NewWeekUtil(730)
		doReport(reportPath, "week", symbols, wd.GetNextStart(), combined, true)
	}
}

func stopLoss(stg strategy.Strategy) strategy.Strategy {
	return decorator.NewStopLossStrategy(stg, 0.1)
}

func doReport(reportPath string, subdir string, symbols []string, days int, stg strategy.Strategy, isWeek bool) {
	var csvPath = "/tmp/csv"
	var err error
	report := backtest.NewHTMLReportWith(reportPath, subdir)
	var repository asset.Repository
	if isWeek {
		repository = asset.NewFileSystemWeekRepository(csvPath)
	} else {
		repository = asset.NewFileSystemRepository(csvPath)
	}
	backtest := backtest.NewBacktest(repository, report)
	backtest.Names = append(backtest.Names, symbols...)
	backtest.Workers = 5
	backtest.LastDays = days
	backtest.Strategies = append(backtest.Strategies, stg)
	err = backtest.Run()
	// PrintDataReport(symbols, report)
	if err != nil {
		log.Fatal(err)
	}
}

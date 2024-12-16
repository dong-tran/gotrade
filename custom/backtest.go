package custom

import (
	"log"
	"os"

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
	var csvPath = "/tmp/csv"
	var err error

	// report := backtest.NewDataReport()
	var reportPath = "/tmp/report"
	err = os.RemoveAll(reportPath + "/")
	if err != nil {
		log.Fatalf("Error when clean report output folder: %v", err)
	}
	report := backtest.NewHTMLReport(reportPath)
	if err != nil {
		log.Fatalf("Error when create report: %v", err)
	}
	repository := asset.NewFileSystemRepository(csvPath)
	backtest := backtest.NewBacktest(repository, report)
	backtest.Names = append(backtest.Names, symbols...)
	backtest.Workers = 1
	backtest.LastDays = days
	//
	backtest.Strategies = append(backtest.Strategies, stopLoss(strategy.NewAndStrategy("Good On UpTrend", trend.NewBopStrategy(), volume.NewForceIndexStrategy(), momentum.NewAwesomeOscillatorStrategy())))
	backtest.Strategies = append(backtest.Strategies, stopLoss(strategy.NewAndStrategy("Good On DownTrend", trend.NewBopStrategy(), momentum.NewRsiStrategy())))
	backtest.Strategies = append(backtest.Strategies, stopLoss(strategy.NewAndStrategy("MACD Stochastic", trend.NewMacdStrategy(), momentum.NewStochasticRsiStrategy())))

	var buyStrategy = strategy.NewAndStrategy("Bop-Kdj-FI", trend.NewBopStrategy(), trend.NewKdjStrategy(), volume.NewForceIndexStrategy())
	var sellStrategy = strategy.NewAndStrategy("AOsi-FI", momentum.NewAwesomeOscillatorStrategy(), volume.NewForceIndexStrategy())
	var combined = strategy.NewSplitStrategy(buyStrategy, sellStrategy)
	backtest.Strategies = append(backtest.Strategies, stopLoss(combined))

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

	err = backtest.Run()
	// PrintDataReport(symbols, report)
	// }
	if err != nil {
		log.Fatal(err)
	}
}

func stopLoss(stg strategy.Strategy) strategy.Strategy {
	return decorator.NewStopLossStrategy(stg, 0.1)
}

package custom

import (
	"fmt"

	"github.com/dong-tran/gotrade/backtest"
	"github.com/dong-tran/gotrade/helper"
)

func PrintDataReport(symbols []string, report *backtest.DataReport) {
	for _, r := range symbols {
		rs := report.Results[r]
		s := r + ":"
		for i := 0; i < len(rs); i++ {
			actionsSplice := helper.SliceToChan(rs[i].Transactions)
			actions := helper.Last(actionsSplice, 1)
			s = s + (<-actions).Annotation()
		}
		fmt.Println(s)
	}
}

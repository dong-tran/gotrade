// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package strategy_test

import (
	"testing"

	"github.com/dong-tran/gotrade/helper"
	"github.com/dong-tran/gotrade/strategy"
)

func TestOutcome(t *testing.T) {
	values := helper.SliceToChan([]float64{
		10, 15, 12, 12, 18,
		20, 22, 25, 24, 20,
	})

	actions := helper.SliceToChan([]strategy.Action{
		strategy.Hold, strategy.Hold, strategy.Buy, strategy.Buy, strategy.Hold,
		strategy.Hold, strategy.Hold, strategy.Sell, strategy.Hold, strategy.Hold,
	})

	expected := helper.SliceToChan([]float64{
		0, 0, 0, 0, 0.5,
		0.67, 0.83, 1.08, 1.08, 1.08,
	})

	actual := helper.RoundDigits(strategy.Outcome(values, actions), 2)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

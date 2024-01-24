package precompiles

import (
	"math/big"
	"testing"
)

func TestArbMathStatsStandardDeviationIdenticalElements(t *testing.T) {
	mathStats := ArbMathStats{}

	data := []int32{1, 1, 1, 1, 1, 1}
	stdev, _ := mathStats.stdDev(nil, nil, data, 0)

	if stdev.Cmp(big.NewInt(0)) != 0 {
		t.Errorf("Expected stdev of 0 for identical elements, got: %s", stdev)
	}
}

func TestArbMathStatsStandardDeviation(t *testing.T) {
	mathStats := ArbMathStats{}

	data := []int32{1, 5, 10, 100}

	stdev, _ := mathStats.stdDev(nil, nil, data, 4)
	expected := big.NewInt(474763)
	if stdev.Cmp(expected) != 0 {
		t.Errorf("Incorrect stddev, expected %s, got %s", expected, stdev)
	}

	stdev, _ = mathStats.stdDev(nil, nil, data, 2)
	expected = big.NewInt(4747)
	if stdev.Cmp(expected) != 0 {
		t.Errorf("Incorrect stddev, expected %s, got %s", expected, stdev)
	}
}

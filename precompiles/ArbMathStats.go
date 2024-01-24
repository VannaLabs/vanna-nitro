package precompiles

import (
	"math"
	"math/big"

	"gonum.org/v1/gonum/stat"
)

// ArbMathStats provides methods for computing statistics
type ArbMathStats struct {
	Address addr // 0x11b
}

func (con *ArbMathStats) stdDev(c ctx, evm mech, input []int32, decimals uint8) (huge, error) {
	floatInput := make([]float64, len(input))
	for i, num := range input {
		floatInput[i] = float64(num)
	}

	stdev := big.NewFloat(stat.StdDev(floatInput, nil))
	stdev.Mul(stdev, big.NewFloat(math.Pow10(int(decimals))))

	truncatedStdev := new(big.Int)
	stdev.Int(truncatedStdev)

	return truncatedStdev, nil
}

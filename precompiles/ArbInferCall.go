// Copyright 2021-2022, Offchain Labs, Inc.
// For license information, see https://github.com/nitro/blob/master/LICENSE

package precompiles

import (
	"fmt"
	"strings"

	inference "github.com/ethereum/go-ethereum/rpc/inference"
)

// ArbGasInfo provides insight into the cost of using the rollup.
type ArbInferCall struct {
	Address addr // 0x11a
}

func (con *ArbInferCall) InferCall(c ctx, evm mech, input []byte) ([]byte, error) {
	inputStr := string(input)
	// Split string into two parts
	inputArray := strings.Split(inputStr, "-")
	modelName := inputArray[0]
	inputData := inputArray[1]
	rc := inference.NewRequestClient(5125)
	tx := inference.InferenceTx{
		Hash:   "0x123456789",
		Model:  modelName,
		Params: inputData,
		TxType: "inference",
	}
	result, err := rc.Emit(tx)
	if err != nil {
		fmt.Println("InferCall Error", err)
		return []byte{}, err
	}

	byteValue := make([]byte, len(result))
	copy(byteValue, result)
	return byteValue, nil
}

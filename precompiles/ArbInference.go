package precompiles

import (
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/rpc/inference"
)

// ArbInference provides precompiles that serve as an entrypoint for AI and ML inference on the Vanna Blockchain
type ArbInference struct {
	Address addr // 0x11a
}

func (con *ArbInference) InferCall(c ctx, evm mech, input []byte) ([]byte, error) {
	inputStr := string(input)
	inputArray := strings.Split(inputStr, "-")
	modelName := inputArray[0]
	inputData := ""
	for _, value := range inputArray[1:] {
		inputData = inputData + value
	}

	rc := inference.NewRequestClient(5125)
	tx := inference.InferenceTx{
		Hash:   c.caller.String() + "#" + strconv.FormatUint(evm.StateDB.GetNonce(c.caller), 10),
		Model:  modelName,
		Params: inputData,
		TxType: "inference",
	}
	result, err := rc.Emit(tx)
	if err != nil {
		return []byte{}, err
	}

	byteValue := make([]byte, len(result))
	copy(byteValue, result)
	return byteValue, nil
}

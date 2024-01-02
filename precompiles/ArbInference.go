package precompiles

import (
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/rpc/inference"
	. "github.com/ethereum/go-ethereum/rpc/inference"
)

// ArbInference provides precompiles that serve as an entrypoint for AI and ML inference on the Vanna Blockchain
type ArbInference struct {
	Address addr // 0x11a
}

const (
	InputSeperator = "<?>"
	TXSeparator    = "#"
)

func (con *ArbInference) InferCall(c ctx, evm mech, input []byte) ([]byte, error) {
	inputStr := string(input)
	inputArray := strings.Split(inputStr, InputSeperator)
	modelName := inputArray[0]
	inputData := inputArray[1]

	rc := inference.NewRequestClient(5125)
	caller := c.txProcessor.Callers[0]
	tx := inference.InferenceTx{
		Hash:   caller.String() + TXSeparator + strconv.FormatUint(evm.StateDB.GetNonce(caller), 10),
		Model:  modelName,
		Params: inputData,
		TxType: Inference,
	}
	result, err := rc.Emit(tx)
	if err != nil {
		return []byte{}, err
	}

	byteValue := make([]byte, len(result))
	copy(byteValue, result)
	return byteValue, nil
}

func (con *ArbInference) InferCallZK(c ctx, evm mech, input []byte) ([]byte, error) {
	inputStr := string(input)
	inputArray := strings.Split(inputStr, InputSeperator)
	modelName := inputArray[0]
	inputData := inputArray[1]

	rc := inference.NewRequestClient(5125)
	caller := c.txProcessor.Callers[0]
	tx := inference.InferenceTx{
		Hash:   caller.String() + TXSeparator + strconv.FormatUint(evm.StateDB.GetNonce(caller), 10),
		Model:  modelName,
		Params: inputData,
		TxType: ZKInference,
	}
	result, err := rc.Emit(tx)
	if err != nil {
		return []byte{}, err
	}

	byteValue := make([]byte, len(result))
	copy(byteValue, result)
	return byteValue, nil
}

func (con *ArbInference) InferCallPipeline(c ctx, evm mech, input []byte) ([]byte, error) {
	inputStr := string(input)
	inputArray := strings.Split(inputStr, InputSeperator)
	modelName := inputArray[0]
	pipelineName := inputArray[1]
	seed := inputArray[2]
	inputData := inputArray[3]

	rc := inference.NewRequestClient(5125)
	caller := c.txProcessor.Callers[0]
	tx := inference.InferenceTx{
		Hash:     caller.String() + TXSeparator + strconv.FormatUint(evm.StateDB.GetNonce(caller), 10),
		Seed:     seed,
		Pipeline: pipelineName,
		Model:    modelName,
		Params:   inputData,
		TxType:   PipelineInference,
	}
	result, err := rc.Emit(tx)
	if err != nil {
		return []byte{}, err
	}

	byteValue := make([]byte, len(result))
	copy(byteValue, result)
	return byteValue, nil
}

func (con *ArbInference) InferCallPrivate(c ctx, evm mech, input []byte) ([]byte, error) {
	inputStr := string(input)
	inputArray := strings.Split(inputStr, InputSeperator)
	IPAddress := inputArray[0]
	modelName := inputArray[1]
	inputData := inputArray[2]

	rc := inference.NewRequestClient(5125)
	caller := c.txProcessor.Callers[0]
	tx := inference.InferenceTx{
		Hash:   caller.String() + TXSeparator + strconv.FormatUint(evm.StateDB.GetNonce(caller), 10),
		Model:  modelName,
		Params: inputData,
		TxType: PrivateInference,
		IP:     IPAddress,
	}
	result, err := rc.Emit(tx)
	if err != nil {
		return []byte{}, err
	}

	byteValue := make([]byte, len(result))
	copy(byteValue, result)
	return byteValue, nil
}

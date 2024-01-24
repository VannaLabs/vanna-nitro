package precompiles

import (
	// #nosec G501
	"crypto/md5"
	"encoding/hex"
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
	InputSeparator = "<?>"
	TXSeparator    = "#"
)

func (con *ArbInference) InferCall(c ctx, evm mech, input []byte) ([]byte, error) {
	inputStr := string(input)
	inputArray := strings.Split(inputStr, InputSeparator)
	modelName := inputArray[0]
	inputData := inputArray[1]

	rc := inference.NewRequestClient(5125)
	caller := c.txProcessor.Callers[0]
	tx := inference.InferenceTx{
		Hash:   HashInferenceTX([]string{caller.String(), strconv.FormatUint(evm.StateDB.GetNonce(caller), 10), modelName, inputData}, TXSeparator),
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
	inputArray := strings.Split(inputStr, InputSeparator)
	modelName := inputArray[0]
	inputData := inputArray[1]

	rc := inference.NewRequestClient(5125)
	caller := c.txProcessor.Callers[0]
	tx := inference.InferenceTx{
		Hash:   HashInferenceTX([]string{caller.String(), strconv.FormatUint(evm.StateDB.GetNonce(caller), 10), modelName, inputData}, TXSeparator),
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
	inputArray := strings.Split(inputStr, InputSeparator)
	modelName := inputArray[0]
	pipelineName := inputArray[1]
	seed := inputArray[2]
	inputData := inputArray[3]

	rc := inference.NewRequestClient(5125)
	caller := c.txProcessor.Callers[0]
	tx := inference.InferenceTx{
		Hash:     HashInferenceTX([]string{caller.String(), strconv.FormatUint(evm.StateDB.GetNonce(caller), 10), modelName, inputData}, TXSeparator),
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
	inputArray := strings.Split(inputStr, InputSeparator)
	IPAddress := inputArray[0]
	modelName := inputArray[1]
	inputData := inputArray[2]

	rc := inference.NewRequestClient(5125)
	caller := c.txProcessor.Callers[0]
	tx := inference.InferenceTx{
		Hash:   HashInferenceTX([]string{caller.String(), strconv.FormatUint(evm.StateDB.GetNonce(caller), 10), modelName, inputData}, TXSeparator),
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

func HashInferenceTX(arr []string, separator string) string {
	hashString := ""
	for i := 0; i < len(arr); i++ {
		hashString += arr[i] + separator
	}

	// #nosec G401
	hasher := md5.New()
	hasher.Write([]byte(hashString))
	return hex.EncodeToString(hasher.Sum(nil))
}

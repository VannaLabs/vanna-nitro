package precompiles

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"

	"github.com/ethereum/go-ethereum/rpc/inference"
	. "github.com/ethereum/go-ethereum/rpc/inference"
)

// ArbInference provides precompiles that serve as an entrypoint for AI and ML inference on the Vanna Blockchain
type ArbInference struct {
	Address addr // 0x11a
}

const (
	TXSeparator = "#"
)

func (con *ArbInference) InferCall(c ctx, evm mech, model []byte, input []byte) ([]byte, error) {
	modelName := string(model)
	inputData := string(input)

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

func (con *ArbInference) InferCallZK(c ctx, evm mech, model []byte, input []byte) ([]byte, error) {
	modelName := string(model)
	inputData := string(input)

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

func (con *ArbInference) InferCallPipeline(c ctx, evm mech, model []byte, pipeline []byte, seed []byte, input []byte) ([]byte, error) {
	modelName := string(model)
	pipelineName := string(pipeline)
	inferenceSeed := string(seed)
	inputData := string(input)

	rc := inference.NewRequestClient(5125)
	caller := c.txProcessor.Callers[0]
	tx := inference.InferenceTx{
		Hash:     HashInferenceTX([]string{caller.String(), strconv.FormatUint(evm.StateDB.GetNonce(caller), 10), modelName, inputData}, TXSeparator),
		Seed:     inferenceSeed,
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

func (con *ArbInference) InferCallPrivate(c ctx, evm mech, ip []byte, model []byte, input []byte) ([]byte, error) {
	IPAddress := string(ip)
	modelName := string(model)
	inputData := string(input)

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
	hasher := md5.New()
	hasher.Write([]byte(hashString))
	return hex.EncodeToString(hasher.Sum(nil))
}

package precompiles

import (
	"github.com/ethereum/go-ethereum/rpc/inference"
)

// ArbInference provides precompiles that serve as an entrypoint for AI and ML inference on the Vanna Blockchain
type ArbInference struct {
	model  string
	params string
}

func (con ArbStatistics) InferCall(c ctx, evm mech) (string, error) {
	rc := inference.NewRequestClient(5125)
	tx := inference.InferenceTx{
		Hash:   "0x123456789",
		Model:  "QmXQpupTphRTeXJMEz3BCt9YUF6kikcqExxPdcVoL1BBhy",
		Params: "[[0.002, 0.005, 0.004056685]]",
		TxType: "inference",
	}
	// Return value should be ("0.0013500629", nil)
	return rc.Emit(tx)
}

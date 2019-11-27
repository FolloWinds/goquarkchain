package native

import (
	"github.com/QuarkChain/goquarkchain/account"
	"github.com/QuarkChain/goquarkchain/core"
	"github.com/stretchr/testify/assert"
	"math/big"

	//"github.com/QuarkChain/goquarkchain/account"
	"github.com/QuarkChain/goquarkchain/common"
	//"github.com/QuarkChain/goquarkchain/core"
	//"math/big"
	"testing"
)

func TestNativeTokenTransfer(t *testing.T) {
	QETH := common.TokenIDEncode("QETH")
	id1, _ := account.CreatRandomIdentity()
	acc1 := account.CreatAddressFromIdentity(id1, 0)
	acc2 := account.CreatEmptyAddress(0)
	acc3 := account.CreatEmptyAddress(0)
	core.TestGenesisMinorTokenBalance["QKC"] = big.NewInt(10000000)
	core.TestGenesisMinorTokenBalance["QETH"] = big.NewInt(99999)
	env := core.GetTestEnv(&acc1, nil, nil, nil, nil, nil)
	shardState := core.CreateDefaultShardState(env, nil, nil, nil, nil)
	val := big.NewInt(12345)
	gas := uint64(21000)
	gasPrice := uint64(1)
	tx1 := core.CreateTransferTransaction(shardState, id1.GetKey().Bytes(), acc1, acc2, val, &gas, &gasPrice, nil, nil, nil, &QETH)
	error := shardState.AddTx(tx1)
	if error != nil {
		t.Errorf("addTx error: %v", error)
	}
	b1, _ := shardState.CreateBlockToMine(nil, &acc3, nil, nil, nil)
	assert.Equal(t, len(b1.Transactions()), 1)
	shardState.FinalizeAndAddBlock(b1)
	a := shardState.CurrentBlock().Header()
	b := b1.Header()
	assert.Equal(t, a, b)

}

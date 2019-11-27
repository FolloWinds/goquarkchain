package native

import (
	"github.com/QuarkChain/goquarkchain/account"
	"github.com/QuarkChain/goquarkchain/core"
	"github.com/QuarkChain/goquarkchain/core/types"
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
	assert.Equal(t, shardState.CurrentHeader(), b1.Header())
	tTxList, _, err := shardState.GetTransactionByAddress(acc1, nil, nil, 0)
	if err != nil {
		//t.Errorf("GetTransactionByAddress error :%v", err)
		t.Log(tTxList)
	}
	//assert.NotEqual(t, tTxList[0].Value, &serialize.Uint256{Value: big.NewInt(12345)})
}

func TestNativeTokenTransferValueSuccess(t *testing.T) {
	MALICIOUS0 := common.TokenIDEncode("MALICIOUS0")
	id1, _ := account.CreatRandomIdentity()
	acc1 := account.CreatAddressFromIdentity(id1, 0)
	acc3 := account.CreatEmptyAddress(0)
	core.TestGenesisMinorTokenBalance["QKC"] = big.NewInt(10000000)
	core.TestGenesisMinorTokenBalance["MALICIOUS0"] = big.NewInt(0)
	env := core.GetTestEnv(&acc1, nil, nil, nil, nil, nil)
	shardState := core.CreateDefaultShardState(env, nil, nil, nil, nil)
	val := big.NewInt(0)
	gas := uint64(21000)
	gasPrice := uint64(1)
	tx1 := core.CreateTransferTransaction(shardState, id1.GetKey().Bytes(), acc1, acc1, val, &gas, &gasPrice, nil, nil, nil, &MALICIOUS0)
	error := shardState.AddTx(tx1)
	if error != nil {
		t.Errorf("addTx error: %v", error)
	}
	b1, _ := shardState.CreateBlockToMine(nil, &acc3, nil, nil, nil)
	assert.Equal(t, len(b1.Transactions()), 1)
	shardState.FinalizeAndAddBlock(b1)
	assert.Equal(t, shardState.CurrentHeader(), b1.Header())
	bl, _ := shardState.GetBalance(id1.GetRecipient(), nil)
	QKC := common.TokenIDEncode("QKC")
	assert.Equal(t, bl.GetTokenBalance(QKC), big.NewInt(10000000-21000))
	b2, _ := shardState.GetBalance(acc1.Recipient, nil)
	assert.Equal(t, b2.GetTokenBalance(MALICIOUS0), big.NewInt(0))
	t1 := types.NewTokenBalancesWithMap(map[uint64]*big.Int{
		MALICIOUS0: big.NewInt(0),
		QKC:        big.NewInt(10000000 - 21000),
	})
	assert.NotEqual(t, b2, t1)
	t2 := types.NewTokenBalancesWithMap(map[uint64]*big.Int{
		QKC: big.NewInt(10000000 - 21000),
	})
	assert.Equal(t, b2, t2)
}

func TestDisallowedUnknownToken(t *testing.T) {
	MALICIOUS0 := common.TokenIDEncode("MALICIOUS0")
	MALICIOUS1 := common.TokenIDEncode("MALICIOUS1")
	id1, _ := account.CreatRandomIdentity()
	acc1 := account.CreatAddressFromIdentity(id1, 0)
	core.TestGenesisMinorTokenBalance["QKC"] = big.NewInt(10000000)
	env := core.GetTestEnv(&acc1, nil, nil, nil, nil, nil)
	shardState := core.CreateDefaultShardState(env, nil, nil, nil, nil)
	val := big.NewInt(0)
	gas := uint64(21000)
	gasPrice := uint64(1)
	tx1 := core.CreateTransferTransaction(shardState, id1.GetKey().Bytes(), acc1, acc1, val, &gas, &gasPrice, nil, nil, nil, &MALICIOUS0)
	assert.Error(t, shardState.AddTx(tx1))
	tx2 := core.CreateTransferTransaction(shardState, id1.GetKey().Bytes(), acc1, acc1, val, &gas, &gasPrice, nil, nil, nil, &MALICIOUS1)
	assert.Error(t, shardState.AddTx(tx2))
}

func TestNativeTokenGas(t *testing.T) {
	QETH := common.TokenIDEncode("QETH")
	id1, _ := account.CreatRandomIdentity()
	acc1 := account.CreatAddressFromIdentity(id1, 0)
	acc2 := account.CreatEmptyAddress(0)
	acc3 := account.CreatEmptyAddress(0)
	core.TestGenesisMinorTokenBalance["QETH"] = big.NewInt(10000000)
	core.TestGenesisMinorTokenBalance["QKC"] = big.NewInt(10000000)
	env := core.GetTestEnv(&acc1, nil, nil, nil, nil, nil)
	shardState := core.CreateDefaultShardState(env, nil, nil, nil, nil)
	val := big.NewInt(12345)
	gas := uint64(21000)
	tx1 := core.CreateTransferTransaction(shardState, id1.GetKey().Bytes(), acc1, acc2, val, &gas, nil, nil, nil, &QETH, &QETH)
	assert.NoError(t, shardState.AddTx(tx1))
	b1, _ := shardState.CreateBlockToMine(nil, &acc3, nil, nil, nil)
	assert.Equal(t, len(b1.Transactions()), 1)
	shardState.FinalizeAndAddBlock(b1)
	assert.Equal(t, shardState.CurrentHeader(), b1.Header())
	bl, _ := shardState.GetBalance(acc1.Recipient, nil)
	b2, _ := shardState.GetBalance(acc2.Recipient, nil)
	bb := new(big.Int).Sub(big.NewInt(10000000), big.NewInt(12345))
	assert.Equal(t, bl.GetTokenBalance(QETH), new(big.Int).Sub(bb, big.NewInt(21000)))
	assert.NotEqual(t, b2.GetTokenBalance(QETH), big.NewInt(12345))
}

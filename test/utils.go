// Package testutils provides general purpose utility functions for unit/integration testing.
package testutils

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"testing"

	geth "github.com/ethereum/go-ethereum"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	abci "github.com/tendermint/tendermint/abci/types"
	evmTypes "github.com/axelarnetwork/axelar-core/x/evm/types"
)

// Func wraps a regular testing function so it can be used as a pointer function receiver
type Func func(t *testing.T)

// Repeat executes the testing function n times
func (f Func) Repeat(n int) Func {
	return func(t *testing.T) {
		for i := 0; i < n; i++ {
			f(t)
		}
	}
}

// Events wraps sdk.Events
type Events []abci.Event

// Filter returns a collection of events filtered by the predicate
func (fe Events) Filter(predicate func(events abci.Event) bool) Events {
	var filtered Events
	for _, event := range fe {
		if predicate(event) {
			filtered = append(filtered, event)
		}
	}
	return filtered
}

// ErrorCache is a struct that can be used to get at the error that is emitted by test assertions when passing it instead ot *testing.T
type ErrorCache struct {
	Error error
}

// Errorf records the given formatted string as an erro
func (ec *ErrorCache) Errorf(format string, args ...interface{}) {
	ec.Error = fmt.Errorf(format, args...)
}

// SetEnv safely sets an OS env var to the specified value and resets it to the original value upon test closure
func SetEnv(t *testing.T, key string, val string) {
	// TODO : enable with Go 1.17 >> it will automatically handle Cleanup
	//t.Setenv(key, val)
	orig := os.Getenv(key)
	os.Setenv(key, val)
	t.Cleanup(func() { os.Setenv(key, orig) })
}

// CreateDeployGatewayTx assembles a transaction for smart contract deployment. See:
// https://goethereumbook.org/en/smart-contract-deploy/
// https://gist.github.com/tomconte/6ce22128b15ba36bb3d7585d5180fba0
func CreateDeployGatewayTx(
	byteCode []byte,
	contractOwner, contractOperator common.Address,
	gasPrice *big.Int,
	gasLimit uint64,
	evmTypes types.RPCClient,
) (*gethTypes.Transaction, error) {
	nonce, err := rpc.PendingNonceAt(context.Background(), contractOwner)
	if err != nil {
		return nil, err
	}

	if gasPrice == nil {
		gasPrice, err = rpc.SuggestGasPrice(context.Background())
		if err != nil {
			return nil, err
		}
	}

	deploymentBytecode, err := evmTypes.GetGatewayDeploymentBytecode(byteCode, contractOperator)
	if err != nil {
		return nil, err
	}

	if gasLimit == 0 {
		gasLimit, err = rpc.EstimateGas(context.Background(), geth.CallMsg{
			To:   nil,
			Data: deploymentBytecode,
		})
	}

	return gethTypes.NewContractCreation(nonce, big.NewInt(0), gasLimit, gasPrice, deploymentBytecode), nil
}


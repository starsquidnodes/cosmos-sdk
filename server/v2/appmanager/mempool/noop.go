package mempool

import (
	"context"

	"github.com/cosmos/cosmos-sdk/server/v2/core/mempool"
	"github.com/cosmos/cosmos-sdk/server/v2/core/transaction"
)

// TODO here until we rebase to get the txcodec items
type TxValidator[T transaction.Tx] interface {
	ValidateTx(ctx context.Context, tx T, simulate bool) (context.Context, error)
}

var _ mempool.Mempool[transaction.Tx] = NoOpMempool[transaction.Tx]{}

// NoOpMempool[T] is a mempool that only validates transactions
type NoOpMempool[T transaction.Tx] struct {
	txValidator TxValidator[T]
}

func NewNoopMempool[T transaction.Tx](txv TxValidator[T]) *NoOpMempool[T] {
	return &NoOpMempool[T]{txValidator: txv}
}

func (s *NoOpMempool[T]) Start() error {
	// NoOpMempool[T] does not require any initialization
	return nil
}

func (s *NoOpMempool[T]) Stop() error {
	// NoOpMempool[T] does not require any cleanup
	return nil
}

func (npm NoOpMempool[T]) Insert(ctx context.Context, tx T) error {
	_, err := npm.txValidator.ValidateTx(ctx, tx, false)
	return err
}

func (NoOpMempool[T]) GetTxs(ctx context.Context, size uint32, txSizeFn mempool.TxSizeFn) (any, error) {
	return nil, nil
}

func (NoOpMempool[T]) Remove(any) error { return nil }

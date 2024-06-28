package tx

import (
	"context"
	_ "embed"
	"encoding/hex"
	"testing"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	stakingtypes "cosmossdk.io/x/staking/types"

	"github.com/cosmos/cosmos-sdk/runtime"

	"cosmossdk.io/core/address"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/x/auth/ante"
	authtypes "cosmossdk.io/x/auth/types"
	"cosmossdk.io/x/gov"
	"cosmossdk.io/x/staking"

	codectestutil "github.com/cosmos/cosmos-sdk/codec/testutil"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	//go:embed testdata/signed_v0.47.tx
	msgEditValidatorV047 string
	//go:embed testdata/signed_v0.50.tx
	msgEditValidatorV050 string
	//go:embed testdata/signed_latest.tx
	msgEditValidatorLatest string
)

func TestRegression(t *testing.T) {
	enc := testutil.MakeTestEncodingConfig(codectestutil.CodecOptions{},
		gov.AppModule{}, staking.AppModule{})
	addressCodec := enc.InterfaceRegistry.SigningContext().AddressCodec()

	specs := map[string]struct {
		src string
	}{
		"MsgEditValidator - v047": {
			src: msgEditValidatorV047,
		},
		"MsgEditValidator - v050": {
			src: msgEditValidatorV050,
		},
		"MsgEditValidator - latest": {
			src: msgEditValidatorLatest,
		},
	}
	for name, spec := range specs {
		t.Run(name, func(t *testing.T) {
			tx, err := enc.TxConfig.TxJSONDecoder()([]byte(spec.src))
			require.NoError(t, err)
			require.NotNil(t, tx)
			require.Len(t, tx.GetMsgs(), 1)

			a := ante.NewSigVerificationDecorator(accountKeeperMock{codec: addressCodec}, enc.TxConfig.SignModeHandler(), ante.DefaultSigVerificationGasConsumer, nil)
			ctx := sdk.Context{}.
				WithChainID("testchain").
				WithBlockHeight(1).
				WithGasMeter(storetypes.NewInfiniteGasMeter()).
				WithEventManager(sdk.NewEventManager())

			var nextCalled bool
			captureNext := func(ctx sdk.Context, tx sdk.Tx, simulate bool) (sdk.Context, error) {
				nextCalled = true
				return ctx, nil
			}
			// when
			_, err = a.AnteHandle(ctx, tx, false, captureNext)

			// then
			require.NoError(t, err)
			assert.True(t, nextCalled)
		})
	}
}

type accountKeeperMock struct {
	codec address.Codec
}

func (a accountKeeperMock) GetEnvironment() appmodule.Environment {
	return runtime.NewEnvironment(runtime.NewKVStoreService(storetypes.NewKVStoreKey(stakingtypes.StoreKey)), log.NewNopLogger())
}

func (a accountKeeperMock) GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI {
	if len(addr) == 0 {
		panic("invalid address")
	}
	r := authtypes.NewBaseAccountWithAddress(addr)
	testkey := &secp256k1.PrivKey{Key: must(hex.DecodeString(`4b1184d5e958060419c9b24f7b42e40619554fcd3bc635c066395306e946bcf8`))}
	_ = r.SetPubKey(testkey.PubKey())
	r.AccountNumber = 1
	r.Sequence = 2
	return r
}

func (a accountKeeperMock) GetParams(ctx context.Context) (params authtypes.Params) {
	return authtypes.DefaultParams()
}
func (a accountKeeperMock) SetAccount(ctx context.Context, acc sdk.AccountI) {}

func (a accountKeeperMock) NewAccountWithAddress(ctx context.Context, addr sdk.AccAddress) sdk.AccountI {
	panic("implement me1")
}

func (a accountKeeperMock) GetModuleAddress(moduleName string) sdk.AccAddress {
	panic("implement me5")
}

func (a accountKeeperMock) AddressCodec() address.Codec {
	return a.codec
}

func must[T any](r T, err error) T {
	if err != nil {
		panic(err)
	}
	return r
}

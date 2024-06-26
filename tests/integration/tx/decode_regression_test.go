package tx

import (
	"context"
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

func TestRegression(t *testing.T) {
	enc := testutil.MakeTestEncodingConfig(codectestutil.CodecOptions{},
		gov.AppModule{}, staking.AppModule{})
	addressCodec := enc.InterfaceRegistry.SigningContext().AddressCodec()

	specs := map[string]struct {
		src string
	}{
		"MsgCreateValidator - v047": {
			src: `{"body":{"messages":[{"@type":"/cosmos.staking.v1beta1.MsgCreateValidator","description":{"moniker":"simapp-test","identity":"","website":"","security_contact":"","details":""},"commission":{"rate":"0.990000000000000000","max_rate":"1.000000000000000000","max_change_rate":"0.100000000000000000"},"min_self_delegation":"100000","delegator_address":"cosmos1a40r995peauyekrkda2uky5h52j0r0ztfut34a","validator_address":"cosmosvaloper1a40r995peauyekrkda2uky5h52j0r0ztvglyew","pubkey":{"@type":"/cosmos.crypto.ed25519.PubKey","key":"4gdqY8r0epCoDDpLrVwo7sWAVwPeqf9PEX4LFECMH3g="},"value":{"denom":"stake","amount":"1000000000"}}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[{"public_key":{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"AymTGWvRvySz1XavUhrz6GIoN8ZyVo278i1TAUfrlJWZ"},"mode_info":{"single":{"mode":"SIGN_MODE_LEGACY_AMINO_JSON"}},"sequence":"2"}],"fee":{"amount":[],"gas_limit":"200000","payer":"","granter":""},"tip":null},"signatures":["CJ/2V5zfl3PB/3CEcYbS9y2HTl+oxw098M64VSeIzH18PuBHN4+ou/NzvBKTWMGdTGPOB/AEN1EYj8TxyjApmA=="]}`,
		},
		"MsgEditValidator - v047": {
			src: `{"body":{"messages":[{"@type":"/cosmos.staking.v1beta1.MsgEditValidator","description":{"moniker":"[do-not-modify]","identity":"[do-not-modify]","website":"[do-not-modify]","security_contact":"[do-not-modify]","details":"[do-not-modify]"},"validator_address":"cosmosvaloper1a40r995peauyekrkda2uky5h52j0r0ztvglyew","commission_rate":"0.250000000000000000","min_self_delegation":null}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[{"public_key":{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"AymTGWvRvySz1XavUhrz6GIoN8ZyVo278i1TAUfrlJWZ"},"mode_info":{"single":{"mode":"SIGN_MODE_LEGACY_AMINO_JSON"}},"sequence":"2"}],"fee":{"amount":[],"gas_limit":"200000","payer":"","granter":""},"tip":null},"signatures":["9u4PviP9rBSymBOztKF+S9s9Kk50juyKX7B+RXdIdRtlHfJZ2pbzizc8X61x0oFJU75Ob4oVbtfnzNjQmdFiQA=="]}`,
		},
		"MsgVoteWeighted - v047": {
			src: `{"body":{"messages":[{"@type":"/cosmos.gov.v1.MsgVoteWeighted","proposal_id":"1","voter":"cosmos1a40r995peauyekrkda2uky5h52j0r0ztfut34a","options":[{"option":"VOTE_OPTION_YES","weight":"0.600000000000000000"},{"option":"VOTE_OPTION_NO","weight":"0.300000000000000000"},{"option":"VOTE_OPTION_ABSTAIN","weight":"0.050000000000000000"},{"option":"VOTE_OPTION_NO_WITH_VETO","weight":"0.050000000000000000"}],"metadata":""}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[{"public_key":{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"AymTGWvRvySz1XavUhrz6GIoN8ZyVo278i1TAUfrlJWZ"},"mode_info":{"single":{"mode":"SIGN_MODE_LEGACY_AMINO_JSON"}},"sequence":"2"}],"fee":{"amount":[],"gas_limit":"200000","payer":"","granter":""},"tip":null},"signatures":["4J2SPFz2RZFwlyw8iMoh/Oq2yX9Zjxg8ULK6Vo6DtDY5CNPEghRUbwLxP2jrPs9vDkn1Y17PrMhjOE9kBg8vpg=="]}`,
		},
		"MsgEditValidator - v050": {
			src: `{"body":{"messages":[{"@type":"/cosmos.staking.v1beta1.MsgEditValidator","description":{"moniker":"[do-not-modify]","identity":"[do-not-modify]","website":"[do-not-modify]","security_contact":"[do-not-modify]","details":"[do-not-modify]"},"validator_address":"cosmosvaloper1a40r995peauyekrkda2uky5h52j0r0ztvglyew","commission_rate":"0.250000000000000000","min_self_delegation":null}],"memo":"","timeout_height":"0","extension_options":[],"non_critical_extension_options":[]},"auth_info":{"signer_infos":[{"public_key":{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"AymTGWvRvySz1XavUhrz6GIoN8ZyVo278i1TAUfrlJWZ"},"mode_info":{"single":{"mode":"SIGN_MODE_LEGACY_AMINO_JSON"}},"sequence":"2"}],"fee":{"amount":[],"gas_limit":"200000","payer":"","granter":""},"tip":null},"signatures":["Jsj4mUDFElz3H9wgDQiF+MktyagWIOvSngn3zONVWuESXStsiUvMQsGCXKdTM9cjF0nuWJPydI2Z72YcF15hhg=="]}`,
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

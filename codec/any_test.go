package codec_test

import (
	"testing"

	"github.com/cosmos/cosmos-proto/anyutil"
	"github.com/cosmos/gogoproto/proto"
	"github.com/cosmos/gogoproto/types/any/test"
	"github.com/stretchr/testify/require"
	protov2 "google.golang.org/protobuf/proto"

	"github.com/cosmos/cosmos-sdk/codec"
	codectestutil "github.com/cosmos/cosmos-sdk/codec/testutil"
	"github.com/cosmos/cosmos-sdk/codec/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	"github.com/cosmos/cosmos-sdk/testutil/testdata/testpb"
	"github.com/cosmos/cosmos-sdk/types/module/testutil"
)

func NewTestInterfaceRegistry() codectypes.InterfaceRegistry {
	registry := codectypes.NewInterfaceRegistry()
	registry.RegisterInterface("Animal", (*test.Animal)(nil))
	registry.RegisterImplementations(
		(*test.Animal)(nil),
		&test.Dog{},
		&test.Cat{},
	)
	return registry
}

func TestMarshalAny(t *testing.T) {
	catRegistry := codectypes.NewInterfaceRegistry()
	catRegistry.RegisterImplementations((*test.Animal)(nil), &test.Cat{})

	registry := codectypes.NewInterfaceRegistry()

	cdc := codec.NewProtoCodec(registry)

	kitty := &test.Cat{Moniker: "Kitty"}
	emptyBz, err := cdc.MarshalInterface(kitty)
	require.ErrorContains(t, err, "does not have a registered interface")

	catBz, err := codec.NewProtoCodec(catRegistry).MarshalInterface(kitty)
	require.NoError(t, err)
	require.NotEmpty(t, catBz)

	var animal test.Animal

	// deserializing cat bytes should error in an empty registry
	err = cdc.UnmarshalInterface(catBz, &animal)
	require.ErrorContains(t, err, "no registered implementations of type test.Animal")

	// deserializing an empty byte array will return nil, but no error
	err = cdc.UnmarshalInterface(emptyBz, &animal)
	require.Nil(t, animal)
	require.NoError(t, err)

	// wrong type registration should fail
	registry.RegisterImplementations((*test.Animal)(nil), &test.Dog{})
	err = cdc.UnmarshalInterface(catBz, &animal)
	require.Error(t, err)

	// should pass
	registry = NewTestInterfaceRegistry()
	cdc = codec.NewProtoCodec(registry)
	err = cdc.UnmarshalInterface(catBz, &animal)
	require.NoError(t, err)
	require.Equal(t, kitty, animal)

	// nil should fail
	_ = NewTestInterfaceRegistry()
	err = cdc.UnmarshalInterface(catBz, nil)
	require.Error(t, err)
}

func TestMarshalProtoPubKey(t *testing.T) {
	require := require.New(t)
	ccfg := testutil.MakeTestEncodingConfig(codectestutil.CodecOptions{})
	privKey := ed25519.GenPrivKey()
	pk := privKey.PubKey()

	// **** test JSON serialization ****

	pkAny, err := codectypes.NewAnyWithValue(pk)
	require.NoError(err)
	bz, err := ccfg.Codec.MarshalJSON(pkAny)
	require.NoError(err)

	var pkAny2 codectypes.Any
	err = ccfg.Codec.UnmarshalJSON(bz, &pkAny2)
	require.NoError(err)
	// Before getting a cached value we need to unpack it.
	// Normally this happens in types which implement UnpackInterfaces
	var pkI cryptotypes.PubKey
	err = ccfg.InterfaceRegistry.UnpackAny(&pkAny2, &pkI)
	require.NoError(err)
	pk2 := pkAny2.GetCachedValue().(cryptotypes.PubKey)
	require.True(pk2.Equals(pk))

	// **** test binary serialization ****

	bz, err = ccfg.Codec.Marshal(pkAny)
	require.NoError(err)

	var pkAny3 codectypes.Any
	err = ccfg.Codec.Unmarshal(bz, &pkAny3)
	require.NoError(err)
	err = ccfg.InterfaceRegistry.UnpackAny(&pkAny3, &pkI)
	require.NoError(err)
	pk3 := pkAny3.GetCachedValue().(cryptotypes.PubKey)
	require.True(pk3.Equals(pk))
}

// TestMarshalProtoInterfacePubKey tests PubKey marshaling using (Un)marshalInterface
// helper functions
func TestMarshalProtoInterfacePubKey(t *testing.T) {
	require := require.New(t)
	ccfg := testutil.MakeTestEncodingConfig(codectestutil.CodecOptions{})
	privKey := ed25519.GenPrivKey()
	pk := privKey.PubKey()

	// **** test JSON serialization ****

	bz, err := ccfg.Codec.MarshalInterfaceJSON(pk)
	require.NoError(err)

	var pk3 cryptotypes.PubKey
	err = ccfg.Codec.UnmarshalInterfaceJSON(bz, &pk3)
	require.NoError(err)
	require.True(pk3.Equals(pk))

	// ** Check unmarshal using JSONCodec **
	// Unpacking won't work straightforward s Any type
	// Any can't implement UnpackInterfacesMessage interface. So Any is not
	// automatically unpacked and we won't get a value.
	var pkAny codectypes.Any
	err = ccfg.Codec.UnmarshalJSON(bz, &pkAny)
	require.NoError(err)
	ifc := pkAny.GetCachedValue()
	require.Nil(ifc)

	// **** test binary serialization ****

	bz, err = ccfg.Codec.MarshalInterface(pk)
	require.NoError(err)

	var pk2 cryptotypes.PubKey
	err = ccfg.Codec.UnmarshalInterface(bz, &pk2)
	require.NoError(err)
	require.True(pk2.Equals(pk))
}

func TestAminoGogoPulsarCompat(t *testing.T) {
	cat := &testdata.Cat{Lives: 10, Moniker: "bongo"}
	catAny, err := types.NewAnyWithValue(cat)
	require.NoError(t, err)
	msg := &testdata.HasAnimal{Animal: catAny, X: 23}
	msgBz, err := proto.Marshal(msg)
	require.NoError(t, err)
	t.Logf("len %d", len(msgBz))

	var msg2 testdata.HasAnimal
	err = proto.Unmarshal(msgBz, &msg2)
	require.NoError(t, err)
	t.Logf("%v", msg2.Animal.GetCachedValue())

	registry := codectypes.NewInterfaceRegistry()
	registry.RegisterInterface("Animal", (*testdata.Animal)(nil))
	registry.RegisterImplementations(
		(*testdata.Animal)(nil),
		&testdata.Dog{},
		&testdata.Cat{},
	)

	err = types.UnpackInterfaces(msg2, registry)
	require.NoError(t, err)
	t.Logf("%v", msg2.Animal.GetCachedValue())

	cdc := codec.NewProtoCodec(registry)
	err = cdc.Unmarshal(msgBz, &msg2)
	require.NoError(t, err)
	t.Logf("%v", msg2.Animal.GetCachedValue())
}

func Test_PulsarBackwardCompat(t *testing.T) {
	cat := &testpb.Cat{Lives: 10, Moniker: "pulsary"}
	catAny, err := anyutil.New(cat)
	require.NoError(t, err)
	msg := &testpb.HasAnimal{Animal: catAny, X: 23}
	msgBz, err := protov2.Marshal(msg)
	require.NoError(t, err)

	registry := codectypes.NewInterfaceRegistry()
	registry.RegisterInterface("Animal", (*testdata.Animal)(nil))
	registry.RegisterImplementations(
		(*testdata.Animal)(nil),
		&testdata.Cat{},
	)

	var msgRoundTrip testpb.HasAnimal
	err = protov2.Unmarshal(msgBz, &msgRoundTrip)
	require.NoError(t, err)

	var animal testdata.Animal
	err = registry.UnpackAny(
		&codectypes.Any{TypeUrl: msg.Animal.TypeUrl, Value: msg.Animal.Value},
		&animal,
	)
	require.NoError(t, err)
	t.Logf(animal.Greet())
}

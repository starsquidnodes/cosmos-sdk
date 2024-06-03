package aminojson

import (
	"fmt"
	"strings"
	"testing"

	gogoproto "github.com/cosmos/gogoproto/proto"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"

	"cosmossdk.io/x/tx/signing/aminojson"

	"github.com/cosmos/cosmos-sdk/codec"
	gogopb "github.com/cosmos/cosmos-sdk/tests/integration/tx/internal/gogo/testpb"
	pulsarpb "github.com/cosmos/cosmos-sdk/tests/integration/tx/internal/pulsar/testpb"
)

func TestRepeatedFields(t *testing.T) {
	cdc := codec.NewLegacyAmino()
	aj := aminojson.NewEncoder(aminojson.EncoderOptions{DoNotSortFields: true})

	cases := map[string]struct {
		gogo    gogoproto.Message
		pulsar  proto.Message
		unequal bool
		errs    bool
	}{
		"unsupported_empty_sets": {
			gogo:    &gogopb.TestRepeatedFields{},
			pulsar:  &pulsarpb.TestRepeatedFields{},
			unequal: true,
		},
		"unsupported_empty_sets_are_set": {
			gogo: &gogopb.TestRepeatedFields{
				NullableDontOmitempty: []*gogopb.Streng{{Value: "foo"}},
				NonNullableOmitempty:  []gogopb.Streng{{Value: "foo"}},
			},
			pulsar: &pulsarpb.TestRepeatedFields{
				NullableDontOmitempty: []*pulsarpb.Streng{{Value: "foo"}},
				NonNullableOmitempty:  []*pulsarpb.Streng{{Value: "foo"}},
			},
		},
		"unsupported_nullable": {
			gogo:   &gogopb.TestNullableFields{},
			pulsar: &pulsarpb.TestNullableFields{},
			errs:   true,
		},
		"unsupported_nullable_set": {
			gogo: &gogopb.TestNullableFields{
				NullableDontOmitempty:    &gogopb.Streng{Value: "foo"},
				NonNullableDontOmitempty: gogopb.Streng{Value: "foo"},
			},
			pulsar: &pulsarpb.TestNullableFields{
				NullableDontOmitempty:    &pulsarpb.Streng{Value: "foo"},
				NonNullableDontOmitempty: &pulsarpb.Streng{Value: "foo"},
			},
			unequal: true,
		},
	}

	for n, tc := range cases {
		t.Run(n, func(t *testing.T) {
			gogoBz, err := cdc.MarshalJSON(tc.gogo)
			require.NoError(t, err)
			pulsarBz, err := aj.Marshal(tc.pulsar)
			if tc.errs {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			fmt.Printf("  gogo: %s\npulsar: %s\n", string(gogoBz), string(pulsarBz))

			if tc.unequal {
				require.NotEqual(t, string(gogoBz), string(pulsarBz))
			} else {
				require.Equal(t, string(gogoBz), string(pulsarBz))
			}
		})
	}
}

func TestConfio_Message(t *testing.T) {
	aj := aminojson.NewEncoder(aminojson.EncoderOptions{DoNotSortFields: true})
	// test pulsar type
	msg := &pulsarpb.MsgStoreCode{
		Sender: "foo",
		InstantiatePermission: &pulsarpb.AccessConfig{
			Permission: pulsarpb.AccessType_ACCESS_TYPE_EVERYBODY,
		},
	}
	bz, err := aj.Marshal(msg)
	require.NotNil(t, bz)
	require.NoError(t, err)

	// test gogo type
	gogoMsg := &gogopb.MsgStoreCode{
		Sender: "foo",
		InstantiatePermission: &gogopb.AccessConfig{
			Permission: gogopb.AccessTypeEverybody,
		},
	}
	// convert the gogo type into an dynamicpb
	typeURL := strings.TrimPrefix(gogoproto.MessageName(gogoMsg), "/")
	msgDesc, err := gogoproto.GogoResolver.FindDescriptorByName(protoreflect.FullName(typeURL))
	require.NoError(t, err)
	dynamicMsg := dynamicpb.NewMessageType(msgDesc.(protoreflect.MessageDescriptor)).New().Interface()
	gogoProtoBz, err := gogoproto.Marshal(gogoMsg)
	require.NoError(t, err)
	// unmarshal into dynamic message
	err = proto.Unmarshal(gogoProtoBz, dynamicMsg)
	require.NoError(t, err)
	dynamicJSONbz, err := aj.Marshal(dynamicMsg)
	require.NoError(t, err)
	require.Equal(t, string(bz), string(dynamicJSONbz))
}

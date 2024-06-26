package tx

import (
	"bytes"
	txv1beta1 "cosmossdk.io/api/cosmos/tx/v1beta1"
	"cosmossdk.io/core/address"
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/x/tx/decode"
	"errors"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"google.golang.org/protobuf/encoding/protojson"
)

// DefaultJSONTxDecoder returns a default protobuf JSON TxDecoder using the provided Marshaler.
func DefaultJSONTxDecoder(addrCodec address.Codec, cdc codec.Codec, decoder *decode.Decoder) sdk.TxDecoder {
	jsonUnmarshaller := protojson.UnmarshalOptions{
		AllowPartial:   false,
		DiscardUnknown: false,
	}
	return func(txBytes []byte) (sdk.Tx, error) {
		jsonTx := new(txv1beta1.Tx)
		err := jsonUnmarshaller.Unmarshal(txBytes, jsonTx)
		if err != nil {
			return nil, err
		}

		// need to convert jsonTx into raw tx.
		bodyBytes, err := marshalOption.Marshal(jsonTx.Body)
		if err != nil {
			return nil, err
		}

		// decode old version
		var oldTx tx.Tx
		err = cdc.UnmarshalJSON(txBytes, &oldTx)
		if err != nil {
			return nil, errorsmod.Wrap(sdkerrors.ErrTxDecode, err.Error())
		}

		if !bytes.Equal(oldTx.Body.Messages[0].Value, jsonTx.Body.Messages[0].Value) {
			return nil, errors.New("no, here we go")
		}

		authInfoBytes, err := marshalOption.Marshal(jsonTx.AuthInfo)
		if err != nil {
			return nil, err
		}

		protoTxBytes, err := marshalOption.Marshal(&txv1beta1.TxRaw{
			BodyBytes:     bodyBytes,
			AuthInfoBytes: authInfoBytes,
			Signatures:    jsonTx.Signatures,
		})
		if err != nil {
			return nil, err
		}

		decodedTx, err := decoder.Decode(protoTxBytes)
		if err != nil {
			return nil, err
		}
		return newWrapperFromDecodedTx(addrCodec, cdc, decodedTx)
	}
}

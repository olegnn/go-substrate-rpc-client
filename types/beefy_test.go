package types_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/olegnn/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/assert"
)

var sig1 = [65]byte{85, 132, 85, 173, 129, 39, 157, 240, 121, 92, 201, 133, 88, 14, 79, 183, 93, 114, 217, 72, 209, 16, 123, 42, 200, 10, 9, 171, 237, 77, 168, 72, 12, 116, 108, 195, 33, 242, 49, 154, 94, 153, 168, 48, 227, 20, 209, 13, 211, 205, 104, 206, 61, 192, 195, 60, 134, 233, 155, 203, 120, 22, 249, 186, 1}
var sig2 = [65]byte{45, 110, 31, 129, 5, 195, 55, 168, 108, 221, 154, 170, 205, 196, 150, 87, 127, 61, 184, 197, 94, 249, 230, 253, 72, 242, 197, 192, 90, 34, 116, 112, 116, 145, 99, 93, 139, 163, 223, 100, 243, 36, 87, 91, 123, 42, 52, 72, 123, 202, 35, 36, 182, 160, 4, 99, 149, 167, 22, 129, 190, 61, 12, 42, 0}

func TestBeefySignature(t *testing.T) {
	empty := types.NewOptionBeefySignatureEmpty()
	assert.True(t, empty.IsNone())
	assert.False(t, empty.IsSome())

	sig := types.NewOptionBeefySignature(types.BeefySignature{})
	sig.SetNone()
	assert.True(t, sig.IsNone())
	sig.SetSome(types.BeefySignature{})
	assert.True(t, sig.IsSome())
	ok, _ := sig.Unwrap()
	assert.True(t, ok)
	assertRoundtrip(t, sig)
}

func makeCommitment() (*types.Commitment, error) {
	data, err := types.EncodeToBytes([]byte("Hello World!"))
	if err != nil {
		return nil, err
	}

	payloadItem := types.PayloadItem{
		ID:   [2]byte{'m', 'h'},
		Data: data,
	}

	commitment := types.Commitment{
		Payload:        []types.PayloadItem{payloadItem},
		BlockNumber:    5,
		ValidatorSetID: 0,
	}

	return &commitment, nil
}

func makeLargeCommitment() (*types.Commitment, error) {
	data := types.MustHexDecodeString("0xb5bb9d8014a0f9b1d61e21e796d78dccdf1352f23cd32812f4850b878ae4944c")

	payloadItem := types.PayloadItem{
		ID:   [2]byte{'m', 'h'},
		Data: data,
	}

	commitment := types.Commitment{
		Payload:        []types.PayloadItem{payloadItem},
		BlockNumber:    5,
		ValidatorSetID: 3,
	}

	return &commitment, nil
}

func TestCommitment_Encode(t *testing.T) {
	c, err := makeCommitment()
	assert.NoError(t, err)
	assertEncode(t, []encodingAssert{
		{c, types.MustHexDecodeString("0x046d68343048656c6c6f20576f726c6421050000000000000000000000")},
	})
}

func TestLargeCommitment_Encode(t *testing.T) {
	c, err := makeLargeCommitment()
	assert.NoError(t, err)
	fmt.Println(len(c.Payload[0].Data))
	fmt.Println(types.EncodeToHexString(c))
}

func TestCommitment_Decode(t *testing.T) {
	c, err := makeCommitment()
	assert.NoError(t, err)

	assertDecode(t, []decodingAssert{
		{
			input:    types.MustHexDecodeString("0x046d68343048656c6c6f20576f726c6421050000000000000000000000"),
			expected: *c,
		},
	})
}

func TestCommitment_EncodeDecode(t *testing.T) {
	c, err := makeCommitment()
	assert.NoError(t, err)

	assertRoundtrip(t, *c)
}

func TestSignedCommitment_Decode(t *testing.T) {
	c, err := makeCommitment()
	assert.NoError(t, err)

	s := types.SignedCommitment{
		Commitment: *c,
		Signatures: []types.OptionBeefySignature{
			types.NewOptionBeefySignatureEmpty(),
			types.NewOptionBeefySignatureEmpty(),
			types.NewOptionBeefySignature(sig1),
			types.NewOptionBeefySignature(sig2),
		},
	}

	assertDecode(t, []decodingAssert{
		{
			input:    types.MustHexDecodeString("0x046d68343048656c6c6f20576f726c642105000000000000000000000004300400000008558455ad81279df0795cc985580e4fb75d72d948d1107b2ac80a09abed4da8480c746cc321f2319a5e99a830e314d10dd3cd68ce3dc0c33c86e99bcb7816f9ba012d6e1f8105c337a86cdd9aaacdc496577f3db8c55ef9e6fd48f2c5c05a2274707491635d8ba3df64f324575b7b2a34487bca2324b6a0046395a71681be3d0c2a00"),
			expected: s,
		},
	})
}

func TestSignedCommitment_EncodeDecode(t *testing.T) {
	c, err := makeCommitment()
	assert.NoError(t, err)

	s := types.SignedCommitment{
		Commitment: *c,
		Signatures: []types.OptionBeefySignature{
			types.NewOptionBeefySignatureEmpty(),
			types.NewOptionBeefySignatureEmpty(),
			types.NewOptionBeefySignature(sig1),
			types.NewOptionBeefySignature(sig1),
			types.NewOptionBeefySignatureEmpty(),
			types.NewOptionBeefySignatureEmpty(),
			types.NewOptionBeefySignatureEmpty(),
			types.NewOptionBeefySignatureEmpty(),
			types.NewOptionBeefySignatureEmpty(),
			types.NewOptionBeefySignature(sig1),
		},
	}

	assertRoundtrip(t, s)
}

func TestOptionBeefySignature_Marshal(t *testing.T) {
	actual, err := json.Marshal(types.NewOptionBeefySignature(sig1))
	assert.NoError(t, err)
	expected, err := json.Marshal(sig1)
	assert.NoError(t, err)
	assert.Equal(t, actual, expected)

	actual, err = json.Marshal(types.NewOptionBeefySignatureEmpty())
	assert.NoError(t, err)
	expected, err = json.Marshal(nil)
	assert.NoError(t, err)
	assert.Equal(t, actual, expected)
}

func TestOptionBeefySignature_MarshalUnmarshal(t *testing.T) {
	expected := types.NewOptionBeefySignature(sig1)

	marshalled, err := json.Marshal(expected)
	assert.NoError(t, err)

	var unmarshalled types.OptionBeefySignature
	err = json.Unmarshal(marshalled, &unmarshalled)
	assert.NoError(t, err)

	assert.Equal(t, expected, unmarshalled)
}

func TestOptionBeefySignature_MarshalUnmarshalEmpty(t *testing.T) {
	expected := types.NewOptionBeefySignatureEmpty()

	marshalled, err := json.Marshal(expected)
	assert.NoError(t, err)

	var unmarshalled types.OptionBeefySignature
	err = json.Unmarshal(marshalled, &unmarshalled)
	assert.NoError(t, err)

	assert.Equal(t, expected, unmarshalled)
}

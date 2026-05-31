package prsr

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFrame(t *testing.T) {
	t.Run("Protocol Header", func(t *testing.T) {
		data, err := hex.DecodeString("414d515003010000")
		assert.Nil(t, err)

		result := ParseFrame(data)
		dumpResult(result)

		assert.NotNil(t, result)
		assert.Equal(t, 0, len(result.RemainingBytes))

		expectedResult := []any{
			map[string]any{
				"_type":       "PROTOCOL_HEADER",
				"protocol_id": uint8(3),
				"major":       uint8(1),
				"minor":       uint8(0),
				"revision":    uint8(0),
			},
		}
		assert.EqualValues(t, expectedResult, result.Results)
	})

	t.Run("SASL INIT", func(t *testing.T) {
		data, err := hex.DecodeString("0000004602010000005341d00000003600000003a309414e4f4e594d4f5553a01a00526f6f744d616e6167655368617265644163636573734b6579a1093132372e302e302e31")
		assert.Nil(t, err)

		result := ParseFrame(data)
		dumpResult(result)

		assert.NotNil(t, result)
		assert.Equal(t, 0, len(result.RemainingBytes))

		expectedResult := []any{
			map[string]any{
				"_type":            DescriptorKey("SASL_INIT"),
				"mechanism":        "ANONYMOUS",
				"initial-response": []uint8{0x0, 0x52, 0x6f, 0x6f, 0x74, 0x4d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x53, 0x68, 0x61, 0x72, 0x65, 0x64, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x4b, 0x65, 0x79},
				"hostname":         "127.0.0.1",
			},
		}
		assert.EqualValues(t, expectedResult, result.Results)
	})

	t.Run("BEGIN", func(t *testing.T) {
		data, err := hex.DecodeString("0000002002000000005311d000000010000000044043700000080070ffffffff")
		assert.Nil(t, err)

		result := ParseFrame(data)
		dumpResult(result)

		assert.NotNil(t, result)
		assert.Equal(t, 0, len(result.RemainingBytes))

		expectedResult := []any{
			map[string]any{
				"_type":            DescriptorKey("BEGIN"),
				"incoming-window":  uint32(2048),
				"next-outgoing-id": uint(0),
				"outgoing-window":  uint32(4294967295),
				"remote-channel":   nil,
			},
		}
		assert.EqualValues(t, expectedResult, result.Results)
	})

	t.Run("ATTACH-1", func(t *testing.T) {
		data, err := hex.DecodeString("0000005702000000005312d0000000470000000aa12433333735373332332d373664652d333434342d613863652d3137663561313137363732344342404000532845005329d00000000a00000001a10424636273404043")
		assert.Nil(t, err)

		result := ParseFrame(data)
		dumpResult(result)

		assert.NotNil(t, result)
		assert.Equal(t, 0, len(result.RemainingBytes))

		expectedResult := []any{
			map[string]any{
				"_type":                  DescriptorKey("ATTACH"),
				"handle":                 uint(0),
				"incomplete-unsettled":   nil,
				"initial-delivery-count": uint(0),
				"name":                   "33757323-76de-3444-a8ce-17f5a1176724",
				"rcv-settle-mode":        nil,
				"role":                   bool(false),
				"snd-settle-mode":        nil,
				"source":                 map[string]any{"_type": DescriptorKey("SOURCE")},
				"target":                 map[string]any{"_type": DescriptorKey("TARGET"), "address": "$cbs"},
				"unsettled":              nil,
			},
		}
		assert.EqualValues(t, expectedResult, result.Results)
	})

	t.Run("ATTACH-2", func(t *testing.T) {
		data, err := hex.DecodeString("0000009802000001005312d0000000880000000aa12c71756575652e312d31393236323965622d326233342d653434332d393939312d39626262643262373366633443424040005328d00000003200000001a12c71756575652e312d39663539616532322d323865352d663534382d383438322d316138336133633230363262005329d00000000d00000001a10771756575652e31404043")
		assert.Nil(t, err)

		result := ParseFrame(data)
		dumpResult(result)

		assert.NotNil(t, result)
		assert.Equal(t, 0, len(result.RemainingBytes))

		expectedResult := []any{
			map[string]any{
				"_type":                  DescriptorKey("ATTACH"),
				"handle":                 uint(0),
				"incomplete-unsettled":   nil,
				"initial-delivery-count": uint(0),
				"name":                   "queue.1-192629eb-2b34-e443-9991-9bbbd2b73fc4",
				"rcv-settle-mode":        nil,
				"role":                   false,
				"snd-settle-mode":        nil,
				"source":                 map[string]any{"_type": DescriptorKey("SOURCE"), "address": "queue.1-9f59ae22-28e5-f548-8482-1a83a3c2062b"},
				"target":                 map[string]any{"_type": DescriptorKey("TARGET"), "address": "queue.1"},
				"unsettled":              nil,
			},
		}
		assert.EqualValues(t, expectedResult, result.Results)
	})

	t.Run("TRANSFER (Queue message)", func(t *testing.T) {
		data, err := hex.DecodeString("000000b002000001005314d00000000c000000064343a0013043424200537045005372d10000000400000000005373d00000001800000001a112746573742d31373731353231313734343932005374d10000003f00000004a1087465737454797065a10a736d6f6b652d74657374a10974696d657374616d70a118323032362d30322d31395431373a31323a35342e3439325a005375a0182248656c6c6f2066726f6d20736d6f6b6520746573742122")
		assert.Nil(t, err)

		result := ParseFrame(data)
		dumpResult(result)

		assert.NotNil(t, result)
		assert.Equal(t, 0, len(result.RemainingBytes))

		expectedResult := []any{
			map[string]any{
				"_type":          DescriptorKey("TRANSFER"),
				"delivery-id":    uint(0),
				"delivery-tag":   []uint8{0x30},
				"handle":         uint(0),
				"message-format": uint(0),
				"more":           false,
				"settled":        false,
			},
			map[string]any{"_type": DescriptorKey("HEADER")},
			map[string]any{"_type": DescriptorKey("MESSAGE_ANNOTATIONS"), "_fields": map[string]any{}},
			map[string]any{"_type": DescriptorKey("MESSAGE_PROPERTIES"), "message-id": "test-1771521174492"},
			map[string]any{
				"_type": DescriptorKey("APPLICATION_PROPERTIES"),
				"_fields": map[string]any{
					"testType":  "smoke-test",
					"timestamp": "2026-02-19T17:12:54.492Z",
				},
			},
			map[string]any{
				"_type":   DescriptorKey("APPLICATION_DATA"),
				"_fields": []uint8{0x22, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x66, 0x72, 0x6f, 0x6d, 0x20, 0x73, 0x6d, 0x6f, 0x6b, 0x65, 0x20, 0x74, 0x65, 0x73, 0x74, 0x21, 0x22},
			},
		}
		assert.EqualValues(t, expectedResult, result.Results)
	})
}

func dumpResult(result ParseFrameResult) {
	json, err := json.Marshal(result.Results)
	check(err, "Failed to marshal result")
	fmt.Printf("--- Result: %s\n", json)
}

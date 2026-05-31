package prsr

import (
	"bufio"
	"fmt"
)

type DescriptorKey string

const (
	DescriptorKeyOPEN        DescriptorKey = "OPEN"
	DescriptorKeyBEGIN       DescriptorKey = "BEGIN"
	DescriptorKeyATTACH      DescriptorKey = "ATTACH"
	DescriptorKeyFLOW        DescriptorKey = "FLOW"
	DescriptorKeyTRANSFER    DescriptorKey = "TRANSFER"
	DescriptorKeyDISPOSITION DescriptorKey = "DISPOSITION"
	DescriptorKeyDETACH      DescriptorKey = "DETACH"
	DescriptorKeyEND         DescriptorKey = "END"
	DescriptorKeyCLOSE       DescriptorKey = "CLOSE"

	DescriptorKeyRECEIVED            DescriptorKey = "RECEIVED"
	DescriptorKeyACCEPTED            DescriptorKey = "ACCEPTED"
	DescriptorKeyREJECTED            DescriptorKey = "REJECTED"
	DescriptorKeyRELEASED            DescriptorKey = "RELEASED"
	DescriptorKeyMODIFIED            DescriptorKey = "MODIFIED"
	DescriptorKeySOURCE              DescriptorKey = "SOURCE"
	DescriptorKeyTARGET              DescriptorKey = "TARGET"
	DescriptorKeyCOORDINATOR         DescriptorKey = "COORDINATOR"
	DescriptorKeyDECLARE             DescriptorKey = "DECLARE"
	DescriptorKeyDISCHARGE           DescriptorKey = "DISCHARGE"
	DescriptorKeyDECLARED            DescriptorKey = "DECLARED"
	DescriptorKeyTRANSACTIONAL_STATE DescriptorKey = "TRANSACTIONAL_STATE"

	DescriptorKeySASL_MECHANISMS DescriptorKey = "SASL_MECHANISMS"
	DescriptorKeySASL_INIT       DescriptorKey = "SASL_INIT"
	DescriptorKeySASL_CHALLENGE  DescriptorKey = "SASL_CHALLENGE"
	DescriptorKeySASL_RESPONSE   DescriptorKey = "SASL_RESPONSE"
	DescriptorKeySASL_OUTCOME    DescriptorKey = "SASL_OUTCOME"

	DescriptorKeyHEADER                 DescriptorKey = "HEADER"
	DescriptorKeyDELIVERY_ANNOTATIONS   DescriptorKey = "DELIVERY_ANNOTATIONS"
	DescriptorKeyMESSAGE_ANNOTATIONS    DescriptorKey = "MESSAGE_ANNOTATIONS"
	DescriptorKeyMESSAGE_PROPERTIES     DescriptorKey = "MESSAGE_PROPERTIES"
	DescriptorKeyAPPLICATION_PROPERTIES DescriptorKey = "APPLICATION_PROPERTIES"
	DescriptorKeyAPPLICATION_DATA       DescriptorKey = "APPLICATION_DATA"
	DescriptorKeyAMQP_SEQUENCE          DescriptorKey = "AMQP_SEQUENCE"
	DescriptorKeyAMQP_VALUE             DescriptorKey = "AMQP_VALUE"
	DescriptorKeyFOOTER                 DescriptorKey = "FOOTER"

	DescriptorKeyUNKNOWN DescriptorKey = "UNKNOWN"
)

var descriptorByUint8 = map[uint8]DescriptorKey{
	0x10: DescriptorKeyOPEN,
	0x11: DescriptorKeyBEGIN,
	0x12: DescriptorKeyATTACH,
	0x13: DescriptorKeyFLOW,
	0x14: DescriptorKeyTRANSFER,
	0x15: DescriptorKeyDISPOSITION,
	0x16: DescriptorKeyDETACH,
	0x17: DescriptorKeyEND,
	0x18: DescriptorKeyCLOSE,

	0x23: DescriptorKeyRECEIVED,
	0x24: DescriptorKeyACCEPTED,
	0x25: DescriptorKeyREJECTED,
	0x26: DescriptorKeyRELEASED,
	0x27: DescriptorKeyMODIFIED,
	0x28: DescriptorKeySOURCE,
	0x29: DescriptorKeyTARGET,
	0x30: DescriptorKeyCOORDINATOR,
	0x31: DescriptorKeyDECLARE,
	0x32: DescriptorKeyDISCHARGE,
	0x33: DescriptorKeyDECLARED,
	0x34: DescriptorKeyTRANSACTIONAL_STATE,

	0x40: DescriptorKeySASL_MECHANISMS,
	0x41: DescriptorKeySASL_INIT,
	0x42: DescriptorKeySASL_CHALLENGE,
	0x43: DescriptorKeySASL_RESPONSE,
	0x44: DescriptorKeySASL_OUTCOME,

	0x70: DescriptorKeyHEADER,
	0x71: DescriptorKeyDELIVERY_ANNOTATIONS,
	0x72: DescriptorKeyMESSAGE_ANNOTATIONS,
	0x73: DescriptorKeyMESSAGE_PROPERTIES,
	0x74: DescriptorKeyAPPLICATION_PROPERTIES,
	0x75: DescriptorKeyAPPLICATION_DATA,
	0x76: DescriptorKeyAMQP_SEQUENCE,
	0x77: DescriptorKeyAMQP_VALUE,
	0x78: DescriptorKeyFOOTER,
}

func descriptorToKey(descriptor any) DescriptorKey {
	switch descriptor := descriptor.(type) {
	case uint8:
		dk, ok := descriptorByUint8[descriptor]
		if !ok {
			fmt.Printf("Unknown descriptor: %v\n", descriptor)
			return DescriptorKeyUNKNOWN
		}
		return dk
	}
	fmt.Printf("Unknown descriptor: %v\n", descriptor)
	return DescriptorKeyUNKNOWN
}

var fieldsByDescriptor = map[DescriptorKey][]string{
	DescriptorKeyOPEN: {
		"container-id",
		"hostname",
		"max-frame-size",
		"channel-max",
		"idle-time-out",
		"outgoing-locales",
		"incoming-locales",
		"offered-capabilities",
		"desired-capabilities",
		"properties",
	},
	DescriptorKeyBEGIN: {
		"remote-channel",
		"next-outgoing-id",
		"incoming-window",
		"outgoing-window",
		"handle-max",
		"offered-capabilities",
		"desired-capabilities",
		"properties",
	},
	DescriptorKeyATTACH: {
		"name",
		"handle",
		"role",
		"snd-settle-mode",
		"rcv-settle-mode",
		"source",
		"target",
		"unsettled",
		"incomplete-unsettled",
		"initial-delivery-count",
		"max-message-size",
		"offered-capabilities",
		"desired-capabilities",
		"properties",
	},
	DescriptorKeyFLOW: {
		"next-incoming-id",
		"incoming-window",
		"next-outgoing-id",
		"outgoing-window",
		"handle",
		"delivery-count",
		"link-credit",
		"available",
		"drain",
		"echo",
		"properties",
	},
	DescriptorKeyTRANSFER: {
		"handle",
		"delivery-id",
		"delivery-tag",
		"message-format",
		"settled",
		"more",
		"rcv-settle-mode",
		"state",
		"resume",
		"aborted",
		"batchable",
	},
	DescriptorKeyDISPOSITION: {
		"role",
		"first",
		"last",
		"settled",
		"state",
		"batchable",
	},
	DescriptorKeyDETACH: {"handle", "closed", "error"},
	DescriptorKeyEND:    {"error"},
	DescriptorKeyCLOSE:  {"error"},

	DescriptorKeyRECEIVED: {"section-number", "section-offset"},
	DescriptorKeyREJECTED: {"error"},
	DescriptorKeyMODIFIED: {"delivery-failed", "undeliverable-here", "message-annotations"},
	DescriptorKeySOURCE: {
		"address",
		"durable",
		"expiry-policy",
		"timeout",
		"dynamic",
		"dynamic-node-properties",
		"distribution-mode",
		"filter",
		"default-outcome",
		"outcomes",
		"capabilities",
	},
	DescriptorKeyTARGET: {
		"address",
		"durable",
		"expiry-policy",
		"timeout",
		"dynamic",
		"dynamic-node-properties",
		"capabilities",
	},
	DescriptorKeyCOORDINATOR:         {"capabilities"},
	DescriptorKeyDECLARE:             {"global-id"},
	DescriptorKeyDISCHARGE:           {"txn-id", "fail"},
	DescriptorKeyDECLARED:            {"txn-id"},
	DescriptorKeyTRANSACTIONAL_STATE: {"txn-id", "outcome"},

	DescriptorKeySASL_MECHANISMS: {"mechanisms"},
	DescriptorKeySASL_INIT:       {"mechanism", "initial-response", "hostname"},
	DescriptorKeySASL_CHALLENGE:  {"challenge"},
	DescriptorKeySASL_RESPONSE:   {"response"},
	DescriptorKeySASL_OUTCOME:    {"code", "additional-data"},
	DescriptorKeyHEADER:          {"durable", "priority", "ttl", "first-acquirer", "delivery-count"},
	DescriptorKeyMESSAGE_PROPERTIES: {
		"message-id",
		"user-id",
		"to",
		"subject",
		"reply-to",
		"correlation-id",
		"content-type",
		"content-encoding",
		"absolute-expiry-time",
		"creation-time",
		"group-id",
		"group-sequence",
		"reply-to-group-id",
	},
}

func parseDescribedPayload(reader *bufio.Reader) any {
	descriptor := descriptorToKey(parseNext(reader))
	fields := parseNext(reader)

	fieldArr, isFieldArray := fields.([]any)
	fieldKeys, ok := fieldsByDescriptor[descriptor]
	if !ok || !isFieldArray {
		return map[string]any{
			"_type":   descriptor,
			"_fields": fields,
		}
	}

	result := map[string]any{
		"_type": descriptor,
	}

	for i, field := range fieldArr {
		if i >= len(fieldKeys) {
			result["_unknown_fields"] = fieldArr[i:]
			break
		}

		result[fieldKeys[i]] = field
	}

	return result
}

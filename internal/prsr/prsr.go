package prsr

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

type ParseFrameResult struct {
	Success        bool
	Channel        int16
	ExtendedHeader []byte
	Results        []any
	RemainingBytes []byte
	FrameType      string
}

var PROTOCOL_HEADER_START = []byte("AMQP")

func ParseFrame(data []byte) ParseFrameResult {
	if len(data) < 8 {
		return ParseFrameResult{
			Success: false,
		}
	}

	if bytes.Equal(data[:len(PROTOCOL_HEADER_START)], PROTOCOL_HEADER_START) {
		return ParseFrameResult{
			Success: true,
			Results: []any{
				map[string]any{
					"_type":       "PROTOCOL_HEADER",
					"protocol_id": data[len(PROTOCOL_HEADER_START)],
					"major":       data[len(PROTOCOL_HEADER_START)+1],
					"minor":       data[len(PROTOCOL_HEADER_START)+2],
					"revision":    data[len(PROTOCOL_HEADER_START)+3],
				},
			},
			RemainingBytes: data[len(PROTOCOL_HEADER_START)+4:],
		}
	}

	reader := bufio.NewReader(bytes.NewReader(data))

	sizeBytes := readBytes(reader, 4)
	size := binary.BigEndian.Uint32(sizeBytes)

	dataOffsetBytes := readBytes(reader, 1)
	dataOffset := uint8(dataOffsetBytes[0])

	frameTypeBytes := readBytes(reader, 1)
	frameType := frameTypeBytes[0]

	result := ParseFrameResult{}
	switch frameType {
	case 0x00:
		// AMQP frame
		result.FrameType = "AMQP"
		result.Channel = int16(binary.BigEndian.Uint16(readBytes(reader, 2)))
	case 0x01:
		// SASL frame
		result.FrameType = "SASL"
		readBytes(reader, 2)
	default:
		result.Success = false
		return result
	}

	// Read the extended header
	extendedHeaderLength := dataOffset*4 - 8
	extendedHeader := readBytes(reader, int(extendedHeaderLength))
	result.ExtendedHeader = extendedHeader

	// Read payload
	payloadSize := size - uint32(dataOffset)*4
	payload := readBytes(reader, int(payloadSize))

	results, remainingBytes := ParseFrameBody(payload)
	result.Results = results
	result.RemainingBytes = remainingBytes
	result.Success = true

	return result
}

func ParseFrameBody(data []byte) ([]any, []byte) {
	reader := bufio.NewReader(bytes.NewReader(data))

	results := make([]any, 0)
	for {
		_, err := reader.Peek(1)
		if err != nil {
			break
		}

		result := parseNext(reader)
		results = append(results, result)
	}

	remaining, _ := io.ReadAll(reader)

	return results, remaining
}

func parseNext(reader *bufio.Reader) any {
	dataType, err := reader.ReadByte()
	check(err, "Failed to read data type")

	return parseByDataType(dataType, reader)
}

func parseByDataType(dataType byte, reader *bufio.Reader) any {
	switch dataType {
	case 0x00: // described
		return parseDescribedPayload(reader)

	case 0xc0: // list8
		sizeByte, err := reader.ReadByte()
		check(err, "Failed to read size byte")
		size := uint8(sizeByte)

		list8payload := readBytes(reader, int(size))
		list8payloadReader := bufio.NewReader(bytes.NewReader(list8payload))

		entryCountBytes := readBytes(list8payloadReader, 1)
		entryCount := uint8(entryCountBytes[0])

		entries := make([]any, entryCount)
		for i := 0; i < int(entryCount) && hasMore(list8payloadReader); i++ {
			entries[i] = parseNext(list8payloadReader)
		}
		return entries
	case 0xd0: // list32
		sizeBytes := readBytes(reader, 4)
		size := binary.BigEndian.Uint32(sizeBytes)

		list32payload := readBytes(reader, int(size))
		list32payloadReader := bufio.NewReader(bytes.NewReader(list32payload))

		entryCountBytes := readBytes(list32payloadReader, 4)
		entryCount := binary.BigEndian.Uint32(entryCountBytes)

		entries := make([]any, entryCount)
		for i := 0; i < int(entryCount) && hasMore(list32payloadReader); i++ {
			entries[i] = parseNext(list32payloadReader)
		}
		return entries

	case 0xc1: // map8
		sizeByte, err := reader.ReadByte()
		check(err, "Failed to read size byte")
		size := uint8(sizeByte)

		map8payload := readBytes(reader, int(size))
		map8payloadReader := bufio.NewReader(bytes.NewReader(map8payload))

		entryCountBytes := readBytes(map8payloadReader, 1)
		entryCount := uint8(entryCountBytes[0])

		entries := make(map[string]any)
		for i := 0; i < int(entryCount) && hasMore(map8payloadReader); i++ {
			key := parseNext(map8payloadReader)
			value := parseNext(map8payloadReader)
			if key, ok := key.(string); ok {
				entries[key] = value
			} else {
				panic("Invalid key: " + fmt.Sprintf("%v", key))
			}
		}
		return entries
	case 0xd1: // map32
		readBytes(reader, 4)

		entryCountBytes := readBytes(reader, 4)
		entryCount := binary.BigEndian.Uint32(entryCountBytes)

		if entryCount%2 != 0 {
			fmt.Printf("Invalid entry count: %d\n", entryCount)
			os.Exit(1)
		}

		entries := make(map[string]any)
		for i := 0; i < int(entryCount); i += 2 {
			key := parseNext(reader)
			value := parseNext(reader)
			if key, ok := key.(string); ok {
				entries[key] = value
			} else {
				panic("Invalid key: " + fmt.Sprintf("%v", key))
			}
		}
		return entries

	case 0xe0: // array8
		sizeByte, err := reader.ReadByte()
		check(err, "Failed to read size byte")
		size := uint8(sizeByte)

		array8payload := readBytes(reader, int(size))
		array8payloadReader := bufio.NewReader(bytes.NewReader(array8payload))

		entryCountBytes := readBytes(array8payloadReader, 1)
		entryCount := uint8(entryCountBytes[0])

		entryTypeByte, err := array8payloadReader.ReadByte()
		check(err, "Failed to read entry type byte")

		entries := make([]any, entryCount)
		for i := 0; i < int(entryCount) && hasMore(array8payloadReader); i++ {
			entries[i] = parseByDataType(entryTypeByte, array8payloadReader)
		}
		return entries
	case 0xf0: // array32
		sizeBytes := readBytes(reader, 4)
		size := binary.BigEndian.Uint32(sizeBytes)

		array32payload := readBytes(reader, int(size))
		array32payloadReader := bufio.NewReader(bytes.NewReader(array32payload))

		entryCountBytes := readBytes(array32payloadReader, 4)
		entryCount := binary.BigEndian.Uint32(entryCountBytes)

		entryTypeByte, err := array32payloadReader.ReadByte()
		check(err, "Failed to read entry type byte")

		entries := make([]any, entryCount)
		for i := 0; i < int(entryCount) && hasMore(array32payloadReader); i++ {
			entries[i] = parseByDataType(entryTypeByte, array32payloadReader)
		}
		return entries

	case 0x40: // nil
		return nil
	case 0x41: // bool = true
		return true
	case 0x42: // bool = false
		return false
	case 0x43: // uint = 0
		return uint(0)
	case 0x44: // ulong = 0
		return uint64(0)
	case 0x45: // empty list
		return []any{}
	case 0x50: // ubyte
		return uint(readBytes(reader, 1)[0])
	case 0x51: // byte
		return readBytes(reader, 1)[0]
	case 0x52: // uint8
		return uint(readBytes(reader, 1)[0])
	case 0x53: // small ulong
		return uint8(readBytes(reader, 1)[0])
	case 0x54: // small int
		return int8(readBytes(reader, 1)[0])
	case 0x55: // small long
		return int8(readBytes(reader, 1)[0])
	case 0x56: // bool
		return readBytes(reader, 1)[0] == 0x01
	case 0x60: // uint16
		return binary.BigEndian.Uint16(readBytes(reader, 2))
	case 0x61: // short
		return int16(binary.BigEndian.Uint16(readBytes(reader, 2)))
	case 0x70: // uint32
		return binary.BigEndian.Uint32(readBytes(reader, 4))
	case 0x71: // int32
		return int32(binary.BigEndian.Uint32(readBytes(reader, 4)))
	case 0x72: // float
		// I can't be bothered to implement this, let's just return a byte array
		return readBytes(reader, 4)
	case 0x73: // utf32 char
		// I can't be bothered to implement this, let's just return a byte array
		return readBytes(reader, 4)
	case 0x74: // decimal32
		// I can't be bothered to implement this, let's just return a byte array
		return readBytes(reader, 4)
	case 0x80: // uint64
		return binary.BigEndian.Uint64(readBytes(reader, 8))
	case 0x81: // long
		return int64(binary.BigEndian.Uint64(readBytes(reader, 8)))
	case 0x82: // double
		// I can't be bothered to implement this, let's just return a byte array
		return readBytes(reader, 8)
	case 0x83: // timestamp
		return int64(binary.BigEndian.Uint64(readBytes(reader, 8)))
	case 0x84: // decimal64
		// I can't be bothered to implement this, let's just return a byte array
		return readBytes(reader, 8)
	case 0x94: // decimal128
		// I can't be bothered to implement this, let's just return a byte array
		return readBytes(reader, 16)
	case 0x98: // UUID
		// I'm unbothered to implement this, let's just return a byte array
		return readBytes(reader, 16)

	case 0xa0: // vbin8 (binary)
		sizeByte, err := reader.ReadByte()
		check(err, "Failed to read size byte")
		size := uint8(sizeByte)
		data := readBytes(reader, int(size))
		return data
	case 0xb0: // vbin32 (binary)
		sizeBytes := readBytes(reader, 4)
		size := binary.BigEndian.Uint32(sizeBytes)
		data := readBytes(reader, int(size))
		return data

	case 0xa1: // str8-utf8
		sizeByte, err := reader.ReadByte()
		check(err, "Failed to read size byte")
		size := uint8(sizeByte)
		data := readBytes(reader, int(size))
		return string(data)
	case 0xb1: // str32-utf8
		sizeBytes := readBytes(reader, 4)
		size := binary.BigEndian.Uint32(sizeBytes)
		data := readBytes(reader, int(size))
		return string(data)

	case 0xa3: // symbol8
		sizeByte, err := reader.ReadByte()
		check(err, "Failed to read size byte")
		size := uint8(sizeByte)
		data := readBytes(reader, int(size))
		return string(data)
	case 0xb3: // symbol32
		sizeBytes := readBytes(reader, 4)
		size := binary.BigEndian.Uint32(sizeBytes)
		data := readBytes(reader, int(size))
		return string(data)
	}

	return nil
}

func check(err error, msg string) {
	if err != nil {
		fmt.Printf("%s: %v", msg, err)
		os.Exit(1)
	}
}

func readBytes(reader *bufio.Reader, n int) []byte {
	bytes := make([]byte, n)
	_, err := io.ReadFull(reader, bytes)
	if err != nil {
		fmt.Printf("Error reading size: %v\n", err)
		os.Exit(0)
	}

	return bytes
}

func hasMore(reader *bufio.Reader) bool {
	_, err := reader.Peek(1)
	return err == nil
}

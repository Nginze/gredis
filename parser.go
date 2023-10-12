// deserialize RESP
package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Value struct {
	typ   string
	str   string
	bulk  string
	err   string
	num   int
	array []Value
}

func ReadLine(reader *bufio.Reader) ([]byte, error) {

	line, err := reader.ReadBytes('\n')
	if err != nil {
		if err == io.EOF {
			return []byte{}, err
		}
		fmt.Println("Error reading line")
		return []byte{}, err
	}

	print(string(line))

	return line, nil
}

func ParseRESP(reader *bufio.Reader) (Value, error) {
	// reader := bufio.NewReader(*conn)
	b, err := reader.ReadByte()

	if err != nil {
		fmt.Println("Error reading RESP commad")
	}

	switch b {
	case STRING:
		return parseString(reader), nil
	case BULK:
		return parseBulk(reader), nil
	// case INTEGER:
	// 	return parseInteger(reader), nil
	case ARRAY:
		return parseArray(reader), nil
		// case ERROR:
		// 	return parseError(reader), nil
	}

	return Value{}, err
}

func SerializeRESP (value Value) []byte {
	switch value.typ {
	case "string":
		var bytes []byte
		bytes = append(bytes, STRING)
		bytes = append(bytes, []byte(value.str)...)
		bytes = append(bytes, '\r', '\n')
		return bytes
	case "bulk":
		var bytes []byte
		bytes = append(bytes, BULK)
		bytes = append(bytes, []byte(strconv.Itoa(len(value.bulk)))...)
		bytes = append(bytes, '\r', '\n')
		bytes = append(bytes, []byte(value.bulk)...)
		bytes = append(bytes, '\r', '\n')

		return bytes

	case "array":
		var bytes []byte
		bytes = append(bytes, ARRAY)
		bytes = append(bytes, []byte(strconv.Itoa(len(value.array)))...)
		bytes = append(bytes, '\r', '\n')

		for _, val := range value.array {
			bytes = append(bytes, SerializeRESP(val)...)
		}

		return bytes
	}

	return []byte{}
}




func parseString(reader *bufio.Reader) Value {
	value := Value{typ: "string"}
	line, err := ReadLine(reader)

	if err != nil {
		fmt.Println("Invalid RESP string format")
		return Value{typ: "error", err: "Invalid RESP string format"}
	}

	value.str = string(line)

	fmt.Println("String: ", value.str)
	return value
}


func parseBulk(reader *bufio.Reader) Value {
	value := Value{typ: "bulk"}
	_, err := ReadLine(reader)

	if err != nil {
		fmt.Println("Invalid RESP string format")
		return Value{typ: "error", err: "Invalid RESP string format"}
	}

	// strLen, err := strconv.ParseInt(string(line), 10, 64)

	// if err != nil {
	// 	fmt.Println("Invalid RESP string format")
	// 	return Value{typ: "error", err: "Invalid RESP string format"}
	// }

	bulk, err := ReadLine(reader)
	value.bulk = string(bulk)

	return value
}

// func parseInteger(reader *bufio.Reader) Value {
// 	fmt.Println("Integer")
// }

func parseArray(reader *bufio.Reader) Value {

	value := Value{typ: "array"}
	strLen, err := ReadLine(reader)

	if err != nil {
		fmt.Println("Invalid RESP string format")
		return Value{typ: "error", err: "Invalid RESP string format"}
	}

	// fmt.Println("Array: ", string(strLen))

	arrayLen, err := strconv.ParseInt(strings.TrimSpace(string(strLen)), 10, 64)

	if err != nil {
		fmt.Println("Error converting arraylen to int")
		return Value{typ: "error", err: "error during string conversion"}
	}

	for i := 0; i < int(arrayLen); i++ {
		val, err := ParseRESP(reader)
		if err != nil {
			continue
		}
		value.array = append(value.array, val)
	}

	return value

}

// func parseError(reader *bufio.Reader) Value {
// 	fmt.Println("Error")
// }

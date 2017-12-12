package main

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
)

// IsDigit returns if the provided character is a digit
func IsDigit(s byte) bool {
	_, e := strconv.ParseUint(string(s), 10, 8)
	return e == nil
}

// ParseBigInt converts a string in 0x prefixed hex or an integer
func ParseBigInt(s string) (*big.Int, error) {
	i := new(big.Int)
	_, err := fmt.Sscan(s, i)
	return i, err
}

// hexBig is like big.Int, except when serialized to JSON it is encoded as hexadecimal
type hexBig big.Int

// UnmarshalBigInt can convert several types of JSON values to big.Int
// it supports hexadecimal, decimal strings, and integers
func UnmarshalBigInt(raw json.RawMessage) (*big.Int, error) {
	var val string
	if raw[0] == '"' {
		err := json.Unmarshal([]byte(raw), &val)
		if err != nil {
			return nil, err
		}
	} else if len(raw) > 0 && IsDigit(raw[0]) {
		val = string(raw)
	} else {
		return nil, fmt.Errorf("Invalid integer value: %v %v", string(raw), len(raw))
	}
	return ParseBigInt(string(val))
}

func (i *hexBig) UnmarshalJSON(data []byte) error {
	result, err := UnmarshalBigInt(data)
	*i = hexBig(*result)
	return err
}

func (i *hexBig) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("0x%x", (*big.Int)(i)))
}

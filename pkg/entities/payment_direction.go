package entities

import (
	"bytes"
	"encoding/json"
	"errors"
)

//go:generate stringer -type PaymentDirection

// PaymentDirection is an enum describing possible payment directions
type PaymentDirection int

const (
	// Incoming payment type is an account income
	Incoming PaymentDirection = iota

	// Outgoing payment type is an account outcome
	Outgoing
)

// MarshalJSON is used for JSON marshaling
func (pd PaymentDirection) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(pd.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON is used for JSON unmarshaling
func (pd *PaymentDirection) UnmarshalJSON(data []byte) error {
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}

	switch str {
	case "outgoing":
		*pd = Outgoing
		return nil

	case "incoming":
		*pd = Incoming
		return nil

	default:
		return errors.New("Unable to deserialize Payment direction")
	}
}

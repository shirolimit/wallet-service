// Code generated by "stringer -type PaymentDirection"; DO NOT EDIT.

package entities

import "strconv"

const _PaymentDirection_name = "IncomingOutgoing"

var _PaymentDirection_index = [...]uint8{0, 8, 16}

func (i PaymentDirection) String() string {
	if i < 0 || i >= PaymentDirection(len(_PaymentDirection_index)-1) {
		return "PaymentDirection(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _PaymentDirection_name[_PaymentDirection_index[i]:_PaymentDirection_index[i+1]]
}
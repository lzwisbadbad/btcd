// Copyright (c) 2016 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package wire

import (
	"fmt"
	"io"
)

// MsgSendHeaders implements the Message interface and represents a bitcoin
// sendheaders message.  It is used to request the peer send block headers
// rather than inventory vectors.
//
// This message has no payload and was not added until protocol versions
// starting with BIP0130Version.
type MsgSendHeaders struct{}

// BtcDecode decodes r using the bitcoin protocol encoding into the receiver.
// This is part of the Message interface implementation.
func (msg *MsgSendHeaders) BtcDecode(r io.Reader, pver uint32) error {
	if pver < BIP0130Version {
		str := fmt.Sprintf("sendheaders message invalid for protocol "+
			"version %d", pver)
		return messageError("MsgSendHeaders.BtcDecode", str)
	}

	return nil
}

// BtcEncode encodes the receiver to w using the bitcoin protocol encoding.
// This is part of the Message interface implementation.
func (msg *MsgSendHeaders) BtcEncode(w io.Writer, pver uint32) error {
	if pver < BIP0130Version {
		str := fmt.Sprintf("sendheaders message invalid for protocol "+
			"version %d", pver)
		return messageError("MsgSendHeaders.BtcEncode", str)
	}

	return nil
}

// Command returns the protocol command string for the message.  This is part
// of the Message interface implementation.
func (msg *MsgSendHeaders) Command() string {
	return CmdSendHeaders
}

// MaxPayloadLength returns the maximum length the payload can be for the
// receiver.  This is part of the Message interface implementation.
func (msg *MsgSendHeaders) MaxPayloadLength(pver uint32) uint32 {
	return 0
}

// NewMsgSendHeaders returns a new bitcoin pong message that conforms to the
// Message interface.  See MsgPong for details.
func NewMsgSendHeaders() *MsgSendHeaders {
	return &MsgSendHeaders{}
}

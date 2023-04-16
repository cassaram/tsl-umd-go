package tsl

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"unicode/utf16"
	"unicode/utf8"
)

// TSL constants for protocol
const (
	_DLE byte = 0xFE
	_STX byte = 0x02

	_OFST_PBC  int = 0
	_OFST_VERS int = 2
	_OFST_FLAG int = 3
	_OFST_SCRN int = 4
	_OFST_INDX int = 6
	_OFST_CTRL int = 8
	_OFST_LENG int = 10

	_SCRN_BCST uint16 = 0xFFFF
)

// Struct which implements a TSL 5.0 endpoint
type TSL5 struct {
	udpListener net.Listener
	stop        chan bool
	callback    func(TSLPacket)
}

func NewTSL5Instance(callback func(TSLPacket)) *TSL5 {
	p := TSL5{
		callback: callback,
	}
	return &p
}

func (p *TSL5) ListenUDP(address string, port string) error {
	var err error
	p.udpListener, err = net.Listen("udp", address+":"+port)
	if err != nil {
		return err
	}

	go p.handleUDP()
	return err
}

func (p *TSL5) handleUDP() {
	for {
		select {
		case <-p.stop:
			return
		default:
			conn, err := p.udpListener.Accept()
			if err != nil {
				// handle
				continue
			}
			go p.handleConnection(conn)
		}
	}
}

func (p *TSL5) handleConnection(conn net.Conn) {
	for {
		select {
		case <-p.stop:
			conn.Close()
			return
		default:
			buf := make([]byte, 100)
			n, err := conn.Read(buf)
			if err == nil {
				// Packet received, decode it
				pkt := DecodePacket(buf[:n])
				// Call callback in new routine
				go p.callback(pkt)
			}
		}
	}
}

func DecodePacket(data []byte) TSLPacket {
	var cleanData []byte

	// Handle DLE and STX
	for i := 0; i < len(data)-1; i += 2 {
		if data[i] != _DLE {
			// Normal packet
			cleanData = append(cleanData, data[i:i+2]...)
		} else if data[i] == _DLE && data[i+1] == _DLE {
			// Stuffed packet, unstuff and add
			cleanData = append(cleanData, data[i])
		} else if data[i] == _DLE && data[i+1] == _STX {
			// Don't add, header packet
			continue
		}
	}

	var pkt TSLPacket

	// Export parts
	pkt.Version = uint8(cleanData[_OFST_VERS])
	pkt.Flags = TSLFlags{
		Unicode:    cleanData[_OFST_FLAG]&0x1 != 0,
		ScreenData: cleanData[_OFST_FLAG]&0x2 != 0,
	}
	pkt.Screen = binary.LittleEndian.Uint16(cleanData[_OFST_SCRN : _OFST_SCRN+2])
	//packetLength := binary.LittleEndian.Uint16(cleanData[_OFST_PBC : _OFST_PBC+2])

	// Export messages from encoded packet
	for len(cleanData) > _OFST_INDX {
		msg := TSLDisplayMessage{
			Index: binary.LittleEndian.Uint16(cleanData[_OFST_INDX : _OFST_INDX+2]),
			Control: TSLControl{
				RightTally: TSLTallyColor(cleanData[_OFST_CTRL] >> 0 & 0x3),
				TextTally:  TSLTallyColor(cleanData[_OFST_CTRL] >> 2 & 0x3),
				LeftTally:  TSLTallyColor(cleanData[_OFST_CTRL] >> 4 & 0x3),
				Brightness: (cleanData[_OFST_CTRL] >> 6 & 0x3),
			},
		}
		msgLength := int(binary.LittleEndian.Uint16(cleanData[_OFST_LENG : _OFST_LENG+2]))

		if cleanData[_OFST_CTRL]&0x80 != 0 {
			// Interpret as control data
			msg.Control.ControlData = TSLControlData{}
		} else {
			// Interpret as display data
			txt, _ := getTextEncoded(cleanData[_OFST_LENG+1:_OFST_LENG+msgLength+1], pkt.Flags.Unicode)
			msg.Control.DisplayData = TSLDisplayData{
				Text: txt,
			}
		}

		// Add decoded message to packet
		pkt.Messages = append(pkt.Messages, msg)
		// Remove message from interation array
		cleanData = append(cleanData[1:_OFST_INDX], cleanData[_OFST_LENG+msgLength+1:]...)
	}

	return pkt
}

func getTextEncoded(data []byte, unicode bool) (string, error) {
	if unicode {
		// UTF-16LE
		if len(data)%2 != 0 {
			return "", fmt.Errorf("must have even length byte slice")
		}

		u16s := make([]uint16, 1)

		ret := &bytes.Buffer{}

		b8buf := make([]byte, 4)

		ldata := len(data)
		for i := 0; i < ldata; i += 2 {
			u16s[0] = uint16(data[i]) + (uint16(data[i+1]) << 8)
			r := utf16.Decode(u16s)
			n := utf8.EncodeRune(b8buf, r[0])
			ret.Write(b8buf[:n])
		}

		return ret.String(), nil

	} else {
		// ASCII
		data = bytes.Trim(data, "\u0000")
		return string(data), nil
	}
}

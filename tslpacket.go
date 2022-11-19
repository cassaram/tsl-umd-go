package tsl

// Represents a TSL 5.0 packet stripped to useful information
type TSLPacket struct {
	Version  uint8
	Flags    TSLFlags
	Screen   uint16
	Messages []TSLDisplayMessage
}

// Struct representing flags within a TSL 5.0 packet
type TSLFlags struct {
	Unicode    bool
	ScreenData bool
}

// Struct representing a display message defined by TSL 5.0
type TSLDisplayMessage struct {
	Index   uint16
	Control TSLControl
}

// Struct representing control data defined by TSL 5.0
type TSLControl struct {
	RightTally  TSLTallyColor
	TextTally   TSLTallyColor
	LeftTally   TSLTallyColor
	Brightness  uint8
	DisplayData TSLDisplayData
	ControlData TSLControlData
}

// Constant defines for tally colors
type TSLTallyColor uint8

const (
	OFF   TSLTallyColor = 0
	RED   TSLTallyColor = 1
	GREEN TSLTallyColor = 2
	AMBER TSLTallyColor = 3
)

// Struct for storing display data
type TSLDisplayData struct {
	Text string
}

// Struct for storing control data
type TSLControlData struct {
	// NOT IMPLEMENTED IN TSL V5.0
}

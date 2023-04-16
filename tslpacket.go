package tsl

// Represents a TSL 5.0 packet stripped to useful information
type TSLPacket struct {
	Version  uint8               `json:"version"`
	Flags    TSLFlags            `json:"flags"`
	Screen   uint16              `json:"screen"`
	Messages []TSLDisplayMessage `json:"messages"`
}

// Struct representing flags within a TSL 5.0 packet
type TSLFlags struct {
	Unicode    bool `json:"unicode"`
	ScreenData bool `json:"screen-data"`
}

// Struct representing a display message defined by TSL 5.0
type TSLDisplayMessage struct {
	Index   uint16     `json:"index"`
	Control TSLControl `json:"control"`
}

// Struct representing control data defined by TSL 5.0
type TSLControl struct {
	RightTally  TSLTallyColor  `json:"right-tally"`
	TextTally   TSLTallyColor  `json:"text-tally"`
	LeftTally   TSLTallyColor  `json:"left-tally"`
	Brightness  uint8          `json:"brightness"`
	DisplayData TSLDisplayData `json:"display-data"`
	ControlData TSLControlData `json:"control-data"`
}

// Constant defines for tally colors
type TSLTallyColor uint8

const (
	TALLY_OFF   TSLTallyColor = 0
	TALLY_RED   TSLTallyColor = 1
	TALLY_GREEN TSLTallyColor = 2
	TALLY_AMBER TSLTallyColor = 3
)

// Struct for storing display data
type TSLDisplayData struct {
	Text string `json:"text"`
}

// Struct for storing control data
type TSLControlData struct {
	// NOT IMPLEMENTED IN TSL V5.0
}

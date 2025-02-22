package common

type Platform struct {
	Name    string
	Version string
}

type Board struct {
	Platform Platform
	Fqbn     string
}

var ArduinoAvr = Platform{
	Name:    "arduino:avr",
	Version: "1.8.6",
}

var ArduinoRenesasUno = Platform{
	Name:    "arduino:renesas_uno",
	Version: "1.3.2",
}

var Boards = map[string]Board{
	"arduino:avr:uno": {
		Platform: ArduinoAvr,
		Fqbn:     "arduino:avr:uno",
	},
	"arduino:renesas_uno:minima": {
		Platform: ArduinoRenesasUno,
		Fqbn:     "arduino:renesas_uno:minima",
	},
	"arduino:renesas_uno:unor4wifi": {
		Platform: ArduinoRenesasUno,
		Fqbn:     "arduino:renesas_uno:unor4wifi",
	},
}

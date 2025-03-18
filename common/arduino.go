package common

type Platform struct {
	Name          string
	Version       string
	AdditionalUrl string
}

type Board struct {
	Platform Platform
	Fqbn     string
	Exts     []string
}

var ArduinoAvr = Platform{
	Name:          "arduino:avr",
	Version:       "1.8.6",
	AdditionalUrl: "",
}

var ArduinoRenesasUno = Platform{
	Name:          "arduino:renesas_uno",
	Version:       "1.3.2",
	AdditionalUrl: "",
}

var Rp2040 = Platform{
	Name:          "rp2040:rp2040",
	Version:       "4.5.0",
	AdditionalUrl: "https://github.com/earlephilhower/arduino-pico/releases/download/global/package_rp2040_index.json",
}

var Boards = map[string]Board{
	"arduino:avr:uno": {
		Platform: ArduinoAvr,
		Fqbn:     "arduino:avr:uno",
		Exts:     []string{"hex", "elf", "eep"},
	},
	"arduino:renesas_uno:minima": {
		Platform: ArduinoRenesasUno,
		Fqbn:     "arduino:renesas_uno:minima",
		Exts:     []string{"hex", "elf", "bin"},
	},
	"rp2040:rp2040:rpipico": {
		Platform: Rp2040,
		Fqbn:     "rp2040:rp2040:rpipico",
		Exts:     []string{"bin", "elf", "uf2"},
	},
	"rp2040:rp2040:rpipicow": {
		Platform: Rp2040,
		Fqbn:     "rp2040:rp2040:rpipicow",
		Exts:     []string{"bin", "elf", "uf2"},
	},
}

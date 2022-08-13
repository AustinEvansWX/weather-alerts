package alerts

import (
	"strconv"
)

func GetHailSizeEquivalent(hailSize string) string {
	equivalent := ""

	hailSizeValue, _ := strconv.ParseFloat(hailSize, 64)

	switch hailSizeValue {
	case 0.25:
		equivalent = "| Pea Size"

	case 0.5:
		equivalent = "| Peanut Size"

	case 0.75:
		equivalent = "| Penny Size"

	case 0.88:
		equivalent = "| Nickel Size"

	case 1:
		equivalent = "| Quarter Size"

	case 1.25:
		equivalent = "| Half Dollar Size"

	case 1.5:
		equivalent = "| Ping Pong Ball Size"

	case 1.75:
		equivalent = "| Golf Ball Size"

	case 2:
		equivalent = "| Lime Size"

	case 2.5:
		equivalent = "| Tennis Ball Size"

	case 2.75:
		equivalent = "| Baseball Size"

	case 3:
		equivalent = "| Apple Size"

	case 4:
		equivalent = "| Softball Size"

	case 4.5:
		equivalent = "| Grapefruit Size"
	}

	if hailSizeValue > 4.5 {
		equivalent = "| GORILLA SIZE"
	}

	return equivalent
}

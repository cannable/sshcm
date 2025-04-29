package misc

import (
	"fmt"
)

func StringTrimmer(s string, len int) string {
	f := fmt.Sprintf("%-*s", len, s)

	return f[:len]
}

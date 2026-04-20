package palette

import (
	"fmt"
	"testing"

	"github.com/mermonia/chromatika/internal/colors"
)

func TestDeltaE00(t *testing.T) {
	a := colors.LCHab{
		L: 7.8757,
		C: 19.412,
		H: 292.71,
	}

	b := colors.LCHab{
		L: 4.481,
		C: 12.9203,
		H: 295.781,
	}

	difference := deltaE00(a, b)
	fmt.Println(difference)
}

package palette

import (
	"fmt"
	"testing"

	"github.com/mermonia/chromatika/internal/colors"
)

func TestDeltaE00(t *testing.T) {
	a := &colors.LCHab{
		L: 1,
		C: 3,
		H: 5,
	}

	b := &colors.LCHab{
		L: 2,
		C: 4,
		H: 6,
	}

	difference := deltaE00(a,b)
	fmt.Print(difference)
}

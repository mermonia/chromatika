package format

import (
	"bytes"
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/mermonia/chromatika/internal/variants"
)

type VariantsFormatter interface {
	Format(*variants.Variants) string
}

type TOMLVariantsFormatter struct{}
type ASCIIVariantsFormatter struct{}

func (*TOMLVariantsFormatter) Format(vars *variants.Variants) string {
	var buf bytes.Buffer
	if err := toml.NewEncoder(&buf).Encode(vars); err != nil {
		return ""
	}
	return buf.String()
}

func (*ASCIIVariantsFormatter) Format(vars *variants.Variants) string {
	res := fmt.Sprintf("%s%s%s%s%s%s%s%s%s%s\nBase Index: %d",
		vars.Variant50,
		vars.Variant100,
		vars.Variant200,
		vars.Variant300,
		vars.Variant400,
		vars.Variant500,
		vars.Variant600,
		vars.Variant700,
		vars.Variant800,
		vars.Variant900,
		vars.BaseIndex)

	return res
}

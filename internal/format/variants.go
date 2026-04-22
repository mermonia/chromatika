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
		vars.Variant_50,
		vars.Variant_100,
		vars.Variant_200,
		vars.Variant_300,
		vars.Variant_400,
		vars.Variant_500,
		vars.Variant_600,
		vars.Variant_700,
		vars.Variant_800,
		vars.Variant_900,
		vars.BaseIndex)

	return res
}

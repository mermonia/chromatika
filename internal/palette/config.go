package palette

type GenerationConfig struct {
	ImagePath string
	DarkMode  bool

	QuantInterval int
	ScaleWidth    int

	Fuzziness float64
	Threshold float64
	MaxIter   int
}

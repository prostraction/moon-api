package phase

type PhaseResp struct {
	BeginDay     *Phase
	Current      *Phase
	EndDay       *Phase
	Illumination Illumination
}

type Phase struct {
	Name          string
	NameLocalized string
	Emoji         string
	IsWaxing      bool
}

type Illumination struct {
	BeginDay float64
	Current  float64
	EndDay   float64
}

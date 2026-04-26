package zodiac

type ZodiacDetailed struct {
	Name          string
	NameLocalized string
	Emoji         string
	Begin         *any
	End           *any
}

type Zodiacs struct {
	Count  int
	Zodiac []ZodiacDetailed
}

type Zodiac struct {
	Name          string
	NameLocalized string
	Emoji         string
}

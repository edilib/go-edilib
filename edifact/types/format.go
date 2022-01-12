package types

type Format struct {
	SkipNewLineAfterSegment       bool
	UnaAllowed                    bool
	ComponentDataElementSeperator rune
	DataElementSeperator          rune
	DecimalMark                   rune
	ReleaseCharacter              rune
	RepetitionSeperator           rune
	SegmentTerminator             rune
}

func UnEdifactFormat() Format {
	return Format{
		SkipNewLineAfterSegment:       false,
		UnaAllowed:                    true,
		ComponentDataElementSeperator: ':',
		DataElementSeperator:          '+',
		DecimalMark:                   '.',
		ReleaseCharacter:              '?',
		RepetitionSeperator:           '*',
		SegmentTerminator:             '\''}
}

func X12Format() Format {
	return Format{
		SkipNewLineAfterSegment:       false,
		UnaAllowed:                    false,
		ComponentDataElementSeperator: '>',
		DataElementSeperator:          '*',
		DecimalMark:                   '.',
		ReleaseCharacter:              '?',
		RepetitionSeperator:           ' ',
		SegmentTerminator:             '~'}
}

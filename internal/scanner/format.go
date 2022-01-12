package scanner

type Format struct {
	componentDataElementSeperator rune
	dataElementSeperator          rune
	decimalMark                   rune
	releaseCharacter              rune
	repetitionSeperator           rune
	segmentTerminator             rune
}

func DefaultFormat() Format {
	return Format{componentDataElementSeperator: ':',
		dataElementSeperator: '+',
		decimalMark:          '.',
		releaseCharacter:     '?',
		repetitionSeperator:  '*',
		segmentTerminator:    '\''}
}

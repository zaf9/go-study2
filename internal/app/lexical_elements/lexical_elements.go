package lexical_elements

// Display is the master function that calls all the individual display functions
// for each lexical element topic.
func Display() {
	DisplayComments()
	DisplayTokens()
	DisplaySemicolons()
	DisplayIdentifiers()
	DisplayKeywords()
	DisplayOperators()
	DisplayIntegers()
	DisplayFloats()
	DisplayImaginary()
	DisplayRunes()
	DisplayStrings()
}

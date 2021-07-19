package utility

import "regexp"

var re = regexp.MustCompile("(?s)//.*?\n|/\\*.*?\\*/")

func TripComment(jsonc []byte) []byte {
	return re.ReplaceAll(jsonc, nil)
}

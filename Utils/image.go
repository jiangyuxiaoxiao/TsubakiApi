package Utils

import "strings"

var ImageString []string = []string{"BMP", "DIB", "PCP", "DIF", "WMF", "GIF", "JPG", "JPEG",
	"TIF", "EPS", "PSD", "CDR", "IFF", "TGA", "PCD", "MPT", "PNG"}

func IsImage(name string) bool {
	name = strings.ToUpper(name)
	for _, end := range ImageString {
		if strings.HasSuffix(name, end) {
			return true
		}
	}
	return false
}

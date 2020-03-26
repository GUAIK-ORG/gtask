package image

import (
	"net/http"
)

func CheckImage(data []byte) bool {
	mine := http.DetectContentType(data)
	switch mine {
	case "image/jpeg", "image/jpg":
		return true
	case "image/gif":
		return true
	case "image/png":
		return true
	}
	return false
}

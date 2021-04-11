package nasa

import (
	"gorm.io/gorm"
)

// Apod saves the information of apod api
/*
mockup data from nasa api
	{
		"copyright": "Daniel L\u00f3pezIAC",
		"date": "2018-01-08",
		"explanation": "What are those red clouds surrounding the Andromeda galaxy? This galaxy, M31, is often imaged by planet Earth-based astronomers. As the nearest large spiral galaxy, it is a familiar sight with dark dust lanes, bright yellowish core, and spiral arms traced by clouds of bright blue stars.  A mosaic of well-exposed broad and narrow-band image data, this colorful portrait of our neighboring island universe offers strikingly unfamiliar features though, faint reddish clouds of glowing ionized hydrogen gas in the same wide field of view. These ionized hydrogen clouds surely lie in the foreground of the scene, well within our Milky Way Galaxy. They are likely associated with the pervasive, dusty interstellar cirrus clouds scattered hundreds of light-years above our own galactic plane.   Free APOD Lecture Tomorrow: January 9 at the National Harbor near Washington, DC",
		"hdurl": "https://apod.nasa.gov/apod/image/1801/M31Clouds_DLopez_1500.jpg",
		"media_type": "image",
		"service_version": "v1",
		"title": "Clouds of Andromeda",
		"url": "https://apod.nasa.gov/apod/image/1801/M31Clouds_DLopez_960.jpg"
	}
*/
type Apod struct {
	gorm.Model

	Title       string `json:"title"`
	MediaType   string `gorm:"index" json:"media_type"`
	ServiceVer  string `json:"service_version"`
	CopyRight   string `json:"copyright"`
	Date        string `gorm:"index" json:"date"`
	Explanation string `json:"explanation"`
	ImgHDURL    string `json:"hdurl"`
	ImgURL      string `json:"url"`
}

// Map returns data bytes as map
func (a Apod) Map() map[string]string {
	result := map[string]string{
		"copyright":   a.CopyRight,
		"date":        a.Date,
		"explanation": a.Explanation,
		"url":         a.ImgURL,
	}

	if a.MediaType == "image" {
		result["hdurl"] = a.ImgHDURL
	}

	return result
}

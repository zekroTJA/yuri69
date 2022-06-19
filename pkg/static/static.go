package static

import "github.com/gabriel-vasile/mimetype"

const (
	BucketSounds = "sounds"
	BucketTemp   = "temp"
	SoundsMime   = "audio/ogg"
)

var SoundsMimeType mimetype.MIME

func init() {
	typ := mimetype.Lookup(SoundsMime)
	if typ == nil {
		panic("Invalid SoundsMime type")
	}
	SoundsMimeType = *typ
}

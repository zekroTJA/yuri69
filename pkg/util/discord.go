package util

import "github.com/bwmarrin/discordgo"

func IsErrCode(err error, code int) bool {
	apiErr, ok := err.(*discordgo.RESTError)
	return ok && apiErr.Message.Code == code
}

package models

type Guild struct {
	GuildID    string
	Name       string
	BannerHash string
	IconHash   string
}

type GuildOut struct {
	GuildID                 string
	Name                    string
	BannerHash              string
	IconHash                string
	GuildBannerResourcePath string
	GuildIconResourcePath   string
}

type Guilds struct {
	Guilds []GuildOut
}

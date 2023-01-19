package utils

type DiscordMember struct {
	Avatar string      `json:"avatar,omitempty"`
	Nick   *string     `json:"nick,omitempty"`
	Roles  []string    `json:"roles,omitempty"`
	User   DiscordUser `json:"user"`
}

type DiscordUser struct {
	Id            string `json:"id"`
	Username      string `json:"username"`
	Avatar        string `json:"avatar,omitempty"`
	Discriminator string `json:"discriminator"`
}

type DiscordInfo struct {
	UserId  string `yaml:"userId"`
	Name    string `yaml:"name"`
	Discord string `yaml:"discord"`
}

func (d DiscordInfo) Valid() error {
	return nil
}

func (d DiscordInfo) Type() string {
	return "discord-user"
}

type DiscordToken struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

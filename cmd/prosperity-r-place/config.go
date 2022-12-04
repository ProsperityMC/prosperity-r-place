package main

type Config struct {
	Listen string       `yaml:"listen"`
	Auth   AuthConfig   `yaml:"auth"`
	Login  LoginConfig  `yaml:"login"`
	Slots  []SlotConfig `yaml:"slots"`
}

type AuthConfig struct {
	Issuer string `yaml:"issuer"`
	Key    string `yaml:"key"`
	Public string `yaml:"public"`
}

type LoginConfig struct {
	Id          string      `yaml:"id"`
	Token       string      `yaml:"token"`
	RedirectUrl string      `yaml:"redirectUrl"`
	BaseUrl     string      `yaml:"baseUrl"`
	Guild       GuildConfig `yaml:"guild"`
}

type GuildConfig struct {
	Id    string   `yaml:"id"`
	Roles []string `yaml:"roles"`
}

type SlotConfig struct {
	Name   string `yaml:"name"`
	Width  int    `yaml:"width"`
	Height int    `yaml:"height"`
}

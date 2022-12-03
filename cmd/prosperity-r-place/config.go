package main

type Config struct {
	Listen string       `yaml:"listen"`
	Login  LoginConfig  `yaml:"login"`
	Slots  []SlotConfig `yaml:"slots"`
}

type LoginConfig struct {
	Id          string `yaml:"id"`
	Token       string `yaml:"token"`
	RedirectUrl string `yaml:"redirectUrl"`
}

type SlotConfig struct {
	Name   string `yaml:"name"`
	Width  int    `yaml:"width"`
	Height int    `yaml:"height"`
}

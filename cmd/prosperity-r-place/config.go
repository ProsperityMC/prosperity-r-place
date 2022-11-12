package main

type Config struct {
	Listen string       `yaml:"listen"`
	Slots  []SlotConfig `yaml:"slots"`
}

type SlotConfig struct {
	Name   string `yaml:"name"`
	Width  int    `yaml:"width"`
	Height int    `yaml:"height"`
}

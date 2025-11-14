package bot

import "github.com/nanachi-sh/go-mc/bot/forge"

type Option interface {
	Merge(*Client)
}

func WithForge(mods forge.ModList) Option {
	return &Forge{
		mods: mods,
	}
}

type Forge struct {
	mods forge.ModList
}

func (s *Forge) Merge(c *Client) {
	c.forge_info = &Forge{
		mods: s.mods,
	}
}

func WithUsername(name string) Option {
	return &username{name}
}

type username struct{ string }

func (s *username) Merge(c *Client) {
	c.username = s.string
}

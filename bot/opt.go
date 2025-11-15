package bot

import (
	"encoding/hex"

	"github.com/nanachi-sh/go-mc/bot/forge"
)

type Option interface {
	Merge(*Client)
}

func WithForge(modchan_payload, modlist_payload string) Option {
	x := new(Forge)
	v, err := hex.DecodeString(modlist_payload)
	if err != nil {
		panic(err)
	}
	x.modlist = v
	v, err = hex.DecodeString(modchan_payload)
	if err != nil {
		panic(err)
	}
	x.modchan = v
	return x
}

type Forge struct {
	mods forge.ModList

	modlist []byte //mod列表
	modchan []byte //mod通信通道
}

func (s *Forge) Merge(c *Client) {
	// c.forge_info = &Forge{
	// 	mods: s.mods,
	// }
	c.forge_info = s
}

func WithUsername(name string) Option {
	return &username{name}
}

type username struct{ string }

func (s *username) Merge(c *Client) {
	c.username = s.string
}

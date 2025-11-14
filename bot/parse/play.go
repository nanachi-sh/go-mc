package parse

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/nanachi-sh/go-mc/constants"
)

type splay struct {
	ChatMessage   *splay_chatMessage
	PluginMessage *struct{}
}

type splay_chatMessage struct {
}

func (p *Parser) parsePlay() (*splay, error) {
	s, e := index(p.buffer)
	if s == -1 || e == -1 {
		return nil, nil
	}
	p.buffer.Next(s)
	buffer := bytes.NewBuffer(p.buffer.Next(e))
	id := constants.ServerID(buffer.Next(1)[0])
	switch id {
	default:
		fmt.Printf("id: %x\n", id)
		return nil, errors.ErrUnsupported
	case constants.SPlguinMessage: //没有实现的必要
		return &splay{PluginMessage: &struct{}{}}, nil
	}
}

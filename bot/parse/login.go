package parse

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/nanachi-sh/go-mc/constants"
)

type slogin struct {
	LoginSuccess *slogin_loginSuccess
}

type slogin_loginSuccess struct {
	Username string
	UUID     string
}

func (p *Parser) parseLogin() (*slogin, error) {
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
	case constants.SLoginSuccess:
		// got uuid, length < 128
		uuid := buffer.Next(int(buffer.Next(1)[0]))
		// got username
		username := buffer.Next(int(buffer.Next(1)[0]))
		return &slogin{
			LoginSuccess: &slogin_loginSuccess{
				Username: string(username),
				UUID:     string(uuid),
			},
		}, nil
	}
}

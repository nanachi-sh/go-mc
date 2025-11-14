package parse

import (
	"bytes"
	"fmt"

	"github.com/nanachi-sh/go-mc/util"
)

type Parser struct {
	buffer  bytes.Buffer
	logined bool
	ch      chan *ServerResponse
}

type ServerResponse struct {
	Login *slogin
	Play  *splay
}

func NewParser() *Parser {
	return &Parser{
		ch: make(chan *ServerResponse, 1),
	}
}

func (p *Parser) Put(b []byte) {
	p.buffer.Write(b)
	p.parse()
}

func (p *Parser) Read() *ServerResponse {
	return <-p.ch
}

func (p *Parser) parse() {
	if p.logined {
		v, err := p.parsePlay()
		if err != nil {
			fmt.Println(err)
			return
		}
		if v != nil {
			go func() { p.ch <- &ServerResponse{Play: v} }()
		}
	} else {
		v, err := p.parseLogin()
		if err != nil {
			fmt.Println(err)
			return
		}
		if v != nil {
			go func() { p.ch <- &ServerResponse{Login: v} }()
		}
	}
}

func (p *Parser) LoginFinish() {
	p.logined = true
}

func index(buffer bytes.Buffer) (int, int) {
	for i := 1; ; i++ {
		//
		if buffer.Len() < i {
			break
		}
		// 获取数据包长度
		l, err := util.ReadVarint(buffer.Bytes()[:i])
		if err != nil {
			continue
		}
		// 数据包完整性检查
		if int(l) < buffer.Len()-i {
			return i, -1
		}
		return i, int(l)
	}
	return -1, -1
}

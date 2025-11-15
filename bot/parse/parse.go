package parse

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/nanachi-sh/go-mc/util"
)

type Parser struct {
	lock    sync.Mutex
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
	p.lock.Lock()
	defer p.lock.Unlock()
	p.buffer.Write(b)
	p.parse()
}

func (p *Parser) Read() *ServerResponse {
	return <-p.ch
}

func (p *Parser) parse() {
	for {
		if p.logined {
			v, err := p.parsePlay()
			if err != nil {
				fmt.Println(err)
				return
			}
			if v != nil {
				go func() { p.ch <- &ServerResponse{Play: v} }()
			} else {
				break
			}
		} else {
			v, err := p.parseLogin()
			if err != nil {
				fmt.Println(err)
				return
			}
			if v != nil {
				go func() { p.ch <- &ServerResponse{Login: v} }()
			} else {
				break
			}
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
		if buffer.Len()-i < int(l) {
			return i, -1
		}
		return i, int(l)
	}
	return -1, -1
}

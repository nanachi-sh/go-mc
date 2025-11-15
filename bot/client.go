package bot

import (
	"encoding/binary"
	"errors"
	"io"
	_net "net"
	"strconv"
	"sync"

	"github.com/nanachi-sh/go-mc/bot/parse"
	"github.com/nanachi-sh/go-mc/constants"
	"github.com/nanachi-sh/go-mc/net"
	"github.com/nanachi-sh/go-mc/util"
)

type Client struct {
	parse *parse.Parser
	conn  *net.Conn
	once  sync.Once

	forge_info *Forge
	dialed     bool
	address    string
	username   string
	uuid       string
	funcs      []func(parse.ServerResponse)
}

func NewClient(address string, opts ...Option) (*Client, error) {
	c := &Client{
		address:  address,
		username: "Stave",
		parse:    parse.NewParser(),
	}
	for _, v := range opts {
		v.Merge(c)
	}
	return c, nil
}

func (c *Client) Register(f func(parse.ServerResponse)) error {
	if !c.dialed {
		return errors.New("请先连接服务器")
	}
	c.funcs = append(c.funcs, f)
	return nil
}

func (c *Client) Dial() error {
	return c.dial()
}

func (c *Client) dial() (err error) {
	c.once.Do(func() {
		conn, err := net.Dial("tcp", c.address)
		if err != nil {
			return
		}
		c.conn = conn
		// 监听数据
		go c.listen()
		// 握手
		if err := c.handshake(); err != nil {
			return
		}
		// 登录
		if err := c.login(); err != nil {
			return
		}
		// 接收UUID和玩家名
		resp := c.parse.Read()
		c.username = resp.Login.LoginSuccess.Username
		c.uuid = resp.Login.LoginSuccess.UUID
		c.parse.LoginFinish()
		// forge
		if c.forge_info != nil {
			// 接收通道
			c.parse.Read()
			// 随便发送
			// if _, err := c.conn.Write(constants.CPlguinMessage, constants.CPluginChannel); err != nil {
			// 	panic(err)
			// }
			// 发送mod通信通道
			if _, err := c.conn.WriteRaw(c.forge_info.modchan); err != nil {
				panic(err)
			}
			// 确定
			if _, err := c.conn.Write(constants.CPlguinMessage, constants.FMLHS); err != nil {
				panic(err)
			}
			// 发送mod列表
			// {
			// 	var buffer bytes.Buffer
			// 	for _, v := range c.forge_info.mods {
			// 		buffer.WriteByte(byte(len(v.ModId)))
			// 		buffer.WriteString(v.ModId)
			// 		buffer.WriteByte(byte(len(v.Version)))
			// 		buffer.WriteString(v.Version)
			// 	}
			// 	data := []byte{}
			// 	// FML|HS
			// 	data = append(data, 6)
			// 	data = append(data, []byte("FML|HS")...)
			// 	// next length
			// 	a := make([]byte, 2)
			// 	binary.BigEndian.PutUint16(a, uint16(2+1+buffer.Len()))
			// 	data = append(data, a...)
			// 	// id?
			// 	data = append(data, 2)
			// 	// mod num
			// 	b := make([]byte, 2)
			// 	binary.BigEndian.PutUint16(b, uint16(len(c.forge_info.mods)))
			// 	data = append(data, b...)
			// 	// mod list
			// 	data = append(data, buffer.Bytes()...)
			// 	if _, err := c.conn.Write(constants.CPlguinMessage, data); err != nil {
			// 		panic(err)
			// 	}
			// }
			if _, err := c.conn.WriteRaw(c.forge_info.modlist); err != nil {
				panic(err)
			}
			// 接收mod列表
			c.parse.Read()
			// 确定
			if _, err := c.conn.Write(constants.CPlguinMessage, constants.FMLHS); err != nil {
				panic(err)
			}
			// 各种方块ID等等
			c.parse.Read()
			if _, err := c.conn.Write(constants.CPlguinMessage, constants.FMLHS); err != nil {
				panic(err)
			}
		} else {
			c.parse.Read()
		}
		// 登录完毕
		c.dialed = true
		go c.back()
	})
	return
}

func (c *Client) SendMessage(msg string) error {
	data := []byte{}
	data = append(data, util.PutVarint(len(msg))...)
	data = append(data, []byte(msg)...)
	_, err := c.conn.Write(constants.CChatMessage, data)
	return err
}

func (c *Client) back() {
	for {
		resp := c.parse.Read()
		for _, v := range c.funcs {
			go v(*resp)
		}
		switch {
		case resp.Play.KeepAlive != nil:
			a := make([]byte, 4)
			binary.BigEndian.PutUint32(a, uint32(*resp.Play.KeepAlive))
			if _, err := c.conn.Write(constants.CKeepAlive, a); err != nil {
				panic(err)
			}
		}
	}
}

func (c *Client) listen() {
	base := 1024
	size := base
	b := make([]byte, size)
	for {
		n, err := c.conn.Read(b)
		if err != nil {
			if err != io.EOF {
				panic(err)
			}
		}
		if n == size {
			size *= 2
		} else {
			size = base
		}
		if n == 0 {
			continue
		}
		c.parse.Put(b[:n])
		b = make([]byte, size)
	}
}

func (c *Client) handshake() error {
	host, port, err := _net.SplitHostPort(c.conn.LocalAddr().String())
	if err != nil {
		return err
	}
	data := []byte{}
	// v
	data = append(data, constants.ProtocolVersion)
	// len(host)
	data = append(data, byte(len(host)))
	// host
	data = append(data, []byte(host)...)
	// port
	i, _ := strconv.ParseInt(port, 10, 0)
	a := make([]byte, 2)
	binary.BigEndian.PutUint16(a, uint16(i))
	data = append(data, a...)
	// login
	data = append(data, 2)
	// 发送
	_, err = c.conn.Write(constants.CHandshake, data)
	return err
}

func (c *Client) login() error {
	data := []byte{}
	// len(username)
	data = append(data, byte(len(c.username)))
	// username
	data = append(data, []byte(c.username)...)
	// w
	_, err := c.conn.Write(constants.CLogin, data)
	return err
}

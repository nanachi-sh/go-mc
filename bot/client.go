package bot

import (
	"bytes"
	"encoding/binary"
	_net "net"
	"strconv"
	"sync"

	"github.com/nanachi-sh/go-mc/bot/parse"
	"github.com/nanachi-sh/go-mc/constants"
	"github.com/nanachi-sh/go-mc/net"
)

type Client struct {
	once sync.Once

	forge_info *Forge
	conn       *net.Conn
	address    string
	username   string
	uuid       string
	parse      *parse.Parser
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
		// forge
		if c.forge_info != nil {
			// 接收通道
			c.parse.Read()
			// 随便发送
			if _, err := c.conn.Write(constants.CPlguinMessage, constants.CPluginChannel); err != nil {
				panic(err)
			}
			// 确定
			if _, err := c.conn.Write(constants.CPlguinMessage, constants.FMLHS); err != nil {
				panic(err)
			}
			// 发送mod列表
			{
				var buffer bytes.Buffer
				for _, v := range c.forge_info.mods {
					buffer.WriteByte(byte(len(v.ModId)))
					buffer.WriteString(v.ModId)
					buffer.WriteByte(byte(len(v.Version)))
					buffer.WriteString(v.Version)
				}
				data := []byte{}
				// FML|HS
				data = append(data, 6)
				data = append(data, []byte("FML|HS")...)
				// next length
				a := make([]byte, 2)
				binary.BigEndian.PutUint16(a, uint16(2+1+buffer.Len()))
				data = append(data, a...)
				// id?
				data = append(data, 2)
				// mod num
				b := make([]byte, 2)
				binary.BigEndian.PutUint16(b, uint16(len(c.forge_info.mods)))
				data = append(data, b...)
				// mod list
				data = append(data, buffer.Bytes()...)
				if _, err := c.conn.Write(constants.CPlguinMessage, data); err != nil {
					panic(err)
				}
			}
			// 接收mod列表
			c.parse.Read()
			// 确定
			if _, err := c.conn.Write(constants.CPlguinMessage, constants.FMLHS); err != nil {
				panic(err)
			}
		}
		// 登录完毕
	})
	return
}

func (c *Client) SendMessage(msg string) error {
	data := []byte{}
	data = append(data, byte(len(msg)))
	data = append(data, []byte(msg)...)
	_, err := c.conn.Write(constants.CChatMessage, data)
	return err
}

func (c *Client) listen() {
	for {
		b, err := c.conn.Read()
		if err != nil {
			panic(err)
		}
		c.parse.Put(b)
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

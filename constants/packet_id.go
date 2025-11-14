package constants

import "encoding/hex"

const ProtocolVersion = 5

// 标准的payload：length + packet_id + data
// 客户端发送的packet_id
type ClientID byte

const (
	CHandshake     ClientID = 0x00 //protocol version + len(host) + host + port + next state(1 status, 2 login)
	CLogin         ClientID = 0x00 //len(username) + username
	CPlguinMessage ClientID = 0x17
	CChatMessage   ClientID = 0x01
)

// 服务端返回的packet_id
type ServerID byte

const (
	SPlguinMessage ServerID = 0x3F
	SLoginSuccess  ServerID = 0x02
	SChatMessage   ClientID = 0x02
)

var (
	FMLHS          []byte
	CPluginChannel []byte
)

func init() {
	data, _ := hex.DecodeString("0852454749535445520014464d4c7c485300464d4c00464d4c00464f524745")
	CPluginChannel = data
	data, _ = hex.DecodeString("06464d4c7c485300020102")
	FMLHS = data
}

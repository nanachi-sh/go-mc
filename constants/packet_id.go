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
	CKeepAlive     ClientID = 0x00
)

// 服务端返回的packet_id
type ServerID byte

const (
	SPlguinMessage             ServerID = 0x3F
	SLoginSuccess              ServerID = 0x02
	SChatMessage               ServerID = 0x02
	SEntityHeadLook            ServerID = 0x19
	SEntityVelocity            ServerID = 0x12
	SMapChunkBulk              ServerID = 0x26
	SEntityRelativeMove        ServerID = 0x15
	STimeUpdate                ServerID = 0x03
	SSpawnMob                  ServerID = 0x0F
	SEntityMetadata            ServerID = 0x1C
	SEntityProperties          ServerID = 0x20
	SSpawnObject               ServerID = 0x0E
	SEntityLookandRelativeMove ServerID = 0x17
	SSpawnPosition             ServerID = 0x05
	SPlayerAbilities           ServerID = 0x39
	SHeldItemChange            ServerID = 0x09
	SStatistics                ServerID = 0x37
	SPlayerListItem            ServerID = 0x38
	SPlayerPositionAndLook     ServerID = 0x08
	SWindowItems               ServerID = 0x30
	SSetSlot                   ServerID = 0x2F
	SSpawnPlayer               ServerID = 0x0C
	SSoundEffect               ServerID = 0x29
	SEntityLook                ServerID = 0x16
	SEntityEquipment           ServerID = 0x04
	SJoinGame                  ServerID = 0x01
	SEntityTeleport            ServerID = 0x18
	SKeepAlive                 ServerID = 0x00
	SChangeGameState           ServerID = 0x2B
	SMultiBlockChange          ServerID = 0x22
	SDestroyEntities           ServerID = 0x13
	SBlockChange               ServerID = 0x23
	SEntityStatus              ServerID = 0x1A
	SSpawnExperienceOrb        ServerID = 0x11
	SCollectItem               ServerID = 0x0D
	SAnimation                 ServerID = 0x0B
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

package parse

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	jsonvalue "github.com/Andrew-M-C/go.jsonvalue"
	"github.com/nanachi-sh/go-mc/constants"
	"github.com/nanachi-sh/go-mc/util"
)

type splay struct {
	ChatMessage               *splay_chatMessage
	PluginMessage             *struct{}
	EntityHeadLook            *struct{}
	EntityVelocity            *struct{}
	MapChunkBulk              *struct{}
	EntityRelativeMove        *struct{}
	TimeUpdate                *struct{}
	SpawnMob                  *struct{}
	EntityMetadata            *struct{}
	EntityProperties          *struct{}
	SpawnObject               *struct{}
	EntityLookandRelativeMove *struct{}
	SpawnPosition             *struct{}
	PlayerAbilities           *struct{}
	HeldItemChange            *struct{}
	Statistics                *struct{}
	PlayerListItem            *struct{}
	PlayerPositionAndLook     *struct{}
	WindowItems               *struct{}
	SetSlot                   *struct{}
	SpawnPlayer               *struct{}
	SoundEffect               *struct{}
	EntityLook                *struct{}
	EntityEquipment           *struct{}
	JoinGame                  *struct{}
	EntityTeleport            *struct{}
	KeepAlive                 *int
	ChangeGameState           *struct{}
	MultiBlockChange          *struct{}
	DestroyEntities           *struct{}
	BlockChange               *struct{}
	EntityStatus              *struct{}
	SpawnExperienceOrb        *struct{}
	CollectItem               *struct{}
	Animation                 *struct{}
}

type splay_chatMessage struct {
	Translate string
	Json      *jsonvalue.V
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
		fmt.Printf("id: %x, b: %x\n", id, buffer.Bytes())
		return nil, errors.ErrUnsupported
	case constants.SPlguinMessage: //没有实现的必要
		return &splay{PluginMessage: &struct{}{}}, nil
	case constants.SChatMessage:
		// got json
		s, e := index(*buffer)
		if s == -1 || e == -1 {
			return nil, errors.New("json不完整，异常错误")
		}
		buffer.Next(s)
		j_str := buffer.Next(e)
		j, err := jsonvalue.Unmarshal(j_str)
		if err != nil {
			return nil, err
		}
		translate, err := j.GetString("translate")
		if err != nil {
			return nil, err
		}
		return &splay{
			ChatMessage: &splay_chatMessage{
				Translate: translate,
				Json:      j,
			},
		}, nil
	case constants.SKeepAlive:
		ka := binary.BigEndian.Uint32(buffer.Bytes())
		return &splay{KeepAlive: util.Int(int(ka))}, nil
	// 以下暂时不实现
	case constants.SEntityHeadLook:
		return &splay{EntityHeadLook: &struct{}{}}, nil
	case constants.SEntityRelativeMove:
		return &splay{EntityRelativeMove: &struct{}{}}, nil
	case constants.SEntityVelocity:
		return &splay{EntityVelocity: &struct{}{}}, nil
	case constants.SMapChunkBulk:
		return &splay{MapChunkBulk: &struct{}{}}, nil
	case constants.STimeUpdate:
		return &splay{TimeUpdate: &struct{}{}}, nil
	case constants.SSpawnMob:
		return &splay{SpawnMob: &struct{}{}}, nil
	case constants.SEntityLookandRelativeMove:
		return &splay{EntityLookandRelativeMove: &struct{}{}}, nil
	case constants.SEntityMetadata:
		return &splay{EntityMetadata: &struct{}{}}, nil
	case constants.SEntityProperties:
		return &splay{EntityProperties: &struct{}{}}, nil
	case constants.SSpawnObject:
		return &splay{SpawnObject: &struct{}{}}, nil
	case constants.SHeldItemChange:
		return &splay{HeldItemChange: &struct{}{}}, nil
	case constants.SPlayerAbilities:
		return &splay{PlayerAbilities: &struct{}{}}, nil
	case constants.SPlayerPositionAndLook:
		return &splay{PlayerPositionAndLook: &struct{}{}}, nil
	case constants.SPlayerListItem:
		return &splay{PlayerListItem: &struct{}{}}, nil
	case constants.SSetSlot:
		return &splay{SetSlot: &struct{}{}}, nil
	case constants.SSoundEffect:
		return &splay{SoundEffect: &struct{}{}}, nil
	case constants.SSpawnPlayer:
		return &splay{SpawnPlayer: &struct{}{}}, nil
	case constants.SSpawnPosition:
		return &splay{SpawnPosition: &struct{}{}}, nil
	case constants.SStatistics:
		return &splay{Statistics: &struct{}{}}, nil
	case constants.SWindowItems:
		return &splay{WindowItems: &struct{}{}}, nil
	case constants.SEntityLook:
		return &splay{EntityLook: &struct{}{}}, nil
	case constants.SEntityEquipment:
		return &splay{EntityEquipment: &struct{}{}}, nil
	case constants.SJoinGame:
		return &splay{JoinGame: &struct{}{}}, nil
	case constants.SEntityTeleport:
		return &splay{EntityTeleport: &struct{}{}}, nil
	case constants.SChangeGameState:
		return &splay{ChangeGameState: &struct{}{}}, nil
	case constants.SMultiBlockChange:
		return &splay{MultiBlockChange: &struct{}{}}, nil
	case constants.SDestroyEntities:
		return &splay{DestroyEntities: &struct{}{}}, nil
	case constants.SBlockChange:
		return &splay{BlockChange: &struct{}{}}, nil
	case constants.SEntityStatus:
		return &splay{EntityStatus: &struct{}{}}, nil
	case constants.SSpawnExperienceOrb:
		return &splay{SpawnExperienceOrb: &struct{}{}}, nil
	case constants.SCollectItem:
		return &splay{CollectItem: &struct{}{}}, nil
	case constants.SAnimation:
		return &splay{Animation: &struct{}{}}, nil
	}
}

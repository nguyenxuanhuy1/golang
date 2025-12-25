package websocket

import (
	"encoding/binary"
	"errors"
)

// EVENT CODES
const (
	EventPlayerJoin = 1
	EventMove       = 2 // Server gửi vị trí thật xuống FE
	EventShoot      = 3
	EventHit        = 4
	EventInput      = 5 // FE gửi input lên server
)

// ============== PLAYER JOIN ==============
// [event][userID(2 byte)]
func EncodePlayerJoin(userID uint16) []byte {
	buf := make([]byte, 3)
	buf[0] = EventPlayerJoin
	binary.BigEndian.PutUint16(buf[1:], userID)
	return buf
}

// =============== MOVE PACKET (server -> FE) ===============
// [event][userID][x][y]
type MovePacket struct {
	UserID uint16
	X      uint16
	Y      uint16
}

func EncodeMove(userID, x, y uint16) []byte {
	buf := make([]byte, 7)
	buf[0] = EventMove
	binary.BigEndian.PutUint16(buf[1:], userID)
	binary.BigEndian.PutUint16(buf[3:], x)
	binary.BigEndian.PutUint16(buf[5:], y)
	return buf
}

func DecodeMove(b []byte) (MovePacket, error) {
	if len(b) < 7 {
		return MovePacket{}, errors.New("invalid move packet")
	}
	return MovePacket{
		UserID: binary.BigEndian.Uint16(b[1:3]),
		X:      binary.BigEndian.Uint16(b[3:5]),
		Y:      binary.BigEndian.Uint16(b[5:7]),
	}, nil
}

// =============== INPUT PACKET (FE -> server) ===============
// [event=5][mask]
//
// mask: bitmask 1 byte
// 1 = up
// 2 = down
// 4 = left
// 8 = right
func DecodeInput(b []byte) (uint8, error) {
	if len(b) < 2 {
		return 0, errors.New("invalid input packet")
	}
	return b[1], nil
}

// =============== SHOOT PACKET ===============
// [event][userID][angle][power]
type ShootPacket struct {
	UserID uint16
	Angle  uint16
	Power  uint16
}

func DecodeShoot(b []byte) (ShootPacket, error) {
	if len(b) < 7 {
		return ShootPacket{}, errors.New("invalid shoot packet")
	}
	return ShootPacket{
		UserID: binary.BigEndian.Uint16(b[1:3]),
		Angle:  binary.BigEndian.Uint16(b[3:5]),
		Power:  binary.BigEndian.Uint16(b[5:7]),
	}, nil
}

// =============== HIT PACKET ===============
// [event][shooter][target][damage]
type HitPacket struct {
	Shooter uint16
	Target  uint16
	Damage  uint16
}

func EncodeHit(shooter, target, damage uint16) []byte {
	buf := make([]byte, 7)
	buf[0] = EventHit
	binary.BigEndian.PutUint16(buf[1:], shooter)
	binary.BigEndian.PutUint16(buf[3:], target)
	binary.BigEndian.PutUint16(buf[5:], damage)
	return buf
}

func DecodeHit(b []byte) (HitPacket, error) {
	if len(b) < 7 {
		return HitPacket{}, errors.New("invalid hit packet")
	}
	return HitPacket{
		Shooter: binary.BigEndian.Uint16(b[1:3]),
		Target:  binary.BigEndian.Uint16(b[3:5]),
		Damage:  binary.BigEndian.Uint16(b[5:7]),
	}, nil
}

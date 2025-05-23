package main

import (
	"bytes"
	"encoding/binary"
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type Option func(*GamePerson)

func WithName(name string) func(*GamePerson) {
	return func(p *GamePerson) {
		n := copy(p.name[:], name)
		for i := n; i < len(p.name); i++ {
			p.name[i] = 0
		}
	}
}

func WithCoordinates(x, y, z int) func(*GamePerson) {
	return func(p *GamePerson) {
		buf := new(bytes.Buffer)
		binary.Write(buf, binary.LittleEndian, int32(x))
		binary.Write(buf, binary.LittleEndian, int32(y))
		binary.Write(buf, binary.LittleEndian, int32(z))
		p.position = [12]byte(buf.Bytes())
	}
}

func WithGold(gold int) func(*GamePerson) {
	return func(p *GamePerson) {
		*(*uint32)(unsafe.Pointer(&p.gold)) = uint32(gold)
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(p *GamePerson) {
		*(*uint16)(unsafe.Pointer(&p.manaHealth[0])) = uint16(mana)
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(p *GamePerson) {
		*(*uint16)(unsafe.Pointer(&p.manaHealth[0])) = uint16(health)
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(p *GamePerson) {
		p.respect[0] = uint8(respect)
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(p *GamePerson) {
		p.strength[0] = uint8(strength)
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(p *GamePerson) {
		p.expLevel[0] = (p.expLevel[0] & 0xF0) | uint8(experience)&15
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(p *GamePerson) {
		p.expLevel[0] = (p.expLevel[0] & 15) | (uint8(level) << 4)
	}
}

func WithHouse() func(*GamePerson) {
	return func(p *GamePerson) {
		p.flags[0] |= 1 << 0
	}
}

func WithGun() func(*GamePerson) {
	return func(p *GamePerson) {
		p.flags[0] |= 1 << 1
	}
}

func WithFamily() func(*GamePerson) {
	return func(p *GamePerson) {
		p.flags[0] |= 1 << 2
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(p *GamePerson) {
		p.flags[0] = (p.flags[0] &^ 24) | (byte(personType) << 3)
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

type GamePerson struct {
	name       [42]byte
	position   [12]byte
	gold       [4]byte
	manaHealth [2]byte
	respect    [1]byte
	strength   [1]byte
	expLevel   [1]byte
	flags      [1]byte // hasHouse(1), hasGun(1), hasFamily(1), type(2)
}

func NewGamePerson(options ...Option) *GamePerson {
	person := &GamePerson{}
	for _, opt := range options {
		opt(person)
	}
	return person
}

func (p *GamePerson) Name() string {
	end := 0
	for end < len(p.name) && p.name[end] != 0 {
		end++
	}
	return string(p.name[:end])
}

func (p *GamePerson) X() int {
	return int(*(*int32)(unsafe.Pointer(&p.position[0])))
}

func (p *GamePerson) Y() int {
	return int(*(*int32)(unsafe.Pointer(&p.position[4])))
}

func (p *GamePerson) Z() int {
	return int(*(*int32)(unsafe.Pointer(&p.position[8])))
}

func (p *GamePerson) Gold() int {
	return int(*(*uint32)(unsafe.Pointer(&p.gold)))
}

func (p *GamePerson) Mana() int {
	return int(*(*uint16)(unsafe.Pointer(&p.manaHealth[0])))
}

func (p *GamePerson) Health() int {
	return int(*(*uint16)(unsafe.Pointer(&p.manaHealth[0])))
}

func (p *GamePerson) Respect() int {
	return int(p.respect[0])
}

func (p *GamePerson) Strength() int {
	return int(p.strength[0])
}

func (p *GamePerson) Experience() int {
	return int(p.expLevel[0] & 15)
}

func (p *GamePerson) Level() int {
	return int((p.expLevel[0] >> 4) & 15)
}

func (p *GamePerson) HasHouse() bool {
	return p.flags[0]&(1<<0) != 0
}

func (p *GamePerson) HasGun() bool {
	return p.flags[0]&(1<<1) != 0
}

func (p *GamePerson) HasFamilty() bool {
	return p.flags[0]&(1<<2) != 0
}

func (p *GamePerson) Type() int {
	return int((p.flags[0] & 24) >> 3)
}

func TestGamePerson(t *testing.T) {
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))

	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = BuilderGamePersonType
	const gold = math.MaxInt32
	const mana = 1000
	const health = 1000
	const respect = 10
	const strength = 10
	const experience = 10
	const level = 10

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(),
		WithFamily(),
		WithType(personType),
	}

	person := NewGamePerson(options...)
	assert.Equal(t, name, person.Name())
	assert.Equal(t, x, person.X())
	assert.Equal(t, y, person.Y())
	assert.Equal(t, z, person.Z())
	assert.Equal(t, gold, person.Gold())
	assert.Equal(t, mana, person.Mana())
	assert.Equal(t, health, person.Health())
	assert.Equal(t, respect, person.Respect())
	assert.Equal(t, strength, person.Strength())
	assert.Equal(t, experience, person.Experience())
	assert.Equal(t, level, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamilty())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())
}

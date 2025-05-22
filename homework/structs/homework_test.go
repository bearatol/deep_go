package main

import (
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type Option func(*GamePerson)

func WithName(name string) func(*GamePerson) {
	return func(person *GamePerson) {
		n := copy(person.name[:], name)
		for i := n; i < len(person.name); i++ {
			person.name[i] = 0
		}
	}
}

func WithCoordinates(x, y, z int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.x = int32(x)
		person.y = int32(y)
		person.z = int32(z)
	}
}

func WithGold(gold int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.gold = uint32(gold)
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.mana = uint16(mana)
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.health = uint16(health)
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.stats = (person.stats &^ respectMask) | (uint32(respect) << respectShift)
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.stats = (person.stats &^ strengthMask) | (uint32(strength) << strengthShift)
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.stats = (person.stats &^ experienceMask) | (uint32(experience) << experienceShift)
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.stats = (person.stats &^ levelMask) | (uint32(level) << levelShift)
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		person.isHouseGunFamily |= hasHouseBit
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		person.isHouseGunFamily |= hasGunBit
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		person.isHouseGunFamily |= hasFamilyBit
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.personType = uint8(personType)
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

const (
	experienceShift = 24
	respectShift    = 16
	strengthShift   = 8
	levelShift      = 0

	experienceMask = 0xFF << experienceShift
	respectMask    = 0xFF << respectShift
	strengthMask   = 0xFF << strengthShift
	levelMask      = 0xFF << levelShift

	hasHouseBit  = 1 << 0
	hasGunBit    = 1 << 1
	hasFamilyBit = 1 << 2
)

type GamePerson struct {
	name [42]byte

	x int32
	y int32
	z int32

	gold  uint32
	stats uint32

	mana   uint16
	health uint16

	personType       uint8
	isHouseGunFamily uint8
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
	return int(p.x)
}

func (p *GamePerson) Y() int {
	return int(p.y)
}

func (p *GamePerson) Z() int {
	return int(p.z)
}

func (p *GamePerson) Gold() int {
	return int(p.gold)
}

func (p *GamePerson) Mana() int {
	return int(p.mana)
}

func (p *GamePerson) Health() int {
	return int(p.health)
}

func (p *GamePerson) Respect() int {
	return int((p.stats & respectMask) >> respectShift)
}

func (p *GamePerson) Strength() int {
	return int((p.stats & strengthMask) >> strengthShift)
}

func (p *GamePerson) Experience() int {
	return int((p.stats & experienceMask) >> experienceShift)
}

func (p *GamePerson) Level() int {
	return int((p.stats & levelMask) >> levelShift)
}

func (p *GamePerson) HasHouse() bool {
	return p.isHouseGunFamily&hasHouseBit != 0
}

func (p *GamePerson) HasGun() bool {
	return p.isHouseGunFamily&hasGunBit != 0
}

func (p *GamePerson) HasFamilty() bool {
	return p.isHouseGunFamily&hasFamilyBit != 0
}

func (p *GamePerson) Type() int {
	return int(p.personType)
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

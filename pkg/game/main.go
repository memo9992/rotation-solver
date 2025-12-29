package game

import (
	"fmt"
	"time"

	"github.com/memo9992/rotation-solver/pkg/utils"
)

type DamageType string

const (
	Physical   DamageType = "physical"
	Magical    DamageType = "magical"
	Unaspected DamageType = "unaspected"
)

type BuffAction struct {
	Action BuffActionType
	Target Buff
}

type BuffActionFunc func(t time.Time, buffs *Buffs) *BuffAction

type Ability struct {
	ID          string
	Name        string
	Description string
	Cooldown    time.Duration
	CastTime    time.Duration
	Damage      int64
	DamageType  DamageType
	Tags        *utils.Set[string]
	OnActivate  BuffActionFunc
	AfterBuffs  func(self Ability, buffs *Buffs) Ability
}

func (a Ability) WithCastTime(d time.Duration) Ability {
	a.CastTime = d
	return a
}

func (a Ability) WithOnActivate(f BuffActionFunc) Ability {
	a.OnActivate = f
	return a
}

func (a Ability) WithAfterBuffs(f func(self Ability, buffs *Buffs) Ability) Ability {
	a.AfterBuffs = f
	return a
}

type Buff struct {
	ID            string
	Name          string
	Description   string
	TimeStart     time.Time
	Duration      time.Duration
	ActualTimeEnd time.Time
	Stacks        int
	Effects       map[BuffEffectType]int64
}

func (b Buff) IsZero() bool {
	return b.ID == ""
}

type Buffs struct {
	buffs []Buff
}

func (b *Buffs) Apply(action *BuffAction) {
	if action == nil {
		return
	}
	switch action.Action {
	case "add":
		b.buffs = append(b.buffs, action.Target)
	}
}

func (b *Buffs) Update(t time.Time) {
	out := []Buff{}

	for _, currentBuff := range b.buffs {
		if currentBuff.TimeStart.Add(currentBuff.Duration).After(t) {
			out = append(out, currentBuff)
		}
	}

	b.buffs = out
}

func (b *Buffs) Get(id string) Buff {
	for _, buff := range b.buffs {
		if buff.ID == id {
			return buff
		}
	}

	return Buff{}
}

func (b *Buffs) State() []Buff {
	out := []Buff{}

	for _, buff := range b.buffs {
		out = append(out, buff)
	}

	return out
}

func (a Ability) IsGCD() bool {
	return a.Tags.Contains("gcd")
}

func (a Ability) IsOGCD() bool {
	return a.Tags.Contains("ogcd")
}

type NewAbilityOptions struct {
	ID          string
	Name        string
	Description string
	Cooldown    time.Duration
	CastTime    time.Duration
	Damage      int64
	DamageType  DamageType
	Tags        *utils.Set[string]
	OnActivate  BuffActionFunc
	AfterBuffs  func(Ability, *Buffs) Ability
}

func NewAbility(opts NewAbilityOptions) (Ability, error) {
	if opts.ID == "" {
		return Ability{}, fmt.Errorf("ability ID cannot be empty")
	}
	if opts.Name == "" {
		return Ability{}, fmt.Errorf("ability Name cannot be empty")
	}
	if opts.DamageType != Physical && opts.DamageType != Magical && opts.DamageType != Unaspected {
		return Ability{}, fmt.Errorf("invalid damage type: %s", opts.DamageType)
	}
	if opts.Tags == nil {
		opts.Tags = utils.NewSet[string]()
	}
	if opts.OnActivate == nil {
		opts.OnActivate = func(t time.Time, buffs *Buffs) *BuffAction {
			return nil
		}
	}
	if opts.AfterBuffs == nil {
		opts.AfterBuffs = func(a Ability, b *Buffs) Ability {
			return a
		}
	}

	return Ability{
		ID:          opts.ID,
		Name:        opts.Name,
		Description: opts.Description,
		Cooldown:    opts.Cooldown,
		CastTime:    opts.CastTime,
		Damage:      opts.Damage,
		DamageType:  opts.DamageType,
		Tags:        opts.Tags,
		OnActivate:  opts.OnActivate,
		AfterBuffs:  opts.AfterBuffs,
	}, nil
}

func MustNewAbility(opts NewAbilityOptions) Ability {
	ability, err := NewAbility(opts)
	if err != nil {
		panic(err)
	}
	return ability
}

func MustNewGCD(id string, name string, damage int64, damageType DamageType) Ability {
	s := utils.NewSet("gcd")
	return MustNewAbility(NewAbilityOptions{
		ID:         id,
		Name:       name,
		CastTime:   0,
		Cooldown:   0,
		Damage:     damage,
		DamageType: damageType,
		Tags:       s,
	})
}

func MustNewOGCD(id string, name string, cooldown time.Duration, damage int64, damageType DamageType) Ability {
	s := utils.NewSet("ogcd")
	return MustNewAbility(NewAbilityOptions{
		ID:         id,
		Name:       name,
		Cooldown:   cooldown,
		Damage:     damage,
		DamageType: damageType,
		Tags:       s,
	})
}

func MustNewDamageBuff(id string, name string, duration time.Duration, cooldown time.Duration, multiplier int64) Ability {
	a := MustNewOGCD(id, name, cooldown, 0, Unaspected)
	a.OnActivate = func(t time.Time, buffs *Buffs) *BuffAction {
		return &BuffAction{
			Action: "add",
			Target: Buff{
				ID:        id,
				Name:      name,
				TimeStart: t,
				Duration:  duration,
				Stacks:    1,
				Effects: map[BuffEffectType]int64{
					BuffEffectDamageMultiplier: multiplier,
				},
			},
		}
	}

	return a
}

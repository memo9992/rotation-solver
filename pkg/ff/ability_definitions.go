package ff

import (
	"fmt"
	"time"

	"github.com/memo9992/rotation-solver/pkg/game"
)

var (
	FastBlade = game.MustNewGCD(
		"fastBlade",
		"Fast Blade",
		120,
		game.Physical,
	).WithOnActivate(func(t time.Time, buffs *game.Buffs) *game.BuffAction {
		return &game.BuffAction{
			Action: game.BuffActionTypeAdd,
			Target: game.Buff{
				ID:        "fastBladeCombo",
				TimeStart: t,
				Duration:  time.Second * 30,
				Effects: map[game.BuffEffectType]int64{
					game.BuffEffectComboPotency: 330,
				},
			},
		}
	})
	RiotBlade = game.MustNewGCD(
		"riotBlade",
		"Riot Blade",
		170,
		game.Physical,
	).WithAfterBuffs(func(self game.Ability, buffs *game.Buffs) game.Ability {
		if combo := buffs.Get("fastBladeCombo"); !combo.IsZero() {
			fmt.Println("Combo active")
			self.Damage = 330
			buffs.Apply(&game.BuffAction{
				Action: game.BuffActionTypeRemove,
				Target: combo,
			})
		}
		return self
	})
	GoringBlade = game.MustNewGCD(
		"DD2",
		"Goring Blade",
		700,
		game.Physical,
	)
	Requiescat = game.MustNewOGCD(
		"requiscat",
		"Requiscat",
		time.Minute,
		320,
		game.Unaspected,
	).WithOnActivate(func(t time.Time, buffs *game.Buffs) *game.BuffAction {
		return &game.BuffAction{
			Action: game.BuffActionTypeAdd,
			Target: game.Buff{
				ID:        "requiscat",
				Stacks:    5,
				TimeStart: t,
				Duration:  30,
				Effects:   map[game.BuffEffectType]int64{
					// IDK
				},
			},
		}
	})
	FightOrFlight = game.MustNewDamageBuff(
		"fightOrFlight",
		"Fight or Flight",
		time.Second*5,
		time.Minute,
		20,
	)
	HolySpirit = game.MustNewGCD(
		"holySpirit",
		"Holy Spirit",
		400,
		game.Magical,
	).WithCastTime(time.Millisecond * 4000)
)

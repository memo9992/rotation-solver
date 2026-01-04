package game

import (
	"fmt"
	"time"
)

const (

	// Damage modifiers

	DECREASE string = "decrease"
	INCREASE string = "increase"
	ADDITIVE string = "additive"

	// Condition actions

	ADD    string = "add"
	REMOVE string = "remove"
)

// Condition is a buff or debuff on the player. It can provide global damage modifiers
// or exist as a pre-condition for other ability modifiers (e.g. combos).
type Condition struct {
	ID        string
	Start     time.Time
	Duration  time.Duration
	Modifiers []DamageModifier
}

// Conditions is the current state of buffs/debuffs on the player.
type Conditions map[string]Condition

// WithAction applies a ConditionAction to the player state.
func (c Conditions) WithAction(a ConditionAction) {
	switch a.Action {
	case ADD:
		c[a.ID] = a.Condition
	case REMOVE:
		delete(c, a.ID)
	default:
		panic(fmt.Sprintf("Unknown condition action %v", a.Action))
	}
}

// ConditionAction is an action to be applied to the player state. This will either add a new
// Condition or modify or remove an existing Condition.
type ConditionAction struct {
	Action string
	Condition
}

// Value is a numerical value with the source from which it originates from. This is used to track
// damage modifiers from the abilities/buffs that caused them.
type Value struct {
	Source string
	Value  int
}

// DamageInstance is a discrete set of damage that was caused by an ability. It includes slices for
// base (additive) damage with the total increased modifiers.
type DamageInstance struct {
	Source    string
	Base      []Value
	Increased []Value
}

// Create a new DamageInstance with the given modifier.
func (d DamageInstance) WithModifier(di DamageModifier) DamageInstance {
	switch di.Action {
	case ADDITIVE:
		d.Base = append(d.Base, Value{
			Source: di.Source,
			Value:  di.Value,
		})
	case INCREASE:
		d.Increased = append(d.Increased, Value{
			Source: di.Source,
			Value:  di.Value,
		})
	default:
		panic(fmt.Sprintf("Unknown damage modifier action: %v", di.Action))
	}

	return d
}

// Calculate the total damage done by this instance.
func (d DamageInstance) Calculate() int {
	base := 0
	for _, b := range d.Base {
		base += b.Value
	}

	increased := 100
	for _, i := range d.Increased {
		increased += i.Value
	}

	return base * increased / 100
}

// DamageModifier affects an Ability before it creates a DamageInstance.
type DamageModifier struct {
	Source string
	Action string
	Value  int

	// ConditionID that must exists on the player to apply this modifier.
	RequiredConditionID string
}

// Ability creates damage instances, condition actions or both.
type Ability struct {
	ID                 string
	DoesDamage         bool
	Base               int
	DamageModifiers    []DamageModifier
	ConditionModifiers []ConditionAction
}

// Activate the Ability to create damage instances or condition actions based on the player's current
// condition state.
func (a Ability) Activate(conditions Conditions) ([]DamageInstance, []ConditionAction) {
	var diOut []DamageInstance
	if a.DoesDamage {
		di := DamageInstance{
			Source: a.ID,
			Base: []Value{
				{
					Source: a.ID,
					Value:  a.Base,
				},
			},
		}

		for _, d := range a.DamageModifiers {
			if d.RequiredConditionID != "" {
				if _, ok := conditions[d.RequiredConditionID]; !ok {
					continue
				}
			}

			di = di.WithModifier(d)
		}

		for _, c := range conditions {
			for _, m := range c.Modifiers {
				di = di.WithModifier(m)
			}
		}

		diOut = append(diOut, di)
	}

	return diOut, a.ConditionModifiers
}

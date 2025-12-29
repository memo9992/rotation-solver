package game

type BuffActionType string

const (
	BuffActionTypeAdd    BuffActionType = "add"
	BuffActionTypeRemove BuffActionType = "remove"
)

type BuffEffectType string

const (
	BuffEffectDamageMultiplier BuffEffectType = "damageMultiplier"
	BuffEffectComboPotency     BuffEffectType = "comboPotency"
)

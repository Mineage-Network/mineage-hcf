package custom

// TripwireHook ...
type TripwireHook struct{}

// EncodeItem ...
func (t TripwireHook) EncodeItem() (name string, meta int16) {
	return "minecraft:tripwire_hook", 0
}

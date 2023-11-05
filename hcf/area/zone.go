package area

import "github.com/mineage-network/mineage-hcf/hcf/util"

// Zone ...
type Zone struct {
	name string
	area util.AreaVec2
}

// NewZone ...
func NewZone(name string, area util.AreaVec2) Zone {
	return Zone{name: name, area: area}
}

// Name ...
func (z Zone) Name() string {
	return z.name
}

// Area ...
func (z Zone) Area() util.AreaVec2 {
	return z.area
}

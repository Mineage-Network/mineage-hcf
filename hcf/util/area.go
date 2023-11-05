package util

import (
	"github.com/go-gl/mathgl/mgl64"
	"math"
)

// AreaVec2 represents a 2D area.
type AreaVec2 struct {
	// minX is the minimum X value.`
	minX,
	// maxX is the maximum X value.
	maxX,
	// minY is the minimum Y value.
	minY,
	// maxY is the maximum Y value.
	maxY float64
}

// Max returns a mgl64.Vec2 with the maximum X and Y values.
func (v AreaVec2) Max() mgl64.Vec2 { return mgl64.Vec2{v.maxX, v.maxY} }

// Min returns a mgl64.Vec2 with the minimum X and Y values.
func (v AreaVec2) Min() mgl64.Vec2 { return mgl64.Vec2{v.minX, v.minY} }

// NewAreaVec2 returns a new AreaVec2 area with the minimum and maximum X and Y values.
func NewAreaVec2(b1, b2 mgl64.Vec2) AreaVec2 {
	return AreaVec2{
		minX: minBound(b1.X(), b2.X()),
		maxX: maxBound(b1.X(), b2.X()),

		minY: minBound(b1.Y(), b2.Y()),
		maxY: maxBound(b1.Y(), b2.Y()),
	}
}

// Vec2Within returns true if the given mgl64.Vec2 is within the area.
func (v AreaVec2) Vec2Within(vec mgl64.Vec2) bool {
	return vec.X() > v.minX && vec.X() < v.maxX && vec.Y() > v.minY && vec.Y() < v.maxY
}

// Vec3WithinXZ returns true if the given mgl64.Vec3 is within the area.
func (v AreaVec2) Vec3WithinXZ(vec mgl64.Vec3) bool {
	return vec.X() > v.minX && vec.X() < v.maxX && vec.Z() > v.minY && vec.Z() < v.maxY
}

// Vec2WithinOrEqual returns true if the given mgl64.Vec2 is within or equal to the bounds of the area.
func (v AreaVec2) Vec2WithinOrEqual(vec mgl64.Vec2) bool {
	return vec.X() >= v.minX && vec.X() <= v.maxX && vec.Y() >= v.minY && vec.Y() <= v.maxY
}

// Vec3WithinOrEqualXZ returns true if the given mgl64.Vec3 is within or equal to the bounds of the area.
func (v AreaVec2) Vec3WithinOrEqualXZ(vec mgl64.Vec3) bool {
	return vec.X() >= v.minX && vec.X() <= v.maxX && vec.Z() >= v.minY && vec.Z() <= v.maxY
}

// Vec2WithinOrEqualFloor returns true if the given mgl64.Vec2 is within or equal to the bounds of the area.
func (v AreaVec2) Vec2WithinOrEqualFloor(vec mgl64.Vec2) bool {
	vec = mgl64.Vec2{math.Floor(vec.X()), math.Floor(vec.Y())}
	return vec.X() >= v.minX && vec.X() <= v.maxX && vec.Y() >= v.minY && vec.Y() <= v.maxY
}

// Vec3WithinOrEqualFloorXZ returns true if the given mgl64.Vec3 is within or equal to the bounds of the area.
func (v AreaVec2) Vec3WithinOrEqualFloorXZ(vec mgl64.Vec3) bool {
	vec = mgl64.Vec3{math.Floor(vec.X()), vec.Y(), math.Floor(vec.Z())}
	return vec.X() >= v.minX && vec.X() <= v.maxX && vec.Z() >= v.minY && vec.Z() <= v.maxY
}

// AreaVec3 represents a 3D area.
type AreaVec3 struct {
	// minX is the minimum X value.
	minX,
	// maxX is the maximum X value.
	maxX,
	// minY is the minimum Y value.
	minY,
	// maxY is the maximum Y value.
	maxY,
	// minZ is the minimum Z value.
	minZ,
	// maxZ is the maximum Z value.
	maxZ float64
}

// Max returns a mgl64.Vec3 with the maximum X, Y, and Z values.
func (v AreaVec3) Max() mgl64.Vec3 { return mgl64.Vec3{v.maxX, v.maxY, v.maxZ} }

// Min returns a mgl64.Vec3 with the minimum X, Y, and Z values.
func (v AreaVec3) Min() mgl64.Vec3 { return mgl64.Vec3{v.minX, v.minY, v.minZ} }

// NewAreaVec3 returns a new AreaVec3 area with the minimum and maximum X, Y, and Z values.
func NewAreaVec3(b1, b2 mgl64.Vec3) AreaVec3 {
	return AreaVec3{
		minX: minBound(b1.X(), b2.X()),
		maxX: maxBound(b1.X(), b2.X()),

		minY: minBound(b1.Y(), b2.Y()),
		maxY: maxBound(b1.Y(), b2.Y()),

		minZ: minBound(b1.Z(), b2.Z()),
		maxZ: maxBound(b1.Z(), b2.Z()),
	}
}

// Vec3Within returns true if the given mgl64.Vec3 is within the area.
func (v AreaVec3) Vec3Within(vec mgl64.Vec3) bool {
	return vec.X() > v.minX && vec.X() < v.maxX && vec.Y() > v.minY && vec.Y() < v.maxY && vec.Z() > v.minZ && vec.Z() < v.maxZ
}

// Vec3WithinOrEqual returns true if the given mgl64.Vec3 is within or equal to the bounds of the area.
func (v AreaVec3) Vec3WithinOrEqual(vec mgl64.Vec3) bool {
	return vec.X() >= v.minX && vec.X() <= v.maxX && vec.Y() >= v.minY && vec.Y() <= v.maxY && vec.Z() >= v.minZ && vec.Z() <= v.maxZ
}

// maxBound returns the maximum of two numbers.
func maxBound(b1, b2 float64) float64 {
	if b1 > b2 {
		return b1
	}
	return b2
}

// minBound returns the minimum of two numbers.
func minBound(b1, b2 float64) float64 {
	if b1 < b2 {
		return b1
	}
	return b2
}

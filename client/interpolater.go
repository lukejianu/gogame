package client

import (
	"time"

	"github.com/lukejianu/gogame/common"
)

type TimestampedClientGameState struct {
	Timestamp time.Time // The time at which this CGS was received.
	common.ClientGameState
}

func TimestampCgs(cgs common.ClientGameState) TimestampedClientGameState {
	return TimestampedClientGameState{
		Timestamp:       time.Now(),
		ClientGameState: cgs,
	}
}

type Interpolater interface {
	Interpolate(t time.Time) (common.ClientGameState, bool)
}

func NewGameStateInterpolater(lastGs, currGs TimestampedClientGameState) Interpolater {
	return interpolater{
		velocity: computeVelocities(lastGs, currGs),
		lastGs:   lastGs,
		currGs:   currGs,
	}
}

func computeVelocities(lastGs, currGs TimestampedClientGameState) map[common.ID]float64 {
	velocities := map[common.ID]float64{}
	dt := max(currGs.Timestamp.Sub(lastGs.Timestamp), 1)
	for id, _ := range currGs.Others {
		dx := computeDx(id, lastGs, currGs)
		velocities[id] = dx / float64(dt.Milliseconds())
	}
	return velocities
}

func computeDx(id common.ID, lastGs, currGs TimestampedClientGameState) float64 {
	x1 := currGs.Others[id]
	x2, ok := lastGs.Others[id]
	if !ok {
		return 0 // The dx is 0 because the ID is new.
	}
	return float64(x1 - x2)
}

type interpolater struct {
	velocity map[common.ID]float64
	lastGs   TimestampedClientGameState
	currGs   TimestampedClientGameState
}

func (i interpolater) Interpolate(t time.Time) (common.ClientGameState, bool) {
	// Recall that from time [t, t + updateGap], we simulate
	// the movement from [t - updateGap, t].
	updateGap := i.currGs.Timestamp.Sub(i.lastGs.Timestamp)
	if t.After(i.currGs.Timestamp.Add(updateGap)) {
		return common.ClientGameState{}, true
	}

	cgs := common.ClientGameState{
		Others: map[common.ID]int{},
	}
	msPassed := float64(t.Sub(i.currGs.Timestamp).Milliseconds())
	for id, x := range i.lastGs.Others {
		v := i.velocity[id]
		cgs.Others[id] = int(float64(x) + (msPassed * v))
	}

	return cgs, false
}

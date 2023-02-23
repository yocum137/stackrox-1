package postgres

import (
	"io"
	"math/rand"

	"github.com/stackrox/rox/pkg/buildinfo"
	"github.com/stackrox/rox/pkg/env"
)

var (
	postgresChaosEnabled = env.PostgresChaosPercent.IntegerSetting() > 0 && !buildinfo.ReleaseBuild
	chaosPercent         = env.PostgresChaosPercent.IntegerSetting()
)

func getChaosError() error {
	if !postgresChaosEnabled {
		return nil
	}
	rate := float64(chaosPercent) / 100
	if rand.Float64() > rate {
		return nil
	}
	return io.EOF
}

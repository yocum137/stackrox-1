package phonehome

import (
	"context"
	"time"

	"github.com/stackrox/rox/pkg/concurrency"
	"github.com/stackrox/rox/pkg/sync"
)

var (
	gathererInstance *gatherer
	onceGatherer     sync.Once
)

// Time period for static data gathering from data sources.
const period = 5 * time.Minute

type gatherer struct {
	telemeter   Telemeter
	period      time.Duration
	stopSig     concurrency.Signal
	ctx         context.Context
	mu          sync.Mutex
	gatherFuncs []GatherFunc
}

// Gatherer interface for interacting with telemetry gatherer.
type Gatherer interface {
	Start()
	Stop()
	AddGatherer(GatherFunc)
}

func (g *gatherer) reset() {
	g.stopSig.Reset()
	g.ctx, _ = concurrency.DependentContext(context.Background(), &g.stopSig)
}

func newGatherer(t Telemeter, p time.Duration) *gatherer {
	return &gatherer{
		telemeter: t,
		period:    p,
	}
}

// GathererSingleton returns the telemetry gatherer instance.
func (cfg *Config) GathererSingleton() Gatherer {
	if Enabled() {
		onceGatherer.Do(func() {
			gathererInstance = newGatherer(cfg.TelemeterSingleton(), period)
		})
	}
	return gathererInstance
}

func (g *gatherer) collect() map[string]any {
	var result map[string]any
	for i, f := range g.gatherFuncs {
		props, err := f(g.ctx)
		if err != nil {
			log.Errorf("gatherer %d failure: %v", i, err)
		}
		if props != nil && result == nil {
			result = make(map[string]any, len(props))
		}
		for k, v := range props {
			result[k] = v
		}
	}
	return result
}

func (g *gatherer) loop() {
	ticker := time.NewTicker(g.period)
	for !g.stopSig.IsDone() {
		select {
		case <-ticker.C:
			go func() {
				g.telemeter.Identify(g.collect())
			}()
		case <-g.stopSig.Done():
			ticker.Stop()
			return
		}
	}
}

func (g *gatherer) Start() {
	if g == nil {
		return
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.stopSig.IsDone() {
		g.reset()
		go g.loop()
	}
}

func (g *gatherer) Stop() {
	if g == nil {
		return
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	g.stopSig.Signal()
}

func (g *gatherer) AddGatherer(f GatherFunc) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.gatherFuncs = append(g.gatherFuncs, f)
}

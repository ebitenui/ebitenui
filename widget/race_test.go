//go:build race

package widget

// raceEnabled is whether the tests run under the race detector, where
// sync.Pool drops items at random, so Ebitengine's allocations per draw
// vary.
const raceEnabled = true

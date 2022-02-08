package booleanpolicy

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/open-policy-agent/opa/rego"
	"github.com/stackrox/rox/pkg/fixtures"
	"github.com/stretchr/testify/require"
)

func BenchmarkOPA(b *testing.B) {
	module := `
package policy.test

test[msg] {
	some i
	input.containers[i].name != ""
	msg := {"name": [input.containers[i].name]}
}
`

	testDeployment := fixtures.GetDeployment()
	res, err := rego.New(
		rego.Query("out = data.policy.test.test"),
		rego.Module("test.policy", module),
		rego.Input(testDeployment),
	).Eval(context.Background())
	_, _ = json.MarshalIndent(res, " ", "   ")
	fmt.Printf("OUT IS %#v\n\n", res)
	return
	q, err := rego.New(
		rego.Query("out = data.policy.test.test"),
		rego.Module("test.policy", module),
		rego.Input(testDeployment),
	).PrepareForEval(context.Background())
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := q.Eval(context.Background())
		require.NoError(b, err)
	}
}

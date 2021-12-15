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
	input.containers[i].security_context.privileged
	some j
	regex.match(".*vol.*", input.containers[i].volumes[j].source)
	msg := sprintf("container '%v' is privileged and has volume %v with source %v",
	[input.containers[i].name, input.containers[i].volumes[j].name, input.containers[i].volumes[j].source])
}
`


	testDeployment := fixtures.GetDeployment()
	res, err := rego.New(
		rego.Query("out = data.policy.test.test"),
		rego.Module("test.policy", module),
		rego.Input(testDeployment),
	).Eval(context.Background())
	out, _ := json.MarshalIndent(res, " ", "   ")
	fmt.Println(string(out))
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


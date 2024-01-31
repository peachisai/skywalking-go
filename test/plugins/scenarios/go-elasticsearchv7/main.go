package go_elasticsearchv7

import (
	"context"
	_ "github.com/apache/skywalking-go"
)

var testCases []struct {
	caseName string
	fn       testFunc
}

type testFunc func(ctx context.Context) error

func init() {
	append(testCases, struct {
		caseName string
		fn       testFunc
	}{caseName: "testIndex", fn: testIndex})
}

func main() {
	for i := range testCases {

	}
}

func testIndex(ctx context.Context) error {

}

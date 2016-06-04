package main

import (
	"context"
	"testing"

	"github.com/kynrai/bop"
)

func TestHello(t *testing.T) {
	for _, test := range []struct {
		name      string
		req, resp *bop.Message
		ctx       context.Context
		err       error
	}{
		{
			name: "success",
			req:  &bop.Message{},
			resp: &bop.Message{},
			ctx:  context.Background(),
		},
	} {
		t.Logf("testing %q", test.name)

		err := Hello(test.ctx, test.req, test.resp)

		if err != test.err {
			t.Errorf("expected: %s got %s", err, test.err)
		}
	}
}

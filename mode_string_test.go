package frontmatter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMode_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		desc string
		give Mode
		want string
	}{
		{
			desc: "SetMetadata",
			give: SetMetadata,
			want: "SetMetadata",
		},
		{
			desc: "unknown",
			give: 0,
			want: "Mode(0)",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, tt.give.String())
		})
	}
}

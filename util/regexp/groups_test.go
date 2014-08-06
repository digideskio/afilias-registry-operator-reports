package regexp

import (
	"github.com/stretchr/testify/assert"
	goregexp "regexp"
	"testing"
)

func TestThatGroupsAreMapped(t *testing.T) {
	assert := assert.New(t)
	r := goregexp.MustCompile(`^(?P<first>[^_]+)_(?P<second>.+)$`)
	groupMap := MapRegexpGroups(r, "my_string")
	assert.Equal("my", groupMap["first"])
	assert.Equal("string", groupMap["second"])
}

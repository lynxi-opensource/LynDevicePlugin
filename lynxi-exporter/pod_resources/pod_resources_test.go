package podresources

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPodResources_Get(t *testing.T) {
	m, err := New()
	assert.Nil(t, err)
	ret, err := m.Get()
	assert.Nil(t, err)
	fmt.Println(ret)
}

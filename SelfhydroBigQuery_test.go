package bigQueryPublisher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldDeserialiseState(t *testing.T) {
	state := []byte("{}")
	expectedState := SeflhydroState{}
	deserialisedState := DeseraliseState(state)
	assert.Equal(t, deserialisedState, expectedState)
}

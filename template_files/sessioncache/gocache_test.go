package sessioncache

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var memcache *GoCache
var testKey = "testKey"
var testValue = "123"

func TestSetup(t *testing.T) {
	memcache = new(GoCache)
	memcache.New()
	memcache.Set(testKey, testValue, 10)
}

func TestGet(t *testing.T) {

	// Test the value is found and returned correctly
	value, _ := memcache.Get(testKey)
	assert.Equal(t, value, testValue)

	// Test error is returned when value not found
	value, err := memcache.Get("myRubbishKey")
	assert.NotNil(t, err, "Key not found error should be returned")

}

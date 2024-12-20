package pokecache

import (
	"testing"
	"time"
)

func TestCacheAddExists(t *testing.T) {
	k := "https://pokeapi.co/api/v2/location-area?offset=40&limit=20"
	v := make([]byte, 5)
	c := NewCache(5 * time.Second)
	c.Add(k, v)
	_, found := c.entries[k]
	if !found {
		t.Errorf("new entry not found after Add")
		t.Fail()
	}
}

func TestCacheAddValues(t *testing.T) {
	k := "https://pokeapi.co/api/v2/location-area?offset=40&limit=20"
	v := make([]byte, 5)
	for i := 0; i < len(v); i++ {
		v[i] = byte('A') + byte(i)
	}
	c := NewCache(5 * time.Second)
	c.Add(k, v)
	for i := 0; i < len(c.entries[k].val); i++ {
		if c.entries[k].val[i] != v[i] {
			t.Errorf("%v does not match %v", c.entries[k].val[i], v[i])
			t.Fail()
		}
	}
}

func TestCacheGetFound(t *testing.T) {
	k := "https://pokeapi.co/api/v2/location-area?offset=40&limit=20"
	v := make([]byte, 5)
	for i := 0; i < len(v); i++ {
		v[i] = byte('A') + byte(i)
	}
	c := NewCache(5 * time.Second)
	c.Add(k, v)
	_, found := c.Get(k)
	if !found {
		t.Errorf("valid entry not found with Get")
		t.Fail()
	}
}

func TestCacheGetNotFound(t *testing.T) {
	k := "https://pokeapi.co/api/v2/location-area?offset=40&limit=20"
	v := make([]byte, 5)
	for i := 0; i < len(v); i++ {
		v[i] = byte('A') + byte(i)
	}
	c := NewCache(5 * time.Second)
	c.Add(k, v)
	_, found := c.Get("https://example.com/api/v1/resource")
	if found {
		t.Errorf("invalid entry found with Get")
		t.Fail()
	}
}

func TestCacheReapAfterInterval(t *testing.T) {
	k := "https://pokeapi.co/api/v2/location-area?offset=40&limit=20"
	c := NewCache(2 * time.Second)
	c.Add(k, []byte{})
	time.Sleep(3 * time.Second)
	_, found := c.Get(k)
	if found {
		t.Errorf("entry should have been reaped")
		t.Fail()
	}
}

func TestCacheReapBeforeInterval(t *testing.T) {
	k := "https://pokeapi.co/api/v2/location-area?offset=40&limit=20"
	c := NewCache(2 * time.Second)
	c.Add(k, []byte{})
	time.Sleep(1 * time.Second)
	_, found := c.Get(k)
	if !found {
		t.Errorf("entry should have been reaped")
		t.Fail()
	}
}

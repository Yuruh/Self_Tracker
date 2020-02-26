package main

import "testing"

func TestDummy(t *testing.T) {
	res := Dummy()
	if res != 1 {
		t.Errorf("Bad dummy")
	}
}


/*
	An example to run several several tests using the same setup / teardown
 */
func TestFoo(t *testing.T) {
	// <setup code>
	t.Run("A=1", func(t *testing.T) {  })
	t.Run("A=2", func(t *testing.T) {  })
	t.Run("B=1", func(t *testing.T) {  })
	// <tear-down code>
}
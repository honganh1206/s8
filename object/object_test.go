package object

import "testing"

func TestStringHasKey(t *testing.T) {
	hello1 := &String{Value: "Hello World"}
	hello2 := &String{Value: "Hello World"}
	diff1 := &String{Value: "My name is John"}
	diff2 := &String{Value: "My name is John"}

	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("strings with same content but have different hash keys")
	}

	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("strings with same content but have different hash keys")
	}

	// Ensure proper hash key distribution
	if hello1.HashKey() == diff1.HashKey() {
		t.Errorf("strings with different content but have same hash keys")
	}
}

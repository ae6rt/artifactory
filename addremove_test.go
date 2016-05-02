package artifactory

import "testing"

func TestAdd(t *testing.T) {
	if !contains([]string{"a", "b"}, "a") {
		t.Fatal("Expecting to contain a")
	}
	if contains([]string{"a", "b"}, "c") {
		t.Fatal("Not expecting to contain c")
	}
}

func TestRemove(t *testing.T) {
	var arr []string

	arr = remove([]string{"a", "b", "a"}, "a")
	if len(arr) != 1 {
		t.Fatalf("Expecting 1 but got %d\n", len(arr))
	}
	if arr[0] != "b" {
		t.Fatalf("Expecting b but got %s\n", arr[0])
	}

	arr = remove([]string{"a", "b"}, "c")
	if len(arr) != 2 {
		t.Fatalf("Expecting 2 but got %d\n", len(arr))
	}
	if arr[0] != "a" {
		t.Fatalf("Expecting a but got %s\n", arr[0])
	}
	if arr[1] != "b" {
		t.Fatalf("Expecting b but got %s\n", arr[1])
	}

}

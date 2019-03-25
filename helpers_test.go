package main

import (
	"errors"
	"testing"
)

func TestValidateHash(t *testing.T) {
	var input = []struct {
		hash     string
		expected error
	}{
		{"8d9d7c59eec27a7aee55536582e45afb18f072c282edd22474a0db0676d74299", nil},
		{"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", errors.New("not 64 bytes long")},
		{"abc", errors.New("not 64 bytes long")},
		{"@#$%$#%^", errors.New("not hexadecimal")},
		{string([]byte{108, 147, 179, 39, 77, 210, 161, 13, 42, 14, 72, 143, 191, 27, 138, 193, 66, 117, 65, 143, 239, 52, 234, 255, 27, 102, 71, 163, 146, 162, 176, 179, 134, 103, 52, 181, 38, 102, 128, 37, 222, 209, 83, 61, 59, 217, 182, 183, 146, 212, 134, 109, 69, 208, 159, 129, 136, 134, 59, 229, 128, 169, 230, 204}), errors.New("not hexadecimal")},
	}

	for _, test := range input {
		err := validateHash(test.hash)

		if err == nil {
			if test.expected != nil {
				t.Errorf(`validateHash(%q) = %q, want: %q`, test.hash, err, test.expected)
			}
		} else if err.Error() != test.expected.Error() {
			t.Errorf(`validateHash(%q) = %q, want: %q`, test.hash, err, test.expected)
		}
	}
}

/*
func TestWriteJsonError(t *testing.T) {
	// TODO test

	w := httptest.NewRecorder()
	writeJsonError(w, "abc")

	response := w.Result()
	body, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		t.Error("http status code was not 200")
	}
	if response.Header.Get("Content-Type") != "application/json" {
		t.Error("Content-Type header was wrong")
	}
	fmt.Println(string(body))
}
*/

func TestReadAndValidateJsonMeetUp(t *testing.T) {
	// TODO test
}

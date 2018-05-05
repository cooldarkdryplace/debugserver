package debugserver

import "testing"

func TestRequestID(t *testing.T) {
	cases := []struct {
		name          string
		input         string
		expected      string
		expectedError error
	}{
		{
			name:          "valid case",
			input:         "/buckets/test/",
			expected:      "test",
			expectedError: nil,
		},
		{
			name:          "valid case",
			input:         "/buckets/_/",
			expected:      "_",
			expectedError: nil,
		},
		{
			name:          "valid case without trailing slash",
			input:         "/buckets/test",
			expected:      "test",
			expectedError: nil,
		},
		{
			name:          "valid case, with multiple levels",
			input:         "/buckets/test1/test2",
			expected:      "test1",
			expectedError: nil,
		},
		{
			name:          "invalid case with no ID and double slashes",
			input:         "/buckets//",
			expected:      "",
			expectedError: errNoID,
		},
		{
			name:          "invalid case with no ID",
			input:         "/buckets",
			expected:      "",
			expectedError: errNoID,
		},
		{
			name:          "invalid case with blank ID",
			input:         "/buckets/ /",
			expected:      "",
			expectedError: errNoID,
		},
	}

	for _, c := range cases {
		actual, err := requestID(c.input)
		if err != c.expectedError {
			t.Errorf("Got error: %v, expected: %v, id: %q, case: %q", err, c.expectedError, actual, c.name)
		}

		if actual != c.expected {
			t.Errorf("Got %q, expected: %q, case: %q", actual, c.expected, c.name)
		}
	}
}

package main

import "testing"

func TestCharPos(t *testing.T) {
	cases := []struct {
        str string
        idx int
        pos int
	}{
        { "", 0, 0 },
        { "", 1, 0 },
        { "", 23, 0 },
        { "a", 0, 0 },
        { "a", 1, 1 },
        { "a", 23, 1 },
        { "abcdef", 0, 0 },
        { "abcdef", 23, 6 },
        { "\tabc", 0, 0 },
        { "\tabc", 1, 4 },
        { "\tabc", 2, 5 },
        { "\tabc", 3, 6 },
        { "\tabc", 23, 7},
        { "\tabc\tde", 0, 0},
        { "\tabc\tde", 1, 4},
        { "\tabc\tde", 4, 7},
        { "\tabc\tde", 5, 11},
        { "\tabc\tde", 23, 13},
        { "\t\tde", 0, 0},
        { "\t\tde", 1, 4},
        { "\t\tde", 2, 8},
	}

	for _, c := range cases {
		expected := c.pos
        actual := charPos(c.str, c.idx, 4)

		if expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	}
}

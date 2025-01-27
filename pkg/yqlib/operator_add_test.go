package yqlib

import (
	"testing"
)

var addOperatorScenarios = []expressionScenario{
	{
		description: "Concatenate and assign arrays",
		document:    `{a: {val: thing, b: [cat,dog]}}`,
		expression:  ".a.b += [\"cow\"]",
		expected: []string{
			"D0, P[], (doc)::{a: {val: thing, b: [cat, dog, cow]}}\n",
		},
	},
	{
		description: "Concatenate arrays",
		document:    `{a: [1,2], b: [3,4]}`,
		expression:  `.a + .b`,
		expected: []string{
			"D0, P[a], (!!seq)::[1, 2, 3, 4]\n",
		},
	},
	{
		skipDoc:    true,
		expression: `[1] + ([2], [3])`,
		expected: []string{
			"D0, P[], (!!seq)::- 1\n- 2\n",
			"D0, P[], (!!seq)::- 1\n- 3\n",
		},
	},
	{
		description: "Concatenate null to array",
		document:    `{a: [1,2]}`,
		expression:  `.a + null`,
		expected: []string{
			"D0, P[a], (!!seq)::[1, 2]\n",
		},
	},
	{
		description: "Add new object to array",
		document:    `a: [{dog: woof}]`,
		expression:  `.a + {"cat": "meow"}`,
		expected: []string{
			"D0, P[a], (!!seq)::[{dog: woof}, {cat: meow}]\n",
		},
	},
	{
		description: "Add string to array",
		document:    `{a: [1,2]}`,
		expression:  `.a + "hello"`,
		expected: []string{
			"D0, P[a], (!!seq)::[1, 2, hello]\n",
		},
	},
	{
		description: "Update array (append)",
		document:    `{a: [1,2], b: [3,4]}`,
		expression:  `.a = .a + .b`,
		expected: []string{
			"D0, P[], (doc)::{a: [1, 2, 3, 4], b: [3, 4]}\n",
		},
	},
	{
		description: "String concatenation",
		document:    `{a: cat, b: meow}`,
		expression:  `.a = .a + .b`,
		expected: []string{
			"D0, P[], (doc)::{a: catmeow, b: meow}\n",
		},
	},
	{
		description: "Relative string concatenation",
		document:    `{a: cat, b: meow}`,
		expression:  `.a += .b`,
		expected: []string{
			"D0, P[], (doc)::{a: catmeow, b: meow}\n",
		},
	},
	{
		description:    "Number addition - float",
		subdescription: "If the lhs or rhs are floats then the expression will be calculated with floats.",
		document:       `{a: 3, b: 4.9}`,
		expression:     `.a = .a + .b`,
		expected: []string{
			"D0, P[], (doc)::{a: 7.9, b: 4.9}\n",
		},
	},
	{
		description:    "Number addition - int",
		subdescription: "If both the lhs and rhs are ints then the expression will be calculated with ints.",
		document:       `{a: 3, b: 4}`,
		expression:     `.a = .a + .b`,
		expected: []string{
			"D0, P[], (doc)::{a: 7, b: 4}\n",
		},
	},
	{
		description: "Increment number",
		document:    `{a: 3}`,
		expression:  `.a += 1`,
		expected: []string{
			"D0, P[], (doc)::{a: 4}\n",
		},
	},
}

func TestAddOperatorScenarios(t *testing.T) {
	for _, tt := range addOperatorScenarios {
		testScenario(t, &tt)
	}
	documentScenarios(t, "Add", addOperatorScenarios)
}

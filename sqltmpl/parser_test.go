package sqltmpl

import (
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

// QueryParsingTest represents a single t of prsr parsing. Given an [Input]
// query, if the actual result of parsing does not match the [Expected]
// string, the t fails
type QueryParsingTest struct {
	Name               string
	Input              string
	Expected           string
	ExpectedParameters int
}

// ParameterParsingTest pepresents a single t of parameter parsing.  Given
// the [prsr] and a set of [Parameters], if the actual parameter output from
// GetParsedParameters() matches the given [ExpectedParameters].  These tests
// specifically check type of output parameters, too.
type ParameterParsingTest struct {
	Name               string
	Query              string
	Parameters         []TestQueryParameter
	ExpectedParameters []interface{}
}

type TestQueryParameter struct {
	Name  string
	Value interface{}
}

func TestQueryParsing(t *testing.T) {
	var prsr Parser

	// Each of these represents a single t.
	QueryParsingTests := []QueryParsingTest{
		{
			Input:    "SELECT * FROM table WHERE col1 = 1",
			Expected: "SELECT * FROM table WHERE col1 = 1",
			Name:     "NoParameter",
		},
		{
			Input:              "SELECT * FROM table WHERE col1 = $name",
			Expected:           "SELECT * FROM table WHERE col1 = $1",
			ExpectedParameters: 1,
			Name:               "SingleParameter",
		},
		{
			Input:              "SELECT * FROM table WHERE col1 = $name::text",
			Expected:           "SELECT * FROM table WHERE col1 = $1::text",
			ExpectedParameters: 1,
			Name:               "With type assertion",
		},
		{
			Input:              "SELECT * FROM table WHERE col1 = $name AND col2 = $occupation",
			Expected:           "SELECT * FROM table WHERE col1 = $1 AND col2 = $2",
			ExpectedParameters: 2,
			Name:               "TwoParameters",
		},
		{
			Input:              "SELECT * FROM table WHERE col1 = $name AND col2 = $name",
			Expected:           "SELECT * FROM table WHERE col1 = $1 AND col2 = $2",
			ExpectedParameters: 2,
			Name:               "OneParameterMultipleTimes",
		},
		{
			Input:              "SELECT * FROM table WHERE col1 IN ($something, $else)",
			Expected:           "SELECT * FROM table WHERE col1 IN ($1, $2)",
			ExpectedParameters: 2,
			Name:               "ParametersInParenthesis",
		},
		{
			Input:    "SELECT * FROM table WHERE col1 = '$literal' AND col2 LIKE '$literal'",
			Expected: "SELECT * FROM table WHERE col1 = '$literal' AND col2 LIKE '$literal'",
			Name:     "ParametersInQuotes",
		},
		{
			Input:              "SELECT * FROM table WHERE col1 = '$literal' AND col2 = $literal AND col3 LIKE '$literal'",
			Expected:           "SELECT * FROM table WHERE col1 = '$literal' AND col2 = $1 AND col3 LIKE '$literal'",
			ExpectedParameters: 1,
			Name:               "ParametersInQuotes2",
		},
		{
			Input:              "SELECT * FROM table WHERE col1 = $foo AND col2 IN (SELECT id FROM tabl2 WHERE col10 = $bar)",
			Expected:           "SELECT * FROM table WHERE col1 = $1 AND col2 IN (SELECT id FROM tabl2 WHERE col10 = $2)",
			ExpectedParameters: 2,
			Name:               "ParametersInSubclause",
		},
		{
			Input:              "SELECT * FROM table WHERE col1 = $1234567890 AND col2 = $0987654321",
			Expected:           "SELECT * FROM table WHERE col1 = $1 AND col2 = $2",
			ExpectedParameters: 2,
			Name:               "NumericParameters",
		},
		{
			Input:              "SELECT * FROM table WHERE col1 = $ABCDEFGHIJKLMNOPQRSTUVWXYZ",
			Expected:           "SELECT * FROM table WHERE col1 = $1",
			ExpectedParameters: 1,
			Name:               "CapsParameters",
		},
		{
			Input:              "SELECT * FROM table WHERE col1 = $abc123ABC098",
			Expected:           "SELECT * FROM table WHERE col1 = $1",
			ExpectedParameters: 1,
			Name:               "AltcapsParameters",
		},
		{
			Input:              "SELECT * FROM table WHERE col1 LIKE %$t%",
			Expected:           "SELECT * FROM table WHERE col1 LIKE %$1%",
			ExpectedParameters: 1,
			Name:               "pattern matching",
		},
		{
			Input:              "ST_GeomFromText('POINT(' || $long $lat || ',4326)'",
			Expected:           "ST_GeomFromText('POINT(' || $1 $2 || ',4326)'",
			ExpectedParameters: 2,
			Name:               "inside quotes",
		},
	}

	for _, parsingTest := range QueryParsingTests {
		prsr = NewParser(parsingTest.Input)
		require.Equal(t, parsingTest.Expected, prsr.GetParsedQuery(), parsingTest.Name)
		require.Equal(t, parsingTest.ExpectedParameters, len(prsr.GetParsedParameters()), parsingTest.Name)
	}
}

/*
	Tests to ensure that setting parameter values turns out correct when using GetParsedParameters().
	These tests ensure correct positioning and type.
*/
func TestParameterReplacement(t *testing.T) {
	var prsr Parser
	var parameterMap map[string]interface{}

	// note that if you're adding or editing these tests,
	// you'll also want to edit the associated struct for this t below,
	// in the next t func.
	QueryVariableTests := []ParameterParsingTest{
		{
			Name:  "SingleStringParameter",
			Query: "SELECT * FROM table WHERE col1 = $foo",
			Parameters: []TestQueryParameter{
				{
					Name:  "foo",
					Value: "bar",
				},
			},
			ExpectedParameters: []interface{}{
				"bar",
			},
		},
		{
			Name:  "TwoStringParameter",
			Query: "SELECT * FROM table WHERE col1 = $foo AND col2 = $foo2",
			Parameters: []TestQueryParameter{
				{
					Name:  "foo",
					Value: "bar",
				},
				{
					Name:  "foo2",
					Value: "bart",
				},
			},
			ExpectedParameters: []interface{}{
				"bar", "bart",
			},
		},
		{
			Name:  "TwiceOccurringParameter",
			Query: "SELECT * FROM table WHERE col1 = $foo AND col2 = $foo",
			Parameters: []TestQueryParameter{
				{
					Name:  "foo",
					Value: "bar",
				},
			},
			ExpectedParameters: []interface{}{
				"bar", "bar",
			},
		},
		{
			Name:  "ParameterTyping",
			Query: "SELECT * FROM table WHERE col1 = $str AND col2 = $int AND col3 = $pi",
			Parameters: []TestQueryParameter{
				{
					Name:  "str",
					Value: "foo",
				},
				{
					Name:  "int",
					Value: 1,
				},
				{
					Name:  "pi",
					Value: 3.14,
				},
			},
			ExpectedParameters: []interface{}{
				"foo", 1, 3.14,
			},
		},
		{
			Name:  "ParameterOrdering",
			Query: "SELECT * FROM table WHERE col1 = $foo AND col2 = $bar AND col3 = $foo AND col4 = $foo AND col5 = $bar",
			Parameters: []TestQueryParameter{
				{
					Name:  "foo",
					Value: "something",
				},
				{
					Name:  "bar",
					Value: "else",
				},
			},
			ExpectedParameters: []interface{}{
				"something", "else", "something", "something", "else",
			},
		},
		{
			Name:  "ParameterCaseSensitivity",
			Query: "SELECT * FROM table WHERE col1 = $foo AND col2 = $FOO",
			Parameters: []TestQueryParameter{
				{
					Name:  "foo",
					Value: "baz",
				},
				{
					Name:  "FOO",
					Value: "quux",
				},
			},
			ExpectedParameters: []interface{}{
				"baz", "quux",
			},
		},
		{
			Name:  "ParameterNil",
			Query: "SELECT * FROM table WHERE col1 = $foo",
			Parameters: []TestQueryParameter{
				{
					Name:  "foo",
					Value: pq.Array([]string{}),
				},
			},
			ExpectedParameters: []interface{}{
				pq.Array([]string{}),
			},
		},
	}

	// run variable tests.
	for _, variableTest := range QueryVariableTests {
		// parse prsr and set values.
		parameterMap = make(map[string]interface{}, 8)
		prsr = NewParser(variableTest.Query)

		for _, queryVariable := range variableTest.Parameters {
			prsr.SetValue(queryVariable.Name, queryVariable.Value)
			parameterMap[queryVariable.Name] = queryVariable.Value
		}

		// t outputs
		for index, queryVariable := range prsr.GetParsedParameters() {
			require.Equal(t, queryVariable, variableTest.ExpectedParameters[index])
		}

		prsr = NewParser(variableTest.Query)
		prsr.SetValuesFromMap(parameterMap)

		// t map parameter outputs.
		for index, queryVariable := range prsr.GetParsedParameters() {
			require.Equal(t, queryVariable, variableTest.ExpectedParameters[index])
		}
	}
}

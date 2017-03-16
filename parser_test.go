package pc

import (
	"encoding/json"
	"fmt"
	"regexp"
)

func Example_tutrial_001_parse_a() {
	a := Rune('a')
	result, _ := a.Parse(NewInMemoryReader([]byte("abc")))
	fmt.Printf("%+v", result)
	// Output: {Value:a TextRange:[1:1,1:2]}
}

func Example_tutrial_002_parse_ab() {
	ab := And(Rune('a'), Rune('b'))
	result, _ := ab.Parse(NewInMemoryReader([]byte("abc")))
	fmt.Printf("%+v", result)
	// Output: {Value:[{Value:a TextRange:[1:1,1:2]} {Value:b TextRange:[1:2,1:3]}] TextRange:[1:1,1:3]}
}

func Example_tutrial_003_parse_ab_and_transform_it() {
	ab := TransformAsResults(
		And(Rune('a'), Rune('b')),
		func(results []ParseResult) (interface{}, error) {
			return results[0].Value.(string) + results[1].Value.(string), nil
		},
	)
	result, _ := ab.Parse(NewInMemoryReader([]byte("abc")))
	fmt.Printf("%+v", result)
	// Output: {Value:ab TextRange:[1:1,1:3]}
}

func Example_tutrial_003_2_parse_string_ab() {
	ab := String("ab")
	result, _ := ab.Parse(NewInMemoryReader([]byte("abc")))
	fmt.Printf("%+v", result)
	// Output: {Value:ab TextRange:[1:1,1:3]}
}

func ExampleCompose() {
	result, _ := Compose(
		func(parser Parser) Parser {
			return Index(parser, 1)
		},
		FilterNil,
		Flatten,
	)(
		And(
			Const(Rune('a'), nil),
			And(
				Rune('b'),
				Const(Rune('c'), nil),
			),
			Rune('d'),
		),
	).Parse(NewInMemoryReader([]byte("abcde")))
	fmt.Println(string(result.Value.(string)))
	// Output: d
}

func ExampleRune() {
	result, _ := Rune('a').Parse(NewInMemoryReader([]byte("abc")))
	fmt.Println(string(result.Value.(string)))
	// Output: a
}

func ExampleRune_not_matched() {
	_, err := Rune('b').Parse(NewInMemoryReader([]byte("abc")))
	fmt.Println(err)
	// Output: not matched
}

func ExampleRuneIn() {
	result, _ := RuneIn("abc").Parse(NewInMemoryReader([]byte("abc")))
	fmt.Println(string(result.Value.(string)))
	// Output: a
}

func ExampleRuneNotIn() {
	result, _ := RuneNotIn("def").Parse(NewInMemoryReader([]byte("abc")))
	fmt.Println(string(result.Value.(string)))
	// Output: a
}

func ExampleRuneInRange() {
	result, _ := RuneInRange('a', 'b').Parse(NewInMemoryReader([]byte("abc")))
	fmt.Println(string(result.Value.(string)))
	// Output: a
}

func ExampleRuneNotInRange() {
	result, _ := RuneNotInRange('b', 'c').Parse(NewInMemoryReader([]byte("abc")))
	fmt.Println(string(result.Value.(string)))
	// Output: a
}

func ExampleAnyRune() {
	result, _ := AnyRune().Parse(NewInMemoryReader([]byte("abc")))
	fmt.Println(string(result.Value.(string)))
	// Output: a
}

func ExampleRegexp() {
	result, _ := Regexp(regexp.MustCompile(`^ab?c[1-3]+d[ef]+`)).Parse(NewInMemoryReader([]byte("abc123def456")))
	fmt.Println(result.Value.(string))
	// Output: abc123def
}

func ExampleRegexp_bad_case() {
	_, err := Regexp(regexp.MustCompile(`bc`)).Parse(NewInMemoryReader([]byte("abc")))
	fmt.Println(err)
	// Output: not matched
}

func ExampleString() {
	result, _ := String("ab").Parse(NewInMemoryReader([]byte("abc")))
	fmt.Println(result.Value.(string))
	// Output: ab
}

func ExampleStringByAnd() {
	result, _ := StringByAnd(
		Rune('a'),
		AnyRune(),
		Rune('c'),
	).Parse(NewInMemoryReader([]byte("abc")))
	fmt.Println(result.Value.(string))
	// Output: abc
}

func ExampleConst() {
	result, _ := Const(String("abc"), "foobar").Parse(NewInMemoryReader([]byte("abc")))
	fmt.Println(result.Value.(string))
	// Output: foobar
}

func ExampleFlatten() {
	result, _ := ConcatString(
		Flatten(
			And(
				Rune('a'),
				Many(Rune('b')),
			),
		),
	).Parse(NewInMemoryReader([]byte("abbc")))
	fmt.Println(result.Value.(string))
	// Output: abb
}

func ExampleManyMinMaxTerminate() {
	result, _ := ConcatString(
		ManyMinMaxTerminate(
			RuneIn("abc"),
			0,
			0,
			Rune('b'),
		),
	).Parse(NewInMemoryReader([]byte("abc")))
	fmt.Println(result.Value.(string))
	// Output: ab
}

func ExampleMaybe() {
	result, _ := StringByAnd(
		Rune('a'),
		Maybe(String("bd")),
		Rune('b'),
		Maybe(Rune('c')),
	).Parse(NewInMemoryReader([]byte("abc")))
	fmt.Println(result.Value.(string))
	// Output: abc
}

func ExampleOr() {
	result, _ := Or(
		String("foo"),
		String("abc"),
	).Parse(NewInMemoryReader([]byte("abc")))
	fmt.Println(result.Value.(string))
	// Output: abc
}

func ExampleSurround() {
	result, _ := Surround("\"", AnyRune(), "\"", "\\").Parse(NewInMemoryReader([]byte("\"foo\\\\bar\\\"fizz\\buzz\"")))
	fmt.Println(result.Value.(string))
	// Output: foo\bar"fizz\buzz
}

func ExampleSurround_no_escape() {
	result, _ := Surround("foo", AnyRune(), "foo", "\\").Parse(NewInMemoryReader([]byte("foo123barfoo456")))
	fmt.Println(result.Value.(string))
	// Output: 123bar
}

func Example_json() {
	ws := RuneIn(" \t\n")
	ws0 := Many(ws)

	stringParser := Surround("\"", AnyRune(), "\"", "\\")

	digit19 := RuneInRange('1', '9')
	digit09 := RuneInRange('0', '9')
	numberParser := TransformToNumber(
		StringByAnd(
			Maybe(Rune('-')),
			Or(
				Rune('0'),
				StringByAnd(digit19, StringByMany(digit09)),
			),
			Maybe(StringByAnd(Rune('.'), StringByMany1(digit09))),
			Maybe(StringByAnd(RuneIn("eE"), Maybe(RuneIn("+-")), StringByMany(digit09))),
		),
	)

	booleanParser := Or(
		Const(String("true"), true),
		Const(String("false"), false),
	)

	nullParser := Const(String("null"), nil)

	var objectParser, arrayParser Parser

	originalValueParser := Lazy(func(me Parser) Parser {
		return Or(
			stringParser,
			numberParser,
			objectParser,
			arrayParser,
			booleanParser,
			nullParser,
		)
	})

	valueParser := Annotate(originalValueParser, "JsonValue", "...JsonValue...")

	type KeyValue struct {
		Key   string
		Value interface{}
	}

	keyValueParser := TransformAsResults(
		And(stringParser, ws0, Rune(':'), ws0, valueParser),
		func(results []ParseResult) (interface{}, error) {
			return KeyValue{results[0].Value.(string), results[4].Value}, nil
		},
	)

	objectParser = TransformAsResults(
		And(
			Rune('{'),
			ws0,
			Maybe(Separated(keyValueParser, And(ws0, Rune(','), ws0))),
			ws0,
			Rune('}'),
		),
		func(results []ParseResult) (interface{}, error) {
			if results[2].Value == nil {
				return nil, nil
			}
			kvMap := make(map[string]interface{})
			kvResults := results[2].Value.([]ParseResult)
			for _, kvResult := range kvResults {
				kv := kvResult.Value.(KeyValue)
				kvMap[kv.Key] = kv.Value
			}
			return kvMap, nil
		},
	)

	arrayParser = TransformAsResults(
		And(
			Rune('['),
			ws0,
			Maybe(Separated(valueParser, And(ws0, Rune(','), ws0))),
			ws0,
			Rune(']'),
		),
		func(results []ParseResult) (interface{}, error) {
			valResults := results[2].Value.([]ParseResult)
			vals := make([]interface{}, 0, len(valResults))
			for _, valResult := range valResults {
				vals = append(vals, valResult.Value)
			}
			return vals, nil
		},
	)

	jsonParser := originalValueParser

	jsonSample := `{"foo": 100, "bar": [1.23, {"fizz": -1.2E-1}, "buzz"]}`
	parseResult, _ := jsonParser.Parse(NewInMemoryReader([]byte(jsonSample)))
	jsonFromResult, _ := json.Marshal(parseResult.Value)
	fmt.Println(string(jsonFromResult))
	// Output: {"bar":[1.23,{"fizz":-0.12},"buzz"],"foo":100}
}

func ExampleStringifyParser_via_fmt_stringer() {
	parser := And(Rune('a'), String("bc"), Maybe(RuneIn("def")))
	fmt.Println(parser)
	// Output: abc([def])?
}

func ExampleStringifyParser_diagnostics() {
	parser := And(Rune('a'), String("bc"), Maybe(RuneIn("def")))
	fmt.Println(StringifyParser(parser, PatternStyleDiagnostics))
	// Output: <And:<Rune:a><String:bc><Transform:<Many:(<RuneIn:[def]>)?>>>
}

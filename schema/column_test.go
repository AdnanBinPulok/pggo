package schema

import "testing"

func TestDefaultValue_KnownExpression_NotQuoted(t *testing.T) {
	def := DataTypeFactory{}.TimestampTz().NotNull().DefaultValue("CURRENT_TIMESTAMP")

	got := def.String()
	want := "TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP"
	if got != want {
		t.Fatalf("unexpected column def string\nwant: %s\ngot:  %s", want, got)
	}
}

func TestDefaultValue_StringLiteral_QuotedAndEscaped(t *testing.T) {
	def := DataTypeFactory{}.Text().DefaultValue("O'Hare")

	got := def.String()
	want := "TEXT DEFAULT 'O''Hare'"
	if got != want {
		t.Fatalf("unexpected column def string\nwant: %s\ngot:  %s", want, got)
	}
}

func TestDefaultValue_BoolAndNumber_NotQuoted(t *testing.T) {
	boolDef := DataTypeFactory{}.Bool().DefaultValue(true)
	if got, want := boolDef.String(), "BOOLEAN DEFAULT TRUE"; got != want {
		t.Fatalf("unexpected bool default\nwant: %s\ngot:  %s", want, got)
	}

	intDef := DataTypeFactory{}.Integer().DefaultValue(42)
	if got, want := intDef.String(), "INTEGER DEFAULT 42"; got != want {
		t.Fatalf("unexpected int default\nwant: %s\ngot:  %s", want, got)
	}
}

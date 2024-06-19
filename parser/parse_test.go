package parser

import (
	"testing"
)

func TestMinimalRule(t *testing.T) {
	cmd := "Start[iri].Eval"
	chain, err := ParseCommand(cmd)
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(chain) != 2 {
		t.Errorf("Expected 2 steps, got %d", len(chain))
	}

	expected := []Step{
		{token: Start, arg: "iri"},
		{token: Eval},
	}

	for i, e := range expected {
		if chain[i].token != e.token || chain[i].arg != e.arg {
			t.Errorf("expected %+v got %+v in step %d", e, chain[i], i)
		}
	}
}

func TestFindNext(t *testing.T) {
	cmd := `field1, value1, "3.14", "value with \" in it"].Eval`
	i, ok := findNext(cmd, ']')
	if !ok {
		t.Errorf("expected to find closing bracket")
	}
	expected := len(cmd) - len(".Eval") - 1
	if i != expected {
		t.Errorf("expected to find ] in position %d got %d", expected, i)
	}

	if cmd[i] != ']' {
		t.Errorf("expected to find ] got %c", cmd[i])
	}
}

func TestRuleWithHasValue(t *testing.T) {
	cmd := `Start[iri].HasValue[field1, value1, "3.14", "value with \" in it"].Eval`
	chain, err := ParseCommand(cmd)
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(chain) != 3 {
		t.Errorf("Expected 3 steps, got %d", len(chain))
	}

	expected := []Step{
		{token: Start, arg: "iri"},
		{token: HasValue, arg: "field1", vals: []string{"value1", "\"3.14\"", "value with \" in it"}},
		{token: Eval},
	}

	for i, e := range expected {
		if chain[i].token != e.token || chain[i].arg != e.arg {
			t.Errorf("expected %+v got %+v in step %d", e, chain[i], i)
		}
	}
}

func TestRuleWithAnd(t *testing.T) {
	cmd := `Start[iri].And(HasType[TypeAnd1].HasType[TypeAnd2]).Eval`
	chain, err := ParseCommand(cmd)
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(chain) != 4 {
		t.Errorf("Expected 4 steps, got %d", len(chain))
	}

	expected := []Step{
		{token: Start, arg: "iri"},
		{token: HasType, arg: "TypeAnd1"},
		{token: HasType, arg: "TypeAnd2"},
		{token: Eval},
	}

	for i, e := range expected {
		if chain[i].token != e.token || chain[i].arg != e.arg {
			t.Errorf("expected %+v got %+v in step %d", e, chain[i], i)
		}
	}
}

func TestRuleWithOr(t *testing.T) {
	cmd := `Start[iri].Or(HasType[TypeOr1].HasType[TypeOr2]).Eval`
	chain, err := ParseCommand(cmd)
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(chain) != 3 {
		t.Errorf("Expected 3 steps, got %d", len(chain))
	}

	expected := []Step{
		{token: Start, arg: "iri"},
		{token: Or, arg: ""},
		{token: Eval},
	}

	for i, e := range expected {
		if chain[i].token != e.token || chain[i].arg != e.arg {
			t.Errorf("expected %+v got %+v in step %d", e, chain[i], i)
		}
	}

	expected = []Step{
		{token: HasType, arg: "TypeOr1"},
		{token: HasType, arg: "TypeOr2"},
	}

	for i, e := range expected {
		if chain[1].subcmd[i].token != e.token || chain[1].subcmd[i].arg != e.arg {
			t.Errorf("expected %+v got %+v in step %d", e, chain[1].subcmd[i], i)
		}
	}
}

func TestRuleWithAndInsideOr(t *testing.T) {
	cmd := `Start[iri].Or(HasType[TypeOr1].And(HasType[TypeAnd1].HasType[TypeAnd2])).Eval`
	chain, err := ParseCommand(cmd)
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(chain) != 3 {
		t.Errorf("Expected 3 steps, got %d", len(chain))
	}

	expected := []Step{
		{token: Start, arg: "iri"},
		{token: Or, arg: ""},
		{token: Eval},
	}

	for i, e := range expected {
		if chain[i].token != e.token || chain[i].arg != e.arg {
			t.Errorf("expected %+v got %+v in step %d", e, chain[i], i)
		}
	}

	expected = []Step{
		{token: HasType, arg: "TypeOr1"},
		{token: HasType, arg: "TypeAnd1"},
		{token: HasType, arg: "TypeAnd2"},
	}

	for i, e := range expected {
		if chain[1].subcmd[i].token != e.token || chain[1].subcmd[i].arg != e.arg {
			t.Errorf("expected %+v got %+v in step %d", e, chain[1].subcmd[i], i)
		}
	}
}

func TestTestInvalidRule1(t *testing.T) {
	cmd := `Start[iri].Or(HasType[TypeOr1].HasType[TypeOr2].Eval`
	_, err := ParseCommand(cmd)
	if err == nil {
		t.Errorf("Expected error when parsing %s", cmd)
	}
}

func TestCmdCleaning(t *testing.T) {
	cmd := `
		Start[iri]
			.HasType[<https://bsm.bloomberg.com/ontology/Company>]
			.HasValue[field1, "value1", "3.14"]
			.Or(
				HasType[TypeOr1]. HasType[TypeOr2]
			)
			.And(HasType[TypeAnd1].HasType[TypeAnd2])
			.Eval`
	clean, err := cleanCmd(cmd)
	if err != nil {
		t.Errorf(err.Error())
	}

	expected := "HasType[bsm:Company].HasValue[field1,\"value1\",\"3.14\"].Or(HasType[TypeOr1].HasType[TypeOr2]).HasType[TypeAnd1].HasType[TypeAnd2]"
	if clean != expected {
		t.Errorf("Expected %s got %s", expected, clean)
	}
}

func TestHasBroader(t *testing.T) {
	cmd := `Start[iri].HasBroader[tax,target].Eval`
	chain, err := ParseCommand(cmd)
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(chain) != 3 {
		t.Errorf("Expected 3 steps, got %d", len(chain))
	}

	expected := []Step{
		{token: Start, arg: "iri"},
		{token: HasBroader, arg: "tax", vals: []string{"target"}},
		{token: Eval},
	}

	for i, e := range expected {
		if chain[i].token != e.token || chain[i].arg != e.arg {
			t.Errorf("expected %+v got %+v in step %d", e, chain[i], i)
		}
	}
}

func TestInScheme(t *testing.T) {
	cmd := `Start[iri].InScheme[tax].Eval`
	chain, err := ParseCommand(cmd)
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(chain) != 3 {
		t.Errorf("Expected 3 steps, got %d", len(chain))
	}

	expected := []Step{
		{token: Start, arg: "iri"},
		{token: InScheme, arg: "tax"},
		{token: Eval},
	}

	for i, e := range expected {
		if chain[i].token != e.token || chain[i].arg != e.arg {
			t.Errorf("expected %+v got %+v in step %d", e, chain[i], i)
		}
	}
}

func TestIsInstance(t *testing.T) {
	cmd := `Start[iri].IsInstance[inst].Eval`
	chain, err := ParseCommand(cmd)
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(chain) != 3 {
		t.Errorf("Expected 3 steps, got %d", len(chain))
	}

	expected := []Step{
		{token: Start, arg: "iri"},
		{token: IsInstance, arg: "inst"},
		{token: Eval},
	}

	for i, e := range expected {
		if chain[i].token != e.token || chain[i].arg != e.arg {
			t.Errorf("expected %+v got %+v in step %d", e, chain[i], i)
		}
	}
}

func TestFollow(t *testing.T) {
	cmd := `Start[iri].Follow[rel].Eval`
	chain, err := ParseCommand(cmd)
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(chain) != 3 {
		t.Errorf("Expected 3 steps, got %d", len(chain))
	}

	expected := []Step{
		{token: Start, arg: "iri"},
		{token: Follow, arg: "rel"},
		{token: Eval},
	}

	for i, e := range expected {
		if chain[i].token != e.token || chain[i].arg != e.arg {
			t.Errorf("expected %+v got %+v in step %d", e, chain[i], i)
		}
	}
}

func TestFollowInverse(t *testing.T) {
	cmd := `Start[iri].FollowInverse[rel].Eval`
	chain, err := ParseCommand(cmd)
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(chain) != 3 {
		t.Errorf("Expected 3 steps, got %d", len(chain))
	}

	expected := []Step{
		{token: Start, arg: "iri"},
		{token: FollowInverse, arg: "rel"},
		{token: Eval},
	}

	for i, e := range expected {
		if chain[i].token != e.token || chain[i].arg != e.arg {
			t.Errorf("expected %+v got %+v in step %d", e, chain[i], i)
		}
	}
}
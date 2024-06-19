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

func TestHasValue(t *testing.T) {
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

func TestAnd(t *testing.T) {
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

func TestOr(t *testing.T) {
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

func TestOrWithModeCmd(t *testing.T) {
	cmd := `Start[iri].Or(HasType[TypeOr1].HasType[TypeOr2]).HasType[Type3].Eval`
	chain, err := ParseCommand(cmd)
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(chain) != 4 {
		t.Errorf("Expected 4 steps, got %d", len(chain))
	}

	expected := []Step{
		{token: Start, arg: "iri"},
		{token: Or, arg: ""},
		{token: HasType, arg: "Type3"},
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

func TestAndInsideOr(t *testing.T) {
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

func TestHasCategory(t *testing.T) {
	cmd := `Start[iri].HasCategory[cat].Eval`
	chain, err := ParseCommand(cmd)
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(chain) != 3 {
		t.Errorf("Expected 3 steps, got %d", len(chain))
	}

	expected := []Step{
		{token: Start, arg: "iri"},
		{token: HasCategory, arg: "cat"},
		{token: Eval},
	}

	for i, e := range expected {
		if chain[i].token != e.token || chain[i].arg != e.arg {
			t.Errorf("expected %+v got %+v in step %d", e, chain[i], i)
		}
	}
}

func TestIsActive(t *testing.T) {
	cmd := `Start[iri].IsActive[].Eval`
	chain, err := ParseCommand(cmd)
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(chain) != 3 {
		t.Errorf("Expected 3 steps, got %d", len(chain))
	}

	expected := []Step{
		{token: Start, arg: "iri"},
		{token: IsActive},
		{token: Eval},
	}

	for i, e := range expected {
		if chain[i].token != e.token || chain[i].arg != e.arg {
			t.Errorf("expected %+v got %+v in step %d", e, chain[i], i)
		}
	}
}

func TestIsInactive(t *testing.T) {
	cmd := `Start[iri].IsInactive[].Eval`
	chain, err := ParseCommand(cmd)
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(chain) != 3 {
		t.Errorf("Expected 3 steps, got %d", len(chain))
	}

	expected := []Step{
		{token: Start, arg: "iri"},
		{token: IsInactive},
		{token: Eval},
	}

	for i, e := range expected {
		if chain[i].token != e.token {
			t.Errorf("expected %+v got %+v in step %d", e, chain[i], i)
		}
	}
}

func TestIriToQnameBsm(t *testing.T) {
	cmd := `Start[iri].HasType[<https://bsm.bloomberg.com/ontology/Company>].Eval`
	cmd, err := cleanCmd(cmd)
	if err != nil {
		t.Errorf(err.Error())
	}

	expected := "HasType[bsm:Company]"
	if cmd != expected {
		t.Errorf("Expected %s got %s", expected, cmd)
	}
}

func TestIriToQnameInstance(t *testing.T) {
	cmd := `Start[iri].IsInstance[<https://bsm.bloomberg.com/instance/0xdecafbad>].Eval`
	cmd, err := cleanCmd(cmd)
	if err != nil {
		t.Errorf(err.Error())
	}

	expected := "IsInstance[bsi:0xdecafbad]"
	if cmd != expected {
		t.Errorf("Expected %s got %s", expected, cmd)
	}
}

func TestIriToQnameOwl(t *testing.T) {
	cmd := `Start[iri].HasType[<http://www.w3.org/2002/07/owl#Thing>].Eval`
	cmd, err := cleanCmd(cmd)
	if err != nil {
		t.Errorf(err.Error())
	}

	expected := "HasType[owl:Thing]"
	if cmd != expected {
		t.Errorf("Expected %s got %s", expected, cmd)
	}
}

func TestIriToQnameRdfs(t *testing.T) {
	cmd := `Start[iri].HasType[<http://www.w3.org/2000/01/rdf-schema#Class>].Eval`
	cmd, err := cleanCmd(cmd)
	if err != nil {
		t.Errorf(err.Error())
	}

	expected := "HasType[rdfs:Class]"
	if cmd != expected {
		t.Errorf("Expected %s got %s", expected, cmd)
	}
}

func TestIriToQnameRdf(t *testing.T) {
	cmd := `Start[iri].HasType[<http://www.w3.org/1999/02/22-rdf-syntax-ns#Property>].Eval`
	cmd, err := cleanCmd(cmd)
	if err != nil {
		t.Errorf(err.Error())
	}

	expected := "HasType[rdf:Property]"
	if cmd != expected {
		t.Errorf("Expected %s got %s", expected, cmd)
	}
}

func TestCommaInsideString(t *testing.T) {
	cmd := `Start[iri].HasValue[field1, "value1, with comma", "3.14"].Eval`
	chain, err := ParseCommand(cmd)
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(chain) != 3 {
		t.Errorf("Expected 3 steps, got %d", len(chain))
	}

	expected := []Step{
		{token: Start, arg: "iri"},
		{token: HasValue, arg: "field1", vals: []string{"value1, with comma", "\"3.14\""}},
		{token: Eval},
	}

	for i, e := range expected {
		if chain[i].token != e.token || chain[i].arg != e.arg {
			t.Errorf("expected %+v got %+v in step %d", e, chain[i], i)
		}
	}
}

func TestSpaceInsideString(t *testing.T) {
	cmd := `Start[iri].HasValue[field1, "value1 with space", "3.14"].Eval`
	chain, err := ParseCommand(cmd)
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(chain) != 3 {
		t.Errorf("Expected 3 steps, got %d", len(chain))
	}

	expected := []Step{
		{token: Start, arg: "iri"},
		{token: HasValue, arg: "field1", vals: []string{"value1 with space", "\"3.14\""}},
		{token: Eval},
	}

	for i, e := range expected {
		if chain[i].token != e.token || chain[i].arg != e.arg {
			t.Errorf("expected %+v got %+v in step %d", e, chain[i], i)
		}
	}
}

func TestParenInTaxonomyName(t *testing.T) {
	cmd := `Start[iri].HasBroader["tax (with paren)", "instance"].Eval`
	chain, err := ParseCommand(cmd)
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(chain) != 3 {
		t.Errorf("Expected 3 steps, got %d", len(chain))
	}

	expected := []Step{
		{token: Start, arg: "iri"},
		{token: HasBroader, arg: "tax (with paren)", vals: []string{"instance"}},
		{token: Eval},
	}

	for i, e := range expected {
		if chain[i].token != e.token || chain[i].arg != e.arg {
			t.Errorf("expected %+v got %+v in step %d", e, chain[i], i)
		}
	}
}

func TestStripQuotesFromValue(t *testing.T) {
	cmd := `Start[iri].HasValue[field1, "3.14"].Eval`
	chain, err := ParseCommand(cmd)
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(chain) != 3 {
		t.Errorf("Expected 3 steps, got %d", len(chain))
	}

	expected := []Step{
		{token: Start, arg: "iri"},
		{token: HasValue, arg: "field1", vals: []string{"3.14"}},
		{token: Eval},
	}

	for i, e := range expected {
		if chain[i].token != e.token || chain[i].arg != e.arg {
			t.Errorf("expected %+v got %+v in step %d", e, chain[i], i)
		}

		if len(chain[i].vals) != len(e.vals) {
			t.Errorf("expected %d values got %d", len(e.vals), len(chain[i].vals))
		}

		for j, v := range e.vals {
			if chain[i].vals[j] != v {
				t.Errorf("expected %s got %s", v, chain[i].vals[j])
			}
		}
	}
}

func TestInvalidRule1(t *testing.T) {
	cmd := `Start[iri].Or(HasType[TypeOr1].HasType[TypeOr2].Eval`
	_, err := ParseCommand(cmd)
	if err == nil {
		t.Errorf("Expected error when parsing %s", cmd)
	}
}

func TestInvalidRuleUnquotedPeriod(t *testing.T) {
	cmd := `Start[iri].HasValue[field, 3.14].Eval`
	_, err := ParseCommand(cmd)
	if err == nil {
		t.Errorf("Expected error when parsing %s", cmd)
	}
}

func TestInvalidRuleMissingValue(t *testing.T) {
	cmd := `Start[iri].HasValue[field].Eval`
	_, err := ParseCommand(cmd)
	if err == nil {
		t.Errorf("Expected error when parsing %s", cmd)
	}
}

func TestInvalidRuleTooManyArgs(t *testing.T) {
	cmd := `Start[iri].InScheme[v1,v2,v3].Eval`
	_, err := ParseCommand(cmd)
	if err == nil {
		t.Errorf("Expected error when parsing %s", cmd)
	}
}

func TestMismatchParen1(t *testing.T) {
	cmd := `Start[iri].Or(HasType[TypeOr1].HasType[TypeOr2].Eval`
	_, err := ParseCommand(cmd)
	if err == nil {
		t.Errorf("Expected error when parsing %s", cmd)
	}
}

func TestMismatchParen2(t *testing.T) {
	cmd := `Start[iri].Or(HasType[TypeOr1].HasType[TypeOr2]).Eval)`
	_, err := ParseCommand(cmd)
	if err == nil {
		t.Errorf("Expected error when parsing %s", cmd)
	}
}

func TestMismatchParent3(t *testing.T) {
	cmd := `Start[iri].Or((HasType[TypeOr1].HasType[TypeOr2].Eval`
	_, err := ParseCommand(cmd)
	if err == nil {
		t.Errorf("Expected error when parsing %s", cmd)
	}
}

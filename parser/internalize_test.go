package parser

import (
	"reflect"
	"testing"
)

type IidStore struct {
	toIid   map[string]Iid
	fromIid map[Iid]string
	Next    Iid
}

func (s *IidStore) GetIid(iri string) (Iid, bool) {
	i, ok := s.toIid[iri]
	return i, ok
}

func (s *IidStore) GetString(i Iid) (string, bool) {
	iri, ok := s.fromIid[i]
	return iri, ok
}

func (s *IidStore) Put(iri string) Iid {
	i, ok := s.toIid[iri]
	if ok {
		return i
	}
	i = s.Next
	s.toIid[iri] = i
	s.fromIid[i] = iri
	s.Next++
	return i
}

func NewIidStore() *IidStore {
	is := IidStore{
		toIid:   make(map[string]Iid),
		fromIid: make(map[Iid]string),
	}

	is.Put("Error")
	is.Put("NoMatch")
	return &is
}

func TestInternalizeIdentityCmd(t *testing.T) {
	is := NewIidStore()
	cmd := `Start[iri].Eval`

	steps, err := ParseCommand(cmd)
	if err != nil {
		t.Error(err)
	}

	isteps, err := InternalizeSteps(steps, is)
	if err != nil {
		t.Error(err)
	}

	if len(isteps) != 2 {
		t.Errorf("Expected 2 steps, got %d", len(isteps))
	}
}

func TestInternalizeIsInstanceCmd(t *testing.T) {
	is := NewIidStore()
	red := is.Put("red")
	cmd := `Start[iri].IsInstance[red].Eval`

	steps, err := ParseCommand(cmd)
	if err != nil {
		t.Error(err)
	}

	isteps, err := InternalizeSteps(steps, is)
	if err != nil {
		t.Error(err)
	}

	if len(isteps) != 3 {
		t.Errorf("Expected 3 step, got %d", len(isteps))
	}

	x := istep{
		Token: IsInstance,
		Arg:   red,
	}
	s := isteps[1]

	if !reflect.DeepEqual(s, x) {
		t.Errorf("Expected %+v got %+v", x, s)
	}
}

func TestInternalizeHasValueCmd(t *testing.T) {
	is := NewIidStore()
	color := is.Put("color")
	cmd := `Start[iri].HasValue[color,red,blue].Eval`

	steps, err := ParseCommand(cmd)
	if err != nil {
		t.Error(err)
	}

	isteps, err := InternalizeSteps(steps, is)
	if err != nil {
		t.Error(err)
	}

	if len(isteps) != 3 {
		t.Errorf("Expected 3 step, got %d", len(isteps))
	}

	x := istep{
		Token: HasValue,
		Arg:   color,
		Svals: []string{"red", "blue"},
	}
	s := isteps[1]

	if !reflect.DeepEqual(s, x) {
		t.Errorf("Expected %+v got %+v", x, s)
	}
}

func TestInternalizeHasBroader(t *testing.T) {
	is := NewIidStore()
	tname := is.Put("TaxonomyName")
	tnode := is.Put("TaxonomyNodeInstance")
	cmd := `Start[iri].HasBroader[TaxonomyName,TaxonomyNodeInstance].Eval`

	steps, err := ParseCommand(cmd)
	if err != nil {
		t.Error(err)
	}

	isteps, err := InternalizeSteps(steps, is)
	if err != nil {
		t.Error(err)
	}

	if len(isteps) != 3 {
		t.Errorf("Expected 3 step, got %d", len(isteps))
	}

	x := istep{
		Token: HasBroader,
		Arg:   tname,
		Ivals: []Iid{tnode},
	}
	s := isteps[1]

	if !reflect.DeepEqual(s, x) {
		t.Errorf("Expected %+v got %+v", x, s)
	}
}

func TestInternalizeOrCmd(t *testing.T) {
	is := NewIidStore()
	red := is.Put("red")
	blue := is.Put("blue")
	cmd := `Start[iri].Or(IsInstance[red].IsInstance[blue]).Eval`

	steps, err := ParseCommand(cmd)
	if err != nil {
		t.Error(err)
	}

	isteps, err := InternalizeSteps(steps, is)
	if err != nil {
		t.Error(err)
	}

	if len(isteps) != 3 {
		t.Errorf("Expected 3 step, got %d", len(isteps))
	}

	x := istep{
		Token: Or,
		Subcmd: []istep{
			{
				Token: IsInstance,
				Arg:   red,
			},
			{
				Token: IsInstance,
				Arg:   blue,
			},
		},
	}
	s := isteps[1]

	if !reflect.DeepEqual(s, x) {
		t.Errorf("Expected %+v got %+v", x, s)
	}
}

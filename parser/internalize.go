package parser

type Iid uint64

// Internalizer is an interface for a store that maps between
// strings and internal ids.
type Internalizer interface {
	GetIid(string) (Iid, bool)
	GetString(Iid) (string, bool)
	Put(string) Iid
}

type istep struct {
	Token  Token
	Arg    Iid
	Ivals  []Iid
	Svals  []string
	Subcmd []istep
}

// Convert a chain of steps to internalized form that is ready for evaluation.
func InternalizeSteps(chain []Step, is Internalizer) ([]istep, error) {
	steps := make([]istep, 0)
	var step istep

	for _, s := range chain {
		switch s.token {
		case Start, Eval, NoOp, IsActive, IsInactive:
			step = istep{
				Token: s.token,
			}
		case HasType, HasCategory, IsInstance, Follow, FollowInverse, InScheme:
			iarg := is.Put(s.arg)
			step = istep{
				Token: s.token,
				Arg:   iarg,
			}
		case HasValue:
			iarg := is.Put(s.arg)
			step = istep{
				Token: HasValue,
				Arg:   iarg,
				Svals: s.vals,
			}
		case HasBroader:
			iagr := is.Put(s.arg)
			ival := is.Put(s.vals[0])
			step = istep{
				Token: HasBroader,
				Arg:   iagr,
				Ivals: []Iid{ival},
			}
		case Or:
			substeps, err := InternalizeSteps(s.subcmd, is)
			if err != nil {
				return nil, err
			}
			step = istep{
				Token:  Or,
				Subcmd: substeps,
			}
		}
		steps = append(steps, step)
	}
	return steps, nil
}

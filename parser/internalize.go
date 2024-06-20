package parser

type Iid uint64

type Token int

const (
	_ Token = iota
	NoOp
	Start
	Eval
	HasType
	HasCategory
	HasValue
	InScheme
	HasBroader
	IsInstance
	Follow
	FollowInverse
	IsActive
	IsInactive
	Or
)

// Internalizer is an interface for a store that maps between
// strings and internal ids.
type Internalizer interface {
	GetIid(string) (Iid, bool)
	GetString(Iid) (string, bool)
	Put(string) Iid
}

// ASCII string to token
func Atot(s string) Token {
	switch s {
	case "Start":
		return Start
	case "Eval":
		return Eval
	case "HasType":
		return HasType
	case "HasCategory":
		return HasCategory
	case "HasValue":
		return HasValue
	case "InScheme":
		return InScheme
	case "HasBroader":
		return HasBroader
	case "IsInstance":
		return IsInstance
	case "Follow":
		return Follow
	case "FollowInverse":
		return FollowInverse
	case "IsActive":
		return IsActive
	case "IsInactive":
		return IsInactive
	case "Or":
		return Or
	default:
		return 0
	}
}

// Token to ASCII string
func Ttoa(t Token) string {
	switch t {
	case Start:
		return "Start"
	case Eval:
		return "Eval"
	case HasType:
		return "HasType"
	case HasCategory:
		return "HasCategory"
	case HasValue:
		return "HasValue"
	case InScheme:
		return "InScheme"
	case HasBroader:
		return "HasBroader"
	case IsInstance:
		return "IsInstance"
	case Follow:
		return "Follow"
	case FollowInverse:
		return "FollowInverse"
	case IsActive:
		return "IsActive"
	case IsInactive:
		return "IsInactive"
	case Or:
		return "Or"
	case NoOp:
		return "NoOp"
	default:
		return "**error**"
	}
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
		case "Start", "Eval", "NoOp", "IsActive", "IsInactive":
			step = istep{
				Token: Atot(s.token),
			}
		case "HasType", "HasCategory", "IsInstance", "Follow", "FollowInverse", "InScheme":
			iarg := is.Put(s.arg)
			step = istep{
				Token: Atot(s.token),
				Arg:   iarg,
			}
		case "HasValue":
			iarg := is.Put(s.arg)
			step = istep{
				Token: Atot(s.token),
				Arg:   iarg,
				Svals: s.vals,
			}
		case "HasBroader":
			iagr := is.Put(s.arg)
			ival := is.Put(s.vals[0])
			step = istep{
				Token: Atot(s.token),
				Arg:   iagr,
				Ivals: []Iid{ival},
			}
		case "Or":
			substeps, err := InternalizeSteps(s.subcmd, is)
			if err != nil {
				return nil, err
			}
			step = istep{
				Token:  Atot(s.token),
				Subcmd: substeps,
			}
		}
		steps = append(steps, step)
	}
	return steps, nil
}

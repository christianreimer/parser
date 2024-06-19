package parser

import (
	"fmt"
	"strings"
)

type Step struct {
	token  Token
	arg    string
	vals   []string
	subcmd []Step
}

type Token int

const (
	_ Token = iota
	NoOp
	Start
	Eval
	HasType
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

// Parses the full command, including the Start and Eval clauses and
// returns a list of steps to be executed.
func ParseCommand(cmd string) ([]Step, error) {
	cmd, err := cleanCmd(cmd)
	if err != nil {
		return nil, err
	}

	subchain, err := parseSubCommand(cmd, false)
	if err != nil {
		return nil, err
	}

	start := Step{
		token: Start,
		arg:   "iri",
	}
	eval := Step{
		token: Eval,
	}
	chain := make([]Step, 2+len(subchain))
	chain[0] = start
	copy(chain[1:], subchain)
	chain[len(chain)-1] = eval

	return chain, nil
}

func parseSubCommand(cmd string, inor bool) ([]Step, error) {
	chain := make([]Step, 0)

	for {
		step, tail, io, err := parseBody(cmd, inor)
		if err != nil {
			return chain, err
		}

		// Ignore empty steps
		if step.token != NoOp {
			chain = append(chain, step)
		}

		if tail == "" {
			break
		}
		cmd = tail
		inor = io
	}

	return chain, nil
}

// Parses the body of a command, breaking it into steps.
func parseBody(cmd string, inor bool) (Step, string, bool, error) {
	iq := false // inside quotes
	io := inor  // inside or
	is := false // inside square brackets
	buf := make([]rune, 0)

	for i, c := range cmd {
		switch c {
		case '"':
			iq = !iq
		case '[', ']':
			is = !is
		case '.':
			if iq {
				continue
			}
			c := cmd[:i]
			cmd = cmd[i+1:]
			s, io, err := parseStep(c, io)
			return s, cmd, io, err
		case ',':
			if io || is {
				continue
			}
			c := cmd[:i]
			cmd = cmd[i+1:]
			s, io, err := parseStep(c, io)
			return s, cmd, io, err
		case '(':
			if string(buf) != "Or" {
				continue
			}
			_ = cmd[:i]
			i, ok := findNextStandaloneRune(cmd, ')')
			if !ok {
				err := fmt.Errorf("expected Or(step1, ...) got %s", cmd)
				return Step{}, cmd, io, err
			}
			subcmd := cmd[3:i]
			s, err := parseSubCommand(subcmd, true)
			if err != nil {
				return Step{}, cmd, io, err
			}
			cmd = cmd[i+1:]
			// if there is more command after or, then eat the . as well
			if len(cmd) > 1 && cmd[0] == '.' {
				cmd = cmd[1:]
			}
			return Step{
				token:  Or,
				subcmd: s,
			}, cmd, io, err
		default:
			buf = append(buf, c)
		}
	}
	s, io, err := parseStep(cmd, io)
	return s, "", io, err
}

// Parses a single step in the command, checking the arguments and values.
func parseStep(cmd string, inor bool) (Step, bool, error) {
	if cmd == "" {
		return Step{
			token: NoOp,
		}, inor, nil
	}

	if strings.HasPrefix(cmd, "Or") {
		// Or takes a subcommand as its only argument
		cmd := strings.Replace(cmd, "Or(", "", 1)
		i, ok := findNextStandaloneRune(cmd, ')')
		var subcmd string
		if !ok {
			subcmd = cmd
		} else {
			subcmd = cmd[:i]
		}
		clist, err := parseSubCommand(subcmd, true)
		if err != nil {
			return Step{}, inor, fmt.Errorf("failed to parse Or (%s)", err)
		}
		return Step{
			token:  Or,
			subcmd: clist,
		}, inor, nil
	}

	if strings.HasPrefix(cmd, "IsActive") {
		// IsActive takes no arguments
		t, err := parseNoArgStep(cmd)
		if err != nil {
			return Step{}, inor, fmt.Errorf("failed to parse IsActive (%s)", err)
		}
		return Step{
			token: atot(t),
		}, inor, nil
	}

	if strings.HasPrefix(cmd, "IsInactive") {
		// IsInactive takes no arguments
		t, err := parseNoArgStep(cmd)
		if err != nil {
			return Step{}, inor, fmt.Errorf("failed to parse IsInactive (%s)", err)
		}
		return Step{
			token: atot(t),
		}, inor, nil
	}

	if strings.HasPrefix(cmd, "HasType") {
		// HasType takes a single argument which is a type iri
		t, arg, err := parseSingeArgStep(cmd)
		if err != nil {
			return Step{}, inor, fmt.Errorf("failed to parse HasType (%s)", cmd)
		}
		return Step{
			token: atot(t),
			arg:   arg,
		}, inor, nil
	}

	if strings.HasPrefix(cmd, "HasValue") {
		// HasValue takes a field name and a list of values
		t, args, err := parseMultiArgStep(cmd)
		if err != nil {
			return Step{}, inor, fmt.Errorf("failed to parse HasValue (%s)", cmd)
		}
		return Step{
			token: atot(t),
			arg:   args[0],
			vals:  args[1:],
		}, inor, nil
	}

	if strings.HasPrefix(cmd, "InScheme") {
		// InScheme takes a single argument which is a taxonomy iri
		t, arg, err := parseSingeArgStep(cmd)
		if err != nil {
			return Step{}, inor, fmt.Errorf("failed to parse InScheme (%s)", cmd)
		}
		return Step{
			token: atot(t),
			arg:   arg,
		}, inor, nil
	}

	if strings.HasPrefix(cmd, "HasBroader") {
		// HasBroader takes two arguments, the first is the taxonomy iri the second
		// is the tatget node
		t, args, err := parseMultiArgStep(cmd)
		if err != nil {
			return Step{}, inor, fmt.Errorf("failed to parse HasBroader (%s)", cmd)
		}
		if len(args) != 2 {
			return Step{}, inor, fmt.Errorf("expected HasBroader[taxonomy, target] got %s", cmd)
		}
		return Step{
			token: atot(t),
			arg:   args[0],
			vals:  []string{args[1]},
		}, inor, nil
	}

	if strings.HasPrefix(cmd, "IsInstance") {
		// IsInstance takes a single argument which is the instance iri
		t, arg, err := parseSingeArgStep(cmd)
		if err != nil {
			return Step{}, inor, fmt.Errorf("failed to parse IsInstance (%s)", cmd)
		}
		return Step{
			token: atot(t),
			arg:   arg,
		}, inor, nil
	}

	if strings.HasPrefix(cmd, "Follow") || strings.HasPrefix(cmd, "FollowInverse") {
		// Follow and FollowInverse both take a single argument which is the
		// relationship iri
		t, arg, err := parseSingeArgStep(cmd)
		if err != nil {
			return Step{}, inor, fmt.Errorf("failed to parse Follow (%s)", cmd)
		}
		return Step{
			token: atot(t),
			arg:   arg,
		}, inor, nil
	}

	return Step{}, inor, nil
}

// Parse a step with no arguments.
func parseNoArgStep(cmd string) (string, error) {
	token, arg := splitOnRune(cmd, '[')

	if arg != "]" {
		return "", fmt.Errorf("expected Token[] got %s", cmd)
	}

	return token, nil
}

// Parse a step with a single argument.
func parseSingeArgStep(cmd string) (string, string, error) {
	token, arg := splitOnRune(cmd, '[')

	i, ok := findNext(arg, ']')
	if !ok {
		return "", "", fmt.Errorf("expected Token[arg] got %s", cmd)
	}
	arg = removeOuterQuotes(arg[:i])
	return token, arg, nil
}

// Parse a step that can take multiple argument.
func parseMultiArgStep(cmd string) (string, []string, error) {
	token, args := splitOnRune(cmd, '[')

	i, ok := findNext(args, ']')
	if !ok {
		return "", nil, fmt.Errorf("expected Token[arg,...] got %s", cmd)
	}
	args = args[:i]

	// Split args on commas
	argl := make([]string, 0)
	var arg string
	for {
		i, ok := findNext(args, ',')
		if !ok && len(args) > 0 {
			// want to get the last argument in the list
			arg = args
		} else {
			arg = args[:i]
			args = args[i+1:]
		}
		arg = removeOuterQuotes(arg)
		argl = append(argl, arg)
		if !ok {
			break
		}
	}

	return token, argl, nil
}

func cleanCmd(cmd string) (string, error) {
	cmd = strings.TrimSpace(cmd)
	cmd = convertIris(cmd)
	cmd = removeWhiteSpace(cmd, nil)

	// Strip the And() clauses from the command since it is just syntatic sugar.
	more := false
	for {
		cmd, more = removeAnd(cmd)
		if !more {
			break
		}
	}

	// Validate Start[iri] and Eval clauses
	if !strings.HasPrefix(cmd, "Start[iri].") {
		return "", fmt.Errorf("invalid cmd, must begin with Start[iri] got %s", cmd)
	}
	cmd = strings.Replace(cmd, "Start[iri]", "", 1)

	if !strings.HasSuffix(cmd, "Eval") {
		return "", fmt.Errorf("invalid cmd, must end with Eval got %s", cmd)
	}
	cmd = strings.Replace(cmd, "Eval", "", 1)

	// Strip outer . if there are any
	if len(cmd) > 0 && cmd[0] == '.' {
		cmd = cmd[1:]
	}
	if len(cmd) > 0 && cmd[len(cmd)-1] == '.' {
		cmd = cmd[:len(cmd)-1]
	}

	return cmd, nil
}

// Converts the IRIs to qnames.
func convertIris(cmd string) string {
	cmd = strings.ReplaceAll(cmd, "<", "")
	cmd = strings.ReplaceAll(cmd, ">", "")
	cmd = strings.ReplaceAll(cmd, "https://bsm.bloomberg.com/ontology/", "bsm:")
	cmd = strings.ReplaceAll(cmd, "https://bsm.bloomberg.com/instance/", "bsi:")
	return cmd
}

// Removes the And() clause from the command.
func removeAnd(cmd string) (string, bool) {
	if !strings.Contains(cmd, "And(") {
		return cmd, false
	}

	iq := false
	buf := make([]rune, 0, len(cmd))
	var start, end int

	for i, c := range cmd {
		switch c {
		case '"':
			iq = !iq
		case '(':
			buf3 := string(buf[len(buf)-3:])
			if !iq && buf3 == "And" {
				// We have found the start of the .And() clause
				start = i
			}
		case ')':
			if !iq && start > 0 {
				// We have found the end of the .And() clause
				end = i
				s1 := cmd[:start-3]
				s2 := cmd[start+1 : end]
				s3 := cmd[end+1:]
				return s1 + s2 + s3, true
			}
		default:
			buf = append(buf, c)
		}
	}
	return cmd, false
}

func removeOuterQuotes(cmd string) string {
	if len(cmd) < 2 {
		return cmd
	}
	if cmd[0] == '"' && cmd[len(cmd)-1] == '"' {
		return cmd[1 : len(cmd)-1]
	}
	return cmd
}

func removeWhiteSpace(cmd string, term *rune) string {
	buf := make([]rune, 0)
	iq := false
	var prev rune

	for _, r := range cmd {
		if term != nil && *term == r && prev != '\\' {
			break
		}

		switch r {
		case '"':
			if prev != '\\' {
				iq = !iq
			}
			buf = append(buf, r)
		case ' ':
			if iq {
				buf = append(buf, r)
			}
		case '\n', '\t':
			// eat the rune
		default:
			buf = append(buf, r)
		}
		prev = r
	}

	return string(buf)
}

// Splits the cmd by the first occurences of the term, accounting for quotes.
func splitOnRune(cmd string, term rune) (string, string) {
	iq := false
	for i, c := range cmd {
		if c == '"' {
			iq = !iq
		}
		if term == c && !iq {
			return cmd[:i], cmd[i+1:]
		}
	}
	return "", ""
}

// Finds the next occurance of rune in the command that is outside
// of quotes or square brackers.
func findNextStandaloneRune(cmd string, r rune) (int, bool) {
	iq := false
	ib := false
	for i, c := range cmd {
		if c == '"' {
			iq = !iq
		}
		if c == '[' || c == ']' {
			ib = !ib
		}
		if c == r && !iq && !ib {
			return i, true
		}
	}
	return 0, false
}

// Finds the next occurance of term in the command, accounting for quotes.
func findNext(cmd string, term rune) (int, bool) {
	iq := false
	var prev rune
	for i, r := range cmd {
		switch r {
		case '"':
			if prev != '\\' {
				iq = !iq
			}
		case term:
			if !iq {
				return i, true
			}
		default:
			prev = r
		}
	}
	return 0, false
}

// String (ASCII)to token
func atot(s string) Token {
	switch s {
	case "Start":
		return Start
	case "Eval":
		return Eval
	case "HasType":
		return HasType
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
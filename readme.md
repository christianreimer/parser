![.](bremlin.jpeg)

## Bremlin Parser

```
Enter Bremlin command:
Start[iri]
.Or(
        HasType[Gremlin]
        .HasType[GooGrok]
)
.HasValue[FurColor, "green", "blue"]
.And(
        InScheme[<http://example.org/Animals>]
        .HasBroader[<http://example.org/Fantasy>, <http://example.org/Preditor>]
)
.Follow[SmellOfFood]
.HasType[TastyMeal]
.Eval

[]parser.Step{
  parser.Step{
    token:  "Start",
    arg:    "iri",
    vals:   []string{},
    subcmd: []parser.Step{},
  },
  parser.Step{
    token:  "Or",
    arg:    "",
    vals:   []string{},
    subcmd: []parser.Step{
      parser.Step{
        token:  "HasType",
        arg:    "Gremlin",
        vals:   []string{},
        subcmd: []parser.Step{},
      },
      parser.Step{
        token:  "HasType",
        arg:    "GooGrok",
        vals:   []string{},
        subcmd: []parser.Step{},
      },
    },
  },
  parser.Step{
    token: "HasValue",
    arg:   "FurColor",
    vals:  []string{
      "green",
      "blue",
    },
    subcmd: []parser.Step{},
  },
  parser.Step{
    token:  "InScheme",
    arg:    "ex:Animals",
    vals:   []string{},
    subcmd: []parser.Step{},
  },
  parser.Step{
    token: "HasBroader",
    arg:   "ex:Fantasy",
    vals:  []string{
      "ex:Preditor",
    },
    subcmd: []parser.Step{},
  },
  parser.Step{
    token:  "Follow",
    arg:    "SmellOfFood",
    vals:   []string{},
    subcmd: []parser.Step{},
  },
  parser.Step{
    token:  "HasType",
    arg:    "TastyMeal",
    vals:   []string{},
    subcmd: []parser.Step{},
  },
  parser.Step{
    token:  "Eval",
    arg:    "",
    vals:   []string{},
    subcmd: []parser.Step{},
  },
}
```
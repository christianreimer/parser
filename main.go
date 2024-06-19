package main

import (
	"bremlin/parser"
	"log"

	"github.com/k0kubun/pp/v3"
)

func main() {
	cmd := `
		Start[iri]
			.HasType[<https://bsm.bloomberg.com/ontology/Company>]
			.HasValue[field1, "value1", "3.14"]
			.Or(
				HasType[TypeOr1]. HasType[TypeOr2]
			)
			.And(HasType[TypeAnd1].HasType[TypeAnd2])
			.Eval`

	chain, err := parser.ParseCommand(cmd)
	if err != nil {
		log.Fatal(err)
	}
	pp.Println(chain)
}

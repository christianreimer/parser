package main

import (
	"bremlin/parser"
	"fmt"
	"log"

	"github.com/k0kubun/pp/v3"
	"github.com/renstrom/dedent"
)

func main() {
	cmd := `
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
		.Eval`

	fmt.Printf("Command:%s\n\n", dedent.Dedent(cmd))

	chain, err := parser.ParseCommand(cmd)
	if err != nil {
		log.Fatal(err)
	}
	pp.Printf("Step Chain:\n%s\n", chain)
}

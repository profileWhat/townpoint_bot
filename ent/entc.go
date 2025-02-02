//go:build ignore
// +build ignore

package main

import (
	"log"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

func main() {
	if err := entc.Generate("./ent/schema", &gen.Config{
		Features: []gen.Feature{
			gen.FeatureVersionedMigration,
			gen.FeatureIntercept,
			gen.FeatureUpsert,
		},
		Target:  "./ent/generated",
		Package: "townpoint_bot/ent/generated",
	}); err != nil {
		log.Fatalf("running ent codegen: %v", err)
	}
}

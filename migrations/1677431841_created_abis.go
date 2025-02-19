package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection := core.NewBaseCollection("abis")
		collection.Id = "dkzs3rh94mj101j"

		collection.Fields.Add(
			&core.NumberField{
				Id:       "f2kljlbw",
				Name:     "chainId",
				Required: true,
			},
			&core.TextField{
				Id:       "gvetwq8a",
				Name:     "address",
				Required: true,
				Pattern:  "",
			},
			&core.JSONField{
				Id:       "qs6mngfo",
				Name:     "abi",
				Required: true,
			},
		)

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("dkzs3rh94mj101j")
		if err != nil {
			return err
		}

		return app.Delete(collection)
	})
}

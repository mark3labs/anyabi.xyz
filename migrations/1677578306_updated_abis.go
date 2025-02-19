package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("dkzs3rh94mj101j")
		if err != nil {
			return err
		}

		// add
		collection.Fields.Add(&core.TextField{
			System:   false,
			Id:       "ftifmr9g",
			Name:     "name",
			Required: true,
			Pattern:  "",
		})

		return app.Save(collection)
	}, func(app core.App) error {
		// Use raw SQL to remove the field
		_, err := app.DB().NewQuery("ALTER TABLE abis DROP COLUMN name").Execute()
		return err
	})
}

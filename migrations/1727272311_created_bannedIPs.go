package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection := core.NewBaseCollection("bannedIPs")
		collection.Id = "ggalodjt1jwro34"
		
		collection.Fields.Add(
			&core.TextField{
				Id:       "pyv6mkxb",
				Name:     "ip",
				Required: true,
				Pattern:  "",
			},
		)

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("ggalodjt1jwro34")
		if err != nil {
			return err
		}

		return app.Delete(collection)
	})
}

package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		jsonData := `{
			"id": "ggalodjt1jwro34",
			"created": "2024-09-25 13:51:51.754Z",
			"updated": "2024-09-25 13:51:51.754Z",
			"name": "bannedIPs",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "pyv6mkxb",
					"name": "ip",
					"type": "text",
					"required": true,
					"unique": true,
					"options": {
						"min": null,
						"max": null,
						"pattern": ""
					}
				}
			],
			"listRule": null,
			"viewRule": null,
			"createRule": null,
			"updateRule": null,
			"deleteRule": null,
			"options": {}
		}`

		collection := &models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return daos.New(db).SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db)

		collection, err := dao.FindCollectionByNameOrId("ggalodjt1jwro34")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}

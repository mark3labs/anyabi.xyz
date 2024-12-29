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
			"id": "dkzs3rh94mj101j",
			"created": "2023-02-26 17:17:21.763Z",
			"updated": "2023-02-26 17:17:21.763Z",
			"name": "abis",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "f2kljlbw",
					"name": "chainId",
					"type": "number",
					"required": true,
					"unique": false,
					"options": {
						"min": null,
						"max": null
					}
				},
				{
					"system": false,
					"id": "gvetwq8a",
					"name": "address",
					"type": "text",
					"required": true,
					"unique": false,
					"options": {
						"min": null,
						"max": null,
						"pattern": ""
					}
				},
				{
					"system": false,
					"id": "qs6mngfo",
					"name": "abi",
					"type": "json",
					"required": true,
					"unique": false,
					"options": {}
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

		collection, err := dao.FindCollectionByNameOrId("dkzs3rh94mj101j")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}

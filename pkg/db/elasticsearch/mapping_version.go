package elasticsearch

var (
	// IndexMappingV1 索引映射v1历史版本
	IndexMappingV1 = `{
				"mappings":{
					"properties":{
						"key":{
							"type":"keyword"
						},
						"version":{
							"type":"keyword"
						},
						"user_id":{
							"type":"keyword"
						}
					}
				}
			}`

	// IndexMappingLatest 索引最新版本
	IndexMappingLatest = `{
				"mappings":{
					"properties":{
						"key":{
							"type":"keyword"
						},
						"version":{
							"type":"keyword"
						},
						"user_id":{
							"type":"keyword"
						},
						"data_type":{
							"type":"keyword"
						},
						"name":{
							"type":"keyword"
						},
						"id":{
							"type":"keyword"
						},
						"tag":{
							"type":"keyword"
						},
						"content":{
							"type":"keyword"
						},
						"created_at":{
							"type":"float"
						}
					}
				}
			}`
)

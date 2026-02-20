package handlers

import (
	"encoding/json"
	"net/http"
)

func SwaggerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	swagger := map[string]interface{}{
		"openapi": "3.0.0",
		"info": map[string]string{
			"title":       "Task API",
			"description": "Simple task management API",
			"version":     "1.0.0",
		},
		"servers": []map[string]string{
			{"url": "http://localhost:8080"},
		},
		"paths": map[string]interface{}{
			"/v1/tasks": map[string]interface{}{
				"get": map[string]interface{}{
					"summary":     "Get all tasks or task by ID",
					"description": "Returns all tasks or a specific task if id parameter is provided",
					"parameters": []map[string]interface{}{
						{
							"name":        "id",
							"in":          "query",
							"description": "Task ID",
							"required":    false,
							"schema":      map[string]string{"type": "integer"},
						},
						{
							"name":        "done",
							"in":          "query",
							"description": "Filter by done status",
							"required":    false,
							"schema":      map[string]string{"type": "boolean"},
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Success",
						},
						"400": map[string]interface{}{
							"description": "Bad Request",
						},
						"404": map[string]interface{}{
							"description": "Not Found",
						},
					},
					"security": []map[string][]string{
						{"ApiKeyAuth": {}},
					},
				},
				"post": map[string]interface{}{
					"summary":     "Create a new task",
					"description": "Creates a new task with the provided title",
					"requestBody": map[string]interface{}{
						"required": true,
						"content": map[string]interface{}{
							"application/json": map[string]interface{}{
								"schema": map[string]interface{}{
									"type": "object",
									"properties": map[string]interface{}{
										"title": map[string]string{"type": "string"},
									},
									"required": []string{"title"},
								},
							},
						},
					},
					"responses": map[string]interface{}{
						"201": map[string]interface{}{
							"description": "Created",
						},
						"400": map[string]interface{}{
							"description": "Bad Request",
						},
					},
					"security": []map[string][]string{
						{"ApiKeyAuth": {}},
					},
				},
				"patch": map[string]interface{}{
					"summary":     "Update task status",
					"description": "Updates the done status of a task",
					"parameters": []map[string]interface{}{
						{
							"name":        "id",
							"in":          "query",
							"description": "Task ID",
							"required":    true,
							"schema":      map[string]string{"type": "integer"},
						},
					},
					"requestBody": map[string]interface{}{
						"required": true,
						"content": map[string]interface{}{
							"application/json": map[string]interface{}{
								"schema": map[string]interface{}{
									"type": "object",
									"properties": map[string]interface{}{
										"done": map[string]string{"type": "boolean"},
									},
									"required": []string{"done"},
								},
							},
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Success",
						},
						"400": map[string]interface{}{
							"description": "Bad Request",
						},
						"404": map[string]interface{}{
							"description": "Not Found",
						},
					},
					"security": []map[string][]string{
						{"ApiKeyAuth": {}},
					},
				},
				"delete": map[string]interface{}{
					"summary":     "Delete a task",
					"description": "Deletes a task by ID",
					"parameters": []map[string]interface{}{
						{
							"name":        "id",
							"in":          "query",
							"description": "Task ID",
							"required":    true,
							"schema":      map[string]string{"type": "integer"},
						},
					},
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Success",
						},
						"400": map[string]interface{}{
							"description": "Bad Request",
						},
						"404": map[string]interface{}{
							"description": "Not Found",
						},
					},
					"security": []map[string][]string{
						{"ApiKeyAuth": {}},
					},
				},
			},
		},
		"components": map[string]interface{}{
			"securitySchemes": map[string]interface{}{
				"ApiKeyAuth": map[string]string{
					"type": "apiKey",
					"in":   "header",
					"name": "X-API-KEY",
				},
			},
		},
	}

	json.NewEncoder(w).Encode(swagger)
}

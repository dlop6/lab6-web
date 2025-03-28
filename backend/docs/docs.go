// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/matches": {
            "get": {
                "description": "Devuelve la lista completa de partidos",
                "produces": [
                    "application/json"
                ],
                "summary": "Obtener todos los partidos",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/main.Match"
                            }
                        }
                    }
                }
            }
        },
        "/api/matches/{id}/extratime": {
            "patch": {
                "description": "Establece minutos de tiempo extra",
                "produces": [
                    "application/json"
                ],
                "summary": "Actualizar tiempo extra",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID del partido",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Minutos de tiempo extra",
                        "name": "minutes",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "integer"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.Match"
                        }
                    }
                }
            }
        },
        "/api/matches/{id}/goals": {
            "patch": {
                "description": "Incrementa los goles de un equipo",
                "produces": [
                    "application/json"
                ],
                "summary": "Actualizar goles",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID del partido",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Equipo (team1 o team2)",
                        "name": "team",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "Goles a añadir",
                        "name": "goals",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "integer"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.Match"
                        }
                    }
                }
            }
        },
        "/api/matches/{id}/redcards": {
            "patch": {
                "description": "Incrementa el contador de tarjetas rojas",
                "produces": [
                    "application/json"
                ],
                "summary": "Registrar tarjeta roja",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID del partido",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.Match"
                        }
                    }
                }
            }
        },
        "/api/matches/{id}/yellowcards": {
            "patch": {
                "description": "Incrementa el contador de tarjetas amarillas",
                "produces": [
                    "application/json"
                ],
                "summary": "Registrar tarjeta amarilla",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID del partido",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.Match"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "main.Match": {
            "description": "Información detallada sobre un partido de La Liga",
            "type": "object",
            "required": [
                "date",
                "team1",
                "team2"
            ],
            "properties": {
                "date": {
                    "type": "string"
                },
                "extraTime": {
                    "type": "integer"
                },
                "id": {
                    "type": "string"
                },
                "redCards": {
                    "type": "integer"
                },
                "score1": {
                    "type": "integer"
                },
                "score2": {
                    "type": "integer"
                },
                "team1": {
                    "type": "string"
                },
                "team2": {
                    "type": "string"
                },
                "yellowCards": {
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "La Liga Tracker API",
	Description:      "API para gestión de partidos de fútbol",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}

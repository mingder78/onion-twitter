package main

import "html/template"

var html = template.Must(template.New("swagger").Parse(`
{
    "swagger": "2.0",
    "info": {
        "title": "{{.lowerCase}} API",
        "description": "Auto generated API from JSON",
        "version": "1.0.0"
    },
    "host": "{{.url}}",
    "chemes": [
        "http"
    ],
    "basePath": "/",
    "produces": [
        "application/json"
    ],
    "paths": {
        "/v2/doc/swagger.json": {
            "get": {
                "summary": "Get swagger json file",
                "description": "",
                "responses": {
                    "200": null,
                    "description": "A json file for swagger api document"
                }
            }
        },
        "/{{.lowerCase}}": {
            "get": {
                "summary": "Get all {{.lowerCase}}",
                "description": "",
                "responses": {
                    "200": {
                        "description": "An array of {{.lowerCase}}",
                        "schema": {
                            "type": "array",
                            "$ref": "#/definitions/{{.lowerCase}}"
                        }
                    }
                }
            },
            "post": {
                "parameters": [
                    {
                        "name": "{{.lowerCase}}",
                        "in": "body",
                        "description": "The {{.lowerCase}} JSON you want to post",
                        "schema": {
                            "$ref": "#/definitions/{{.lowerCase}}"
                        },
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Make a new {{.lowerCase}}",
                        "schema": {
                            "$ref": "#/definitions/{{.lowerCase}}Record"
                        }
                    }
                }
            }
        },
        "/{{.lowerCase}}/{id}": {
            "get": {
                "summary": "Get {{.lowerCase}} from Id",
                "parameters": [
                    {
                        "name": "id",
                        "in": "path",
                        "type": "string",
                        "description": "ID of the {{.lowerCase}}",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Get the {{.lowerCase}} with Id",
                        "schema": {
                            "$ref": "#/definitions/{{.lowerCase}}Record"
                        }
                    }
                }
            },
            "delete": {
                "summary": "Remove a {{.lowerCase}}",
                "description": "",
                "parameters": [
                    {
                        "name": "id",
                        "in": "path",
                        "type": "string",
                        "description": "ID of the {{.lowerCase}}",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "deleted"
                    }
                }
            },
            "put": {
                "summary": "Update a {{.lowerCase}}",
                "description": "",
                "parameters": [
                    {
                        "name": "id",
                        "in": "path",
                        "type": "string",
                        "description": "ID of the {{.lowerCase}}",
                        "required": true
                    },
                    {
                        "name": "{{.lowerCase}}",
                        "in": "body",
                        "description": "The {{.lowerCase}} JSON you want to post",
                        "schema": {
                            "$ref": "#/definitions/{{.lowerCase}}"
                        },
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Updated {{.lowerCase}}",
                        "schema": {
                            "$ref": "#/definitions/{{.lowerCase}}Record"
                        }
                    }
                }
            },
            "patch": {
                "summary": "Update a {{.lowerCase}}",
                "description": "",
                "parameters": [
                    {
                        "name": "id",
                        "in": "path",
                        "type": "string",
                        "description": "ID of the {{.lowerCase}}",
                        "required": true
                    },
                    {
                        "name": "{{.lowerCase}}",
                        "in": "body",
                        "description": "The {{.lowerCase}} JSON you want to post",
                        "schema": {
                            "$ref": "#/definitions/{{.lowerCase}}"
                        },
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Updated {{.lowerCase}}",
                        "schema": {
                            "$ref": "#/definitions/{{.lowerCase}}Record"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "{{.lowerCase}}": {
            "type": "object",
            "properties": {
{{.dataType}}
            }
        },
        "{{.lowerCase}}Record": {
            "type": "object",
            "properties": {
                "ginger_created": {
                    "type": "string"
                },
                "ginger_id": {
                    "type": "string"
                },
{{.dataType}}
            }
        }
    }
}
`))

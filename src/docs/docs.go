// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2020-04-07 17:49:19.3659571 +0200 CEST m=+0.072959701

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "email": "ad5880mu-s@student.lu.se"
        },
        "license": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/company": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "Gets all CompanyPublic from database",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/db.CompanyPublic"
                            }
                        }
                    }
                }
            },
            "post": {
                "produces": [
                    "application/json"
                ],
                "summary": "Adds a company to the database",
                "parameters": [
                    {
                        "description": "Company",
                        "name": "company",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/db.Company"
                        }
                    }
                ],
                "responses": {
                    "200": {}
                }
            },
            "patch": {
                "produces": [
                    "application/json"
                ],
                "summary": "Updates a company in the database, then returns the updated company. Requires authorization.",
                "parameters": [
                    {
                        "description": "Company",
                        "name": "company",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/db.Company"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/db.Company"
                        }
                    }
                }
            }
        },
        "/company/code/{ticketCode}/verify": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "summary": "Verifies a ticket code for a company. Requires authorization.",
                "parameters": [
                    {
                        "description": "Company",
                        "name": "company",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/db.Company"
                        }
                    },
                    {
                        "type": "string",
                        "description": "Code",
                        "name": "ticketCode",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Ticket was verified."
                    },
                    "401": {
                        "description": "Ticket could not be verified."
                    }
                }
            }
        },
        "/company/distance": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "summary": "Gets companies within a certain distance.",
                "parameters": [
                    {
                        "description": "Distance",
                        "name": "distance",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/db.Distance"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/db.Company"
                            }
                        }
                    }
                }
            }
        },
        "/company/info": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "Gets a full company by id, no password required. Requires authorization.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Token",
                        "name": "authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/db.Company"
                        }
                    }
                }
            }
        },
        "/company/login": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "summary": "Takes in a company as parameter, looks for password hash in database.",
                "parameters": [
                    {
                        "description": "Company",
                        "name": "company",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/db.Company"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Returns \"success\" and token as body. Also sets cookie with token."
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            }
        },
        "/company/slots": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "Gets all slots for a certain company. Requires authorization.",
                "parameters": [
                    {
                        "description": "Company",
                        "name": "company",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/db.Company"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/db.Slot"
                            }
                        }
                    }
                }
            },
            "post": {
                "produces": [
                    "application/json"
                ],
                "summary": "Adds slots to database. Requires authorization.",
                "parameters": [
                    {
                        "description": "Slots",
                        "name": "slots",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/db.Slot"
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {}
                }
            },
            "patch": {
                "produces": [
                    "application/json"
                ],
                "summary": "Updates a certain slot, then returns updated slot. Requires authorization.",
                "parameters": [
                    {
                        "description": "Slot",
                        "name": "slot",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/db.Slot"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/db.Slot"
                        }
                    }
                }
            }
        },
        "/company/slots/id": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "Gets a full slot. Requires a slot as parameter, but an id in body will suffice.",
                "parameters": [
                    {
                        "description": "Slot",
                        "name": "slot",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/db.Slot"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/db.Slot"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "db.Company": {
            "type": "object",
            "properties": {
                "adress": {
                    "type": "string"
                },
                "city": {
                    "type": "string"
                },
                "contact_firstname": {
                    "type": "string"
                },
                "contact_lastname": {
                    "type": "string"
                },
                "contact_number": {
                    "type": "string"
                },
                "country": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "latitude": {
                    "type": "string"
                },
                "longitude": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "post_code": {
                    "type": "string"
                },
                "verified": {
                    "type": "string"
                }
            }
        },
        "db.CompanyPublic": {
            "type": "object",
            "properties": {
                "adress": {
                    "type": "string"
                },
                "city": {
                    "type": "string"
                },
                "country": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "latitude": {
                    "type": "string"
                },
                "longitude": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "post_code": {
                    "type": "string"
                }
            }
        },
        "db.Distance": {
            "type": "object",
            "properties": {
                "distance": {
                    "type": "integer"
                },
                "latMax": {
                    "type": "number"
                },
                "latMin": {
                    "type": "number"
                },
                "latitude": {
                    "type": "number"
                },
                "lonMax": {
                    "type": "number"
                },
                "lonMin": {
                    "type": "number"
                },
                "longitude": {
                    "type": "number"
                },
                "r": {
                    "type": "number"
                }
            }
        },
        "db.Slot": {
            "type": "object",
            "properties": {
                "company_id": {
                    "type": "string"
                },
                "day": {
                    "type": "string"
                },
                "end_time": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "max": {
                    "type": "string"
                },
                "start_time": {
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "",
	BasePath:    "/",
	Schemes:     []string{},
	Title:       "ShopAlone API",
	Description: "Swagger API for ShopAlone API",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}

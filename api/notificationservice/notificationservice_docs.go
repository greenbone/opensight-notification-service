// Package notificationservice Code generated by swaggo/swag. DO NOT EDIT
package notificationservice

import "github.com/swaggo/swag"

const docTemplatenotificationservice = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "license": {
            "name": "AGPL-3.0-or-later"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/notifications": {
            "put": {
                "security": [
                    {
                        "KeycloakAuth": []
                    }
                ],
                "description": "Returns a list of notifications matching the provided filters",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "notification"
                ],
                "summary": "List Notifications",
                "parameters": [
                    {
                        "type": "string",
                        "example": "Bearer eyJhbGciOiJSUzI1NiIs",
                        "description": "Authentication header",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "filters, paging and sorting",
                        "name": "MatchCriterias",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/query.ResultSelector"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/query.ResponseListWithMetadata-models_Notification"
                        },
                        "headers": {
                            "api-version": {
                                "type": "string",
                                "description": "API version"
                            }
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "KeycloakAuth": []
                    }
                ],
                "description": "Create a new notification",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "notification"
                ],
                "summary": "Create Notification",
                "parameters": [
                    {
                        "type": "string",
                        "example": "Bearer eyJhbGciOiJSUzI1NiIs",
                        "description": "Authentication header",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "notification to add",
                        "name": "Notification",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Notification"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/query.ResponseWithMetadata-models_Notification"
                        },
                        "headers": {
                            "api-version": {
                                "type": "string",
                                "description": "API version"
                            }
                        }
                    }
                }
            }
        },
        "/notifications/options": {
            "get": {
                "security": [
                    {
                        "KeycloakAuth": []
                    }
                ],
                "description": "Get filter options for listing notifications",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "notification"
                ],
                "summary": "Notification filter options",
                "parameters": [
                    {
                        "type": "string",
                        "example": "Bearer eyJhbGciOiJSUzI1NiIs",
                        "description": "Authentication header",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/query.ResponseWithMetadata-array_query_FilterOption"
                        },
                        "headers": {
                            "api-version": {
                                "type": "string",
                                "description": "API version"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "filter.CompareOperator": {
            "type": "string",
            "enum": [
                "beginsWith",
                "doesNotBeginWith",
                "contains",
                "doesNotContain",
                "isNumberEqualTo",
                "isEqualTo",
                "isIpEqualTo",
                "isStringEqualTo",
                "isStringCaseInsensitiveEqualTo",
                "isNotEqualTo",
                "isNumberNotEqualTo",
                "isIpNotEqualTo",
                "isStringNotEqualTo",
                "isGreaterThan",
                "isGreaterThanOrEqualTo",
                "isLessThan",
                "isLessThanOrEqualTo",
                "beforeDate",
                "afterDate",
                "exists",
                "isEqualToRating",
                "isNotEqualToRating",
                "isGreaterThanRating",
                "isLessThanRating",
                "isGreaterThanOrEqualToRating",
                "isLessThanOrEqualToRating",
                "betweenDates"
            ],
            "x-enum-varnames": [
                "CompareOperatorBeginsWith",
                "CompareOperatorDoesNotBeginWith",
                "CompareOperatorContains",
                "CompareOperatorDoesNotContain",
                "CompareOperatorIsNumberEqualTo",
                "CompareOperatorIsEqualTo",
                "CompareOperatorIsIpEqualTo",
                "CompareOperatorIsStringEqualTo",
                "CompareOperatorIsStringCaseInsensitiveEqualTo",
                "CompareOperatorIsNotEqualTo",
                "CompareOperatorIsNumberNotEqualTo",
                "CompareOperatorIsIpNotEqualTo",
                "CompareOperatorIsStringNotEqualTo",
                "CompareOperatorIsGreaterThan",
                "CompareOperatorIsGreaterThanOrEqualTo",
                "CompareOperatorIsLessThan",
                "CompareOperatorIsLessThanOrEqualTo",
                "CompareOperatorBeforeDate",
                "CompareOperatorAfterDate",
                "CompareOperatorExists",
                "CompareOperatorIsEqualToRating",
                "CompareOperatorIsNotEqualToRating",
                "CompareOperatorIsGreaterThanRating",
                "CompareOperatorIsLessThanRating",
                "CompareOperatorIsGreaterThanOrEqualToRating",
                "CompareOperatorIsLessThanOrEqualToRating",
                "CompareOperatorBetweenDates"
            ]
        },
        "filter.ControlType": {
            "type": "string",
            "enum": [
                "bool",
                "enum",
                "float",
                "integer",
                "string",
                "dateTime",
                "uuid",
                "autocomplete"
            ],
            "x-enum-varnames": [
                "ControlTypeBool",
                "ControlTypeEnum",
                "ControlTypeFloat",
                "ControlTypeInteger",
                "ControlTypeString",
                "ControlTypeDateTime",
                "ControlTypeUuid",
                "ControlTypeAutocomplete"
            ]
        },
        "filter.LogicOperator": {
            "type": "string",
            "enum": [
                "and",
                "or"
            ],
            "x-enum-varnames": [
                "LogicOperatorAnd",
                "LogicOperatorOr"
            ]
        },
        "filter.ReadableValue-string": {
            "type": "object",
            "properties": {
                "label": {
                    "description": "Label is the human-readable form of the value",
                    "type": "string"
                },
                "value": {
                    "description": "Value is the value for the backend",
                    "type": "string"
                }
            }
        },
        "filter.Request": {
            "type": "object",
            "required": [
                "operator"
            ],
            "properties": {
                "fields": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/filter.RequestField"
                    }
                },
                "operator": {
                    "$ref": "#/definitions/filter.LogicOperator"
                }
            }
        },
        "filter.RequestField": {
            "type": "object",
            "required": [
                "name",
                "operator",
                "value"
            ],
            "properties": {
                "keys": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "name": {
                    "type": "string"
                },
                "operator": {
                    "$ref": "#/definitions/filter.CompareOperator"
                },
                "value": {
                    "description": "Value can be a list of values or a value"
                }
            }
        },
        "filter.RequestOptionType": {
            "type": "object",
            "properties": {
                "type": {
                    "enum": [
                        "string",
                        "float",
                        "integer",
                        "enum",
                        "bool"
                    ],
                    "allOf": [
                        {
                            "$ref": "#/definitions/filter.ControlType"
                        }
                    ]
                }
            }
        },
        "models.Notification": {
            "type": "object",
            "required": [
                "detail",
                "level",
                "origin",
                "timestamp",
                "title"
            ],
            "properties": {
                "customFields": {
                    "description": "can contain arbitrary structured information about the notification",
                    "type": "object",
                    "additionalProperties": {}
                },
                "detail": {
                    "type": "string"
                },
                "id": {
                    "type": "string",
                    "readOnly": true
                },
                "level": {
                    "type": "string",
                    "enum": [
                        "info",
                        "warning",
                        "error"
                    ]
                },
                "origin": {
                    "type": "string"
                },
                "originUri": {
                    "description": "can be used to provide a link to the origin",
                    "type": "string"
                },
                "timestamp": {
                    "type": "string",
                    "format": "date-time"
                },
                "title": {
                    "description": "can also be seen as the 'type'",
                    "type": "string"
                }
            }
        },
        "paging.Request": {
            "type": "object",
            "properties": {
                "index": {
                    "type": "integer"
                },
                "size": {
                    "type": "integer"
                }
            }
        },
        "paging.Response": {
            "type": "object",
            "required": [
                "index",
                "size",
                "totalDisplayableResults"
            ],
            "properties": {
                "index": {
                    "type": "integer"
                },
                "size": {
                    "type": "integer"
                },
                "totalDisplayableResults": {
                    "type": "integer"
                },
                "totalResults": {
                    "type": "integer"
                }
            }
        },
        "query.FilterOption": {
            "type": "object",
            "required": [
                "control",
                "multiSelect",
                "name",
                "operators"
            ],
            "properties": {
                "control": {
                    "$ref": "#/definitions/filter.RequestOptionType"
                },
                "multiSelect": {
                    "type": "boolean"
                },
                "name": {
                    "$ref": "#/definitions/filter.ReadableValue-string"
                },
                "operators": {
                    "type": "array",
                    "items": {
                        "type": "object",
                        "properties": {
                            "label": {
                                "description": "Label is the human-readable form of the value",
                                "type": "string"
                            },
                            "value": {
                                "description": "Value is the value for the backend",
                                "allOf": [
                                    {
                                        "$ref": "#/definitions/filter.CompareOperator"
                                    }
                                ]
                            }
                        }
                    }
                },
                "values": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "query.Metadata": {
            "type": "object",
            "properties": {
                "filter": {
                    "$ref": "#/definitions/filter.Request"
                },
                "paging": {
                    "$ref": "#/definitions/paging.Response"
                },
                "sorting": {
                    "$ref": "#/definitions/sorting.Request"
                }
            }
        },
        "query.ResponseListWithMetadata-models_Notification": {
            "type": "object",
            "required": [
                "data",
                "metadata"
            ],
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Notification"
                    }
                },
                "metadata": {
                    "$ref": "#/definitions/query.Metadata"
                }
            }
        },
        "query.ResponseWithMetadata-array_query_FilterOption": {
            "type": "object",
            "required": [
                "data",
                "metadata"
            ],
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/query.FilterOption"
                    }
                },
                "metadata": {
                    "$ref": "#/definitions/query.Metadata"
                }
            }
        },
        "query.ResponseWithMetadata-models_Notification": {
            "type": "object",
            "required": [
                "data",
                "metadata"
            ],
            "properties": {
                "data": {
                    "$ref": "#/definitions/models.Notification"
                },
                "metadata": {
                    "$ref": "#/definitions/query.Metadata"
                }
            }
        },
        "query.ResultSelector": {
            "type": "object",
            "properties": {
                "filter": {
                    "$ref": "#/definitions/filter.Request"
                },
                "paging": {
                    "$ref": "#/definitions/paging.Request"
                },
                "sorting": {
                    "$ref": "#/definitions/sorting.Request"
                }
            }
        },
        "sorting.Request": {
            "type": "object",
            "properties": {
                "column": {
                    "type": "string"
                },
                "direction": {
                    "$ref": "#/definitions/sorting.SortDirection"
                }
            }
        },
        "sorting.SortDirection": {
            "type": "string",
            "enum": [
                "desc",
                "asc",
                ""
            ],
            "x-enum-varnames": [
                "DirectionDescending",
                "DirectionAscending",
                "NoDirection"
            ]
        }
    },
    "securityDefinitions": {
        "KeycloakAuth": {
            "type": "oauth2",
            "flow": "implicit",
            "authorizationUrl": "{{.KeycloakAuthUrl}}/realms/{{.KeycloakRealm}}/protocol/openid-connect/auth"
        }
    },
    "externalDocs": {
        "description": "OpenAPI",
        "url": "https://swagger.io/resources/open-api/"
    }
}`

// SwaggerInfonotificationservice holds exported Swagger Info so clients can modify it
var SwaggerInfonotificationservice = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/api/notification-service",
	Schemes:          []string{},
	Title:            "Notification Service API",
	Description:      "HTTP API of the Notification service",
	InfoInstanceName: "notificationservice",
	SwaggerTemplate:  docTemplatenotificationservice,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfonotificationservice.InstanceName(), SwaggerInfonotificationservice)
}

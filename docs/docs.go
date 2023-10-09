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
        "/healthcheck": {
            "get": {
                "description": "Returns the status and version of the application",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Health check endpoint of the Photostorage app",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/common.Health"
                        }
                    }
                }
            }
        },
        "/login": {
            "post": {
                "description": "Logs in the user, sets up the JWT authorization",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "User login endpoint",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/common.StatusMessage"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/common.StatusMessage"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/common.StatusMessage"
                        }
                    }
                }
            }
        },
        "/logout": {
            "get": {
                "description": "Logs out of the application, deletes the JWT token uased for authorization",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Logout endpoint",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/common.StatusMessage"
                        }
                    }
                }
            }
        },
        "/photos": {
            "get": {
                "description": "Returns all photo descriptors for the current user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "photos"
                ],
                "summary": "List user's photo descriptors endpoint",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/photo.Response"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/common.StatusMessage"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/common.StatusMessage"
                        }
                    }
                }
            }
        },
        "/photos/:id": {
            "get": {
                "description": "Returns the photo descriptor with the provided ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "photos"
                ],
                "summary": "Get photo endpoint",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID of the photo information to collect",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/photo.Response"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/common.StatusMessage"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/common.StatusMessage"
                        }
                    }
                }
            }
        },
        "/photos/:id/download": {
            "get": {
                "description": "Returns the RAW file for the provided ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "photos"
                ],
                "summary": "Download RAW file endpoint",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID of the RAW photo to download",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "integer"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/common.StatusMessage"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/common.StatusMessage"
                        }
                    }
                }
            }
        },
        "/profile": {
            "get": {
                "description": "Gets the current logged in user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Get user profile endpoint",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/auth.Profile"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/common.StatusMessage"
                        }
                    }
                }
            }
        },
        "/reset": {
            "post": {
                "description": "Returns the status and version of the application",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Reset password endpoint",
                "responses": {
                    "501": {
                        "description": "Not Implemented",
                        "schema": {
                            "$ref": "#/definitions/common.StatusMessage"
                        }
                    }
                }
            }
        },
        "/signup": {
            "post": {
                "description": "Registers the user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "User registration endpoint",
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/common.StatusMessage"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/common.StatusMessage"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/common.StatusMessage"
                        }
                    }
                }
            }
        },
        "/upload": {
            "post": {
                "description": "Upload a RAW file with descriptor",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "photos"
                ],
                "summary": "Photo upload endpoint",
                "parameters": [
                    {
                        "type": "file",
                        "description": "Photo to store",
                        "name": "photo",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/common.StatusMessage"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/common.StatusMessage"
                        }
                    },
                    "415": {
                        "description": "Unsupported Media Type",
                        "schema": {
                            "$ref": "#/definitions/common.StatusMessage"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/common.StatusMessage"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "auth.Profile": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "phone": {
                    "type": "string"
                }
            }
        },
        "common.Health": {
            "type": "object",
            "properties": {
                "status": {
                    "type": "string"
                },
                "version": {
                    "type": "string"
                }
            }
        },
        "common.StatusMessage": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "descriptor.Response": {
            "type": "object",
            "properties": {
                "filename": {
                    "type": "string"
                },
                "format": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "metadata": {
                    "$ref": "#/definitions/image.Response"
                },
                "tags": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "thumbnail": {
                    "type": "string"
                },
                "uploaded": {
                    "type": "string"
                }
            }
        },
        "image.Response": {
            "type": "object",
            "properties": {
                "aperture": {
                    "type": "number"
                },
                "cameraMake": {
                    "type": "string"
                },
                "cameraModel": {
                    "type": "string"
                },
                "cameraSW": {
                    "type": "string"
                },
                "colors": {
                    "type": "integer"
                },
                "dataSize": {
                    "type": "integer"
                },
                "height": {
                    "type": "integer"
                },
                "iso": {
                    "type": "integer"
                },
                "lensMake": {
                    "type": "string"
                },
                "lensModel": {
                    "type": "string"
                },
                "shutter": {
                    "type": "number"
                },
                "timestamp": {
                    "type": "string"
                },
                "width": {
                    "type": "integer"
                }
            }
        },
        "photo.Response": {
            "type": "object",
            "properties": {
                "descriptor": {
                    "$ref": "#/definitions/descriptor.Response"
                },
                "id": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}

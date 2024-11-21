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
        "/api/v1/products/{product_id}": {
            "get": {
                "description": "Retrieve a single product's details by its ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "products"
                ],
                "summary": "Get a product by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Product ID",
                        "name": "product_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/payload.ResGetProductByID"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/errorHandler.HttpError"
                        }
                    },
                    "404": {
                        "description": "Product not found",
                        "schema": {
                            "$ref": "#/definitions/errorHandler.HttpError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/errorHandler.HttpError"
                        }
                    }
                }
            },
            "put": {
                "description": "Update an existing product's details",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "products"
                ],
                "summary": "Update a product",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Product ID",
                        "name": "product_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Product information to update",
                        "name": "product",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/payload.ReqUpdateProduct"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/payload.ResUpdateProduct"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/errorHandler.HttpError"
                        }
                    },
                    "404": {
                        "description": "Product not found",
                        "schema": {
                            "$ref": "#/definitions/errorHandler.HttpError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/errorHandler.HttpError"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete a product by its ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "products"
                ],
                "summary": "Delete a product",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Product ID",
                        "name": "product_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/payload.ResDeleteProductByID"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/errorHandler.HttpError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/errorHandler.HttpError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/errorHandler.HttpError"
                        }
                    }
                }
            }
        },
        "/categories": {
            "get": {
                "description": "Retrieve all categories with optional pagination",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "categories"
                ],
                "summary": "Get all categories",
                "parameters": [
                    {
                        "minimum": 1,
                        "type": "integer",
                        "default": 1,
                        "description": "Page number (default: 1)",
                        "name": "page_number",
                        "in": "query"
                    },
                    {
                        "maximum": 100,
                        "minimum": 1,
                        "type": "integer",
                        "default": 20,
                        "description": "Page size (default: 20)",
                        "name": "page_size",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/payload.ResGetAllCategory"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/errorHandler.HttpError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/errorHandler.HttpError"
                        }
                    }
                }
            },
            "post": {
                "description": "Create a new category with the provided name",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "categories"
                ],
                "summary": "Create a new category",
                "parameters": [
                    {
                        "description": "Category to create",
                        "name": "category",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/payload.ReqCreateCategory"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/payload.ResCreateCategory"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/errorHandler.HttpError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/errorHandler.HttpError"
                        }
                    }
                }
            }
        },
        "/categories/{category_id}": {
            "get": {
                "description": "Retrieve a specific category by its ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "categories"
                ],
                "summary": "Get a category by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Category ID",
                        "name": "category_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/payload.ResGetCategoryByID"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/errorHandler.HttpError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/errorHandler.HttpError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/errorHandler.HttpError"
                        }
                    }
                }
            },
            "put": {
                "description": "Update a specific category's name by its ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "categories"
                ],
                "summary": "Edit a category by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Category ID",
                        "name": "category_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Updated category information",
                        "name": "category",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/payload.ReqUpdateCategory"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/payload.ResUpdateCategory"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/errorHandler.HttpError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/errorHandler.HttpError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/errorHandler.HttpError"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete a specific category by its ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "categories"
                ],
                "summary": "Delete a category",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Category ID",
                        "name": "category_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/payload.ResDeleteCategory"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/errorHandler.HttpError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/errorHandler.HttpError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/errorHandler.HttpError"
                        }
                    }
                }
            }
        },
        "/products": {
            "get": {
                "description": "Retrieve a list of products with optional filtering and pagination",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "products"
                ],
                "summary": "Get all products",
                "parameters": [
                    {
                        "type": "integer",
                        "default": 10,
                        "description": "Number of items per page",
                        "name": "page_size",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 1,
                        "description": "Page number",
                        "name": "page_number",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Category ID to filter products",
                        "name": "category_id",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/payload.ResGetAllProducts"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/errorHandler.HttpError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/errorHandler.HttpError"
                        }
                    }
                }
            },
            "post": {
                "description": "Create a new product with the provided details",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "products"
                ],
                "summary": "Create a new product",
                "parameters": [
                    {
                        "description": "Product details",
                        "name": "product",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/payload.ReqCreateProduct"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/payload.ResCreateProduct"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/errorHandler.HttpError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/errorHandler.HttpError"
                        }
                    }
                }
            }
        },
        "/products/generate-barcode": {
            "post": {
                "description": "Generate a PDF containing barcodes for a single product",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/pdf"
                ],
                "tags": [
                    "products"
                ],
                "summary": "Generate barcode PDF for a single product",
                "parameters": [
                    {
                        "description": "Product ID",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/payload.GenerateBarcodeRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "file"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/errorHandler.HttpError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/errorHandler.HttpError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "errorHandler.HttpError": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "info": {
                    "type": "string"
                },
                "message": {}
            }
        },
        "payload.CategoryData": {
            "type": "object",
            "properties": {
                "category_name": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "deleted_at": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "payload.GenerateBarcodeRequest": {
            "type": "object",
            "properties": {
                "product_id": {
                    "type": "string"
                }
            }
        },
        "payload.Pagination": {
            "type": "object",
            "properties": {
                "page_number": {
                    "type": "integer"
                },
                "page_size": {
                    "type": "integer"
                },
                "total_elements": {
                    "type": "integer"
                },
                "total_pages": {
                    "type": "integer"
                }
            }
        },
        "payload.ProductData": {
            "type": "object",
            "properties": {
                "base_price": {
                    "type": "number"
                },
                "category_id": {
                    "type": "string"
                },
                "category_name": {
                    "description": "From categories table",
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "deleted_at": {
                    "type": "string"
                },
                "price": {
                    "type": "number"
                },
                "product_id": {
                    "type": "string"
                },
                "product_name": {
                    "type": "string"
                },
                "stock_quantity": {
                    "type": "integer"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "payload.ReqCreateCategory": {
            "type": "object",
            "properties": {
                "category_name": {
                    "type": "string"
                }
            }
        },
        "payload.ReqCreateProduct": {
            "type": "object",
            "required": [
                "category_id",
                "product_name"
            ],
            "properties": {
                "base_price": {
                    "type": "number",
                    "minimum": 0
                },
                "category_id": {
                    "type": "string"
                },
                "price": {
                    "type": "number",
                    "minimum": 0
                },
                "product_name": {
                    "type": "string",
                    "maxLength": 40
                },
                "stock_quantity": {
                    "type": "integer",
                    "minimum": 0
                }
            }
        },
        "payload.ReqUpdateCategory": {
            "type": "object",
            "properties": {
                "category_name": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                }
            }
        },
        "payload.ReqUpdateProduct": {
            "type": "object",
            "required": [
                "category_id",
                "product_id",
                "product_name"
            ],
            "properties": {
                "base_price": {
                    "type": "number",
                    "minimum": 0
                },
                "category_id": {
                    "type": "string"
                },
                "price": {
                    "type": "number",
                    "minimum": 0
                },
                "product_id": {
                    "type": "string"
                },
                "product_name": {
                    "type": "string",
                    "maxLength": 40
                },
                "stock_quantity": {
                    "type": "integer",
                    "minimum": 0
                }
            }
        },
        "payload.ResCreateCategory": {
            "type": "object"
        },
        "payload.ResCreateProduct": {
            "type": "object"
        },
        "payload.ResDeleteCategory": {
            "type": "object"
        },
        "payload.ResDeleteProductByID": {
            "type": "object"
        },
        "payload.ResGetAllCategory": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/payload.CategoryData"
                    }
                },
                "metadata": {
                    "$ref": "#/definitions/payload.Pagination"
                }
            }
        },
        "payload.ResGetAllProducts": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/payload.ProductData"
                    }
                },
                "metadata": {
                    "$ref": "#/definitions/payload.Pagination"
                }
            }
        },
        "payload.ResGetCategoryByID": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/payload.CategoryData"
                }
            }
        },
        "payload.ResGetProductByID": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/payload.ProductData"
                }
            }
        },
        "payload.ResUpdateCategory": {
            "type": "object"
        },
        "payload.ResUpdateProduct": {
            "type": "object"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "",
	Description:      "rizki plastik point of sale api server.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}

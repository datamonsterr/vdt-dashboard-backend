-- Migration: 002_seed_data.sql
-- Description: Insert sample schema data for testing and development
-- Created: 2024-06-08

-- Insert sample blog schema
INSERT INTO schemas (
    id,
    name,
    description,
    database_name,
    status,
    version,
    schema_definition
) VALUES (
    uuid_generate_v4(),
    'Blog Schema',
    'A simple blog database schema with users, posts, and comments',
    'schema_blog_example',
    'created',
    '1.0',
    '{
        "tables": [
            {
                "id": "users_table",
                "name": "users",
                "columns": [
                    {
                        "id": "user_id",
                        "name": "id",
                        "dataType": "INT",
                        "nullable": false,
                        "primaryKey": true,
                        "autoIncrement": true
                    },
                    {
                        "id": "user_email",
                        "name": "email",
                        "dataType": "VARCHAR",
                        "length": 255,
                        "nullable": false,
                        "unique": true
                    },
                    {
                        "id": "user_name",
                        "name": "name",
                        "dataType": "VARCHAR",
                        "length": 100,
                        "nullable": false
                    },
                    {
                        "id": "user_created_at",
                        "name": "created_at",
                        "dataType": "TIMESTAMP",
                        "nullable": false
                    }
                ],
                "position": {"x": 100, "y": 100}
            },
            {
                "id": "posts_table",
                "name": "posts",
                "columns": [
                    {
                        "id": "post_id",
                        "name": "id",
                        "dataType": "INT",
                        "nullable": false,
                        "primaryKey": true,
                        "autoIncrement": true
                    },
                    {
                        "id": "post_user_id",
                        "name": "user_id",
                        "dataType": "INT",
                        "nullable": false
                    },
                    {
                        "id": "post_title",
                        "name": "title",
                        "dataType": "VARCHAR",
                        "length": 255,
                        "nullable": false
                    },
                    {
                        "id": "post_content",
                        "name": "content",
                        "dataType": "TEXT",
                        "nullable": true
                    },
                    {
                        "id": "post_published",
                        "name": "published",
                        "dataType": "BOOLEAN",
                        "nullable": false,
                        "defaultValue": false
                    }
                ],
                "position": {"x": 400, "y": 100}
            },
            {
                "id": "comments_table",
                "name": "comments",
                "columns": [
                    {
                        "id": "comment_id",
                        "name": "id",
                        "dataType": "INT",
                        "nullable": false,
                        "primaryKey": true,
                        "autoIncrement": true
                    },
                    {
                        "id": "comment_post_id",
                        "name": "post_id",
                        "dataType": "INT",
                        "nullable": false
                    },
                    {
                        "id": "comment_user_id",
                        "name": "user_id",
                        "dataType": "INT",
                        "nullable": false
                    },
                    {
                        "id": "comment_content",
                        "name": "content",
                        "dataType": "TEXT",
                        "nullable": false
                    }
                ],
                "position": {"x": 250, "y": 350}
            }
        ],
        "foreignKeys": [
            {
                "id": "fk_posts_users",
                "sourceTableId": "posts_table",
                "sourceColumnId": "post_user_id",
                "targetTableId": "users_table",
                "targetColumnId": "user_id",
                "onDelete": "CASCADE",
                "onUpdate": "CASCADE"
            },
            {
                "id": "fk_comments_posts",
                "sourceTableId": "comments_table",
                "sourceColumnId": "comment_post_id",
                "targetTableId": "posts_table",
                "targetColumnId": "post_id",
                "onDelete": "CASCADE",
                "onUpdate": "CASCADE"
            },
            {
                "id": "fk_comments_users",
                "sourceTableId": "comments_table",
                "sourceColumnId": "comment_user_id",
                "targetTableId": "users_table",
                "targetColumnId": "user_id",
                "onDelete": "CASCADE",
                "onUpdate": "CASCADE"
            }
        ],
        "version": "1.0",
        "exportedAt": "2024-06-08T12:00:00.000Z"
    }'::jsonb
) ON CONFLICT (name) DO NOTHING;

-- Insert simple e-commerce schema
INSERT INTO schemas (
    id,
    name,
    description,
    database_name,
    status,
    version,
    schema_definition
) VALUES (
    uuid_generate_v4(),
    'E-commerce Schema',
    'Basic e-commerce database with products and orders',
    'schema_ecommerce_example',
    'created',
    '1.0',
    '{
        "tables": [
            {
                "id": "products_table",
                "name": "products",
                "columns": [
                    {
                        "id": "product_id",
                        "name": "id",
                        "dataType": "INT",
                        "nullable": false,
                        "primaryKey": true,
                        "autoIncrement": true
                    },
                    {
                        "id": "product_name",
                        "name": "name",
                        "dataType": "VARCHAR",
                        "length": 255,
                        "nullable": false
                    },
                    {
                        "id": "product_price",
                        "name": "price",
                        "dataType": "DECIMAL",
                        "precision": 10,
                        "scale": 2,
                        "nullable": false
                    },
                    {
                        "id": "product_stock",
                        "name": "stock",
                        "dataType": "INT",
                        "nullable": false,
                        "defaultValue": 0
                    }
                ],
                "position": {"x": 100, "y": 100}
            },
            {
                "id": "orders_table",
                "name": "orders",
                "columns": [
                    {
                        "id": "order_id",
                        "name": "id",
                        "dataType": "INT",
                        "nullable": false,
                        "primaryKey": true,
                        "autoIncrement": true
                    },
                    {
                        "id": "order_total",
                        "name": "total",
                        "dataType": "DECIMAL",
                        "precision": 10,
                        "scale": 2,
                        "nullable": false
                    },
                    {
                        "id": "order_status",
                        "name": "status",
                        "dataType": "VARCHAR",
                        "length": 50,
                        "nullable": false,
                        "defaultValue": "pending"
                    },
                    {
                        "id": "order_created_at",
                        "name": "created_at",
                        "dataType": "TIMESTAMP",
                        "nullable": false
                    }
                ],
                "position": {"x": 400, "y": 100}
            }
        ],
        "foreignKeys": [],
        "version": "1.0",
        "exportedAt": "2024-06-08T12:00:00.000Z"
    }'::jsonb
) ON CONFLICT (name) DO NOTHING; 
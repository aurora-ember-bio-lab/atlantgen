schema "public" {}

# Target data model that every connector (WordPress, Shopify, WooCommerce,
# Magento) normalizes into. Mirrors the shape MedusaJS expects downstream.

table "products" {
  schema = schema.public

  column "id" {
    type = uuid
  }
  column "title" {
    type = text
  }
  column "source_type" {
    type = text
  }
  column "source_id" {
    type = text
  }
  column "price_cents" {
    type = int
    null = true
  }
  column "created_at" {
    type    = timestamptz
    default = sql("now()")
  }

  primary_key {
    columns = [column.id]
  }

  index "products_source_idx" {
    columns = [column.source_type, column.source_id]
    unique  = true
  }
}

table "orders" {
  schema = schema.public

  column "id" {
    type = uuid
  }
  column "source_type" {
    type = text
  }
  column "source_id" {
    type = text
  }
  column "total_cents" {
    type = int
  }
  column "created_at" {
    type    = timestamptz
    default = sql("now()")
  }

  primary_key {
    columns = [column.id]
  }

  index "orders_source_idx" {
    columns = [column.source_type, column.source_id]
    unique  = true
  }
}

table "customers" {
  schema = schema.public

  column "id" {
    type = uuid
  }
  column "email" {
    type = text
  }
  column "source_type" {
    type = text
  }
  column "source_id" {
    type = text
  }

  primary_key {
    columns = [column.id]
  }

  index "customers_email_idx" {
    columns = [column.email]
    unique  = true
  }
}

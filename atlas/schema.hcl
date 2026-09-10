schema "public" {}

table "migrations" {
  schema = schema.public

  column "id" {
    type = varchar(64)
  }
  column "name" {
    type = varchar(255)
  }
  column "source_type" {
    type = varchar(64)
  }
  column "status" {
    type    = varchar(32)
    default = sql("'queued'")
  }
  column "worker_job_id" {
    type    = varchar(64)
    null    = true
  }
  column "target_db_url" {
    type = text
  }
  column "created_at" {
    type    = timestamptz
    default = sql("now()")
  }
  column "updated_at" {
    type    = timestamptz
    default = sql("now()")
  }

  primary_key {
    columns = [column.id]
  }

  index "migrations_source_type_idx" {
    columns = [column.source_type]
  }

  index "migrations_status_idx" {
    columns = [column.status]
  }
}

table "accounts" {
  schema = schema.public

  column "id" {
    type = uuid
  }
  column "login" {
    type = varchar(255)
  }
  column "source_type" {
    type    = varchar(64)
    default = sql("'github_marketplace'")
  }
  column "stripe_customer_id" {
    type = varchar(255)
    null = true
  }
  column "github_installation_id" {
    type = int8
    null = true
  }
  column "status" {
    type    = varchar(32)
    default = sql("'active'")
  }
  column "plan_name" {
    type = varchar(255)
    null = true
  }
  column "email" {
    type = varchar(255)
    null = true
  }
  column "created_at" {
    type    = timestamptz
    default = sql("now()")
  }
  column "updated_at" {
    type    = timestamptz
    default = sql("now()")
  }

  primary_key {
    columns = [column.id]
  }

  index "accounts_login_idx" {
    columns = [column.login]
    unique  = true
  }
}

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
  column "created_at" {
    type    = timestamptz
    default = sql("now()")
  }

  primary_key {
    columns = [column.id]
  }

  index "customers_email_idx" {
    columns = [column.email]
    unique  = true
  }
}

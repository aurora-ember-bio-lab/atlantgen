schema "public" {
  charset = "utf8"
}

table "migrations" {
  column "id" {
    type = int
    auto_increment = true
  }
  column "source" {
    type = varchar(64)
  }
  column "status" {
    type = varchar(32)
  }
  column "created_at" {
    type = timestamp
  }
  primary_key {
    columns = [column.id]
  }
}

table "accounts" {
  column "id" {
    type = uuid
  }
  column "login" {
    type = varchar(255)
  }
  column "source_type" {
    type = varchar(64)
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
    type = varchar(32)
    default = sql("'active'")
  }
  column "plan_name" {
    type = varchar(255)
    null = true
  }
  column "created_at" {
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

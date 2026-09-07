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

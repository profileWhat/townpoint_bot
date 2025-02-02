-- create "users" table
CREATE TABLE "users" ("id" uuid NOT NULL, "create_time" timestamptz NOT NULL, "update_time" timestamptz NOT NULL, "tg_id" character varying NOT NULL, "role" character varying NOT NULL, PRIMARY KEY ("id"));
-- create index "users_tg_id_key" to table: "users"
CREATE UNIQUE INDEX "users_tg_id_key" ON "users" ("tg_id");

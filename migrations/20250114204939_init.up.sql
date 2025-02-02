-- create "regions" table
CREATE TABLE "regions" ("id" uuid NOT NULL, "create_time" timestamptz NOT NULL, "update_time" timestamptz NOT NULL, "name" character varying NOT NULL, "description" character varying NULL, "status" character varying NOT NULL, PRIMARY KEY ("id"));
-- create index "region_id" to table: "regions"
CREATE UNIQUE INDEX "region_id" ON "regions" ("id");
-- create "towns" table
CREATE TABLE "towns" ("id" uuid NOT NULL, "create_time" timestamptz NOT NULL, "update_time" timestamptz NOT NULL, "name" character varying NOT NULL, "description" character varying NULL, "status" character varying NOT NULL, "town_region" uuid NULL, PRIMARY KEY ("id"), CONSTRAINT "towns_regions_region" FOREIGN KEY ("town_region") REFERENCES "regions" ("id") ON UPDATE NO ACTION ON DELETE SET NULL);
-- create index "town_id" to table: "towns"
CREATE UNIQUE INDEX "town_id" ON "towns" ("id");
-- create "points" table
CREATE TABLE "points" ("id" uuid NOT NULL, "create_time" timestamptz NOT NULL, "update_time" timestamptz NOT NULL, "pictures" jsonb NOT NULL, "videos" jsonb NOT NULL, "address" character varying NOT NULL, "description" character varying NULL, "status" character varying NOT NULL, "point_town" uuid NULL, PRIMARY KEY ("id"), CONSTRAINT "points_towns_town" FOREIGN KEY ("point_town") REFERENCES "towns" ("id") ON UPDATE NO ACTION ON DELETE SET NULL);
-- create index "point_id" to table: "points"
CREATE UNIQUE INDEX "point_id" ON "points" ("id");

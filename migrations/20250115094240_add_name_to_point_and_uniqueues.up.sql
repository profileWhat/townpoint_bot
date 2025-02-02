-- modify "points" table
ALTER TABLE "points" ADD COLUMN "name" character varying NOT NULL;
-- create index "points_address_key" to table: "points"
CREATE UNIQUE INDEX "points_address_key" ON "points" ("address");
-- create index "points_name_key" to table: "points"
CREATE UNIQUE INDEX "points_name_key" ON "points" ("name");
-- create index "regions_name_key" to table: "regions"
CREATE UNIQUE INDEX "regions_name_key" ON "regions" ("name");
-- create index "towns_name_key" to table: "towns"
CREATE UNIQUE INDEX "towns_name_key" ON "towns" ("name");

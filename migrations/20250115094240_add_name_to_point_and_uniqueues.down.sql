-- reverse: create index "towns_name_key" to table: "towns"
DROP INDEX "towns_name_key";
-- reverse: create index "regions_name_key" to table: "regions"
DROP INDEX "regions_name_key";
-- reverse: create index "points_name_key" to table: "points"
DROP INDEX "points_name_key";
-- reverse: create index "points_address_key" to table: "points"
DROP INDEX "points_address_key";
-- reverse: modify "points" table
ALTER TABLE "points" DROP COLUMN "name";

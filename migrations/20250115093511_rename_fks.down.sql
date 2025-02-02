-- reverse: modify "points" table
ALTER TABLE "points" DROP CONSTRAINT "points_towns_town", DROP COLUMN "town_id", ADD COLUMN "point_town" uuid NULL, ADD CONSTRAINT "points_towns_town" FOREIGN KEY ("point_town") REFERENCES "towns" ("id") ON UPDATE NO ACTION ON DELETE SET NULL;
-- reverse: modify "towns" table
ALTER TABLE "towns" DROP CONSTRAINT "towns_regions_region", DROP COLUMN "region_id", ADD COLUMN "town_region" uuid NULL, ADD CONSTRAINT "towns_regions_region" FOREIGN KEY ("town_region") REFERENCES "regions" ("id") ON UPDATE NO ACTION ON DELETE SET NULL;

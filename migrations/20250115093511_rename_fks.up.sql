-- modify "towns" table
ALTER TABLE "towns" DROP CONSTRAINT "towns_regions_region", DROP COLUMN "town_region", ADD COLUMN "region_id" uuid NOT NULL, ADD CONSTRAINT "towns_regions_region" FOREIGN KEY ("region_id") REFERENCES "regions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- modify "points" table
ALTER TABLE "points" DROP CONSTRAINT "points_towns_town", DROP COLUMN "point_town", ADD COLUMN "town_id" uuid NOT NULL, ADD CONSTRAINT "points_towns_town" FOREIGN KEY ("town_id") REFERENCES "towns" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;

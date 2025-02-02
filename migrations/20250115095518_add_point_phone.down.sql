-- reverse: modify "points" table
ALTER TABLE "points" DROP COLUMN "phone", ALTER COLUMN "videos" SET NOT NULL, ALTER COLUMN "pictures" SET NOT NULL;

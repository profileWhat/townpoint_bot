-- modify "points" table
ALTER TABLE "points" ALTER COLUMN "pictures" DROP NOT NULL, ALTER COLUMN "videos" DROP NOT NULL, ADD COLUMN "phone" character varying NOT NULL;

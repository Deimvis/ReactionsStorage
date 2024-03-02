DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM "pg_type" WHERE typname = 'reaction_type') THEN
		CREATE TYPE "reaction_type" AS ENUM (
			'unicode',
			'custom'
		);
    END IF;
END$$;

CREATE TABLE IF NOT EXISTS "reaction" (
	-- "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	"id" TEXT PRIMARY KEY,
	"short_name" VARCHAR(32),
	"type" reaction_type NOT NULL,
	"code" CHAR(1),
	"url" TEXT,

	CHECK (
		(type = 'unicode' AND "code" IS NOT NULL) OR
		(type = 'custom' AND "url" IS NOT NULL)
	)
);

CREATE TABLE IF NOT EXISTS "reaction_set" (
	-- "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	"id" TEXT PRIMARY KEY,
	"reaction_ids" TEXT[] NOT NULL
);

CREATE TABLE IF NOT EXISTS "namespace" (
	-- "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	"id" TEXT PRIMARY KEY,
	"reaction_set_id" TEXT NOT NULL,
	"max_uniq_reactions" INT NOT NULL,
	"mutually_exclusive_reactions" TEXT[][] NOT NULL
);
COMMENT ON COLUMN namespace.mutually_exclusive_reactions is
'Ids of mutually exlcusive reactions';

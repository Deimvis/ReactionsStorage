CREATE TABLE IF NOT EXISTS "user_reaction" (
	"namespace_id" TEXT NOT NULL,
	"entity_id" TEXT NOT NULL,
	"reaction_id" TEXT NOT NULL,
	"user_id" TEXT NOT NULL,
	"timestamp" BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())::BIGINT,
	PRIMARY KEY ("namespace_id", "entity_id", "reaction_id", "user_id")
);

CREATE MATERIALIZED VIEW IF NOT EXISTS "entity_reactions"
AS
	WITH tmp AS (
		SELECT
			"namespace_id",
			"entity_id",
			"reaction_id",
			COUNT(*) AS "reaction_count"
		FROM "user_reaction"
		GROUP BY namespace_id, entity_id, reaction_id
	)
	SELECT
		"namespace_id" || '__' || "entity_id" AS "namespace_id__entity_id",
		json_object_agg(
			"reaction_id", "reaction_count"
		) AS "reactions_count"
	FROM tmp
	GROUP BY namespace_id, entity_id
WITH NO DATA;

CREATE INDEX IF NOT EXISTS "entity_reactions_pkey" ON "entity_reactions" USING HASH (
	"namespace_id__entity_id"
);

REFRESH MATERIALIZED VIEW "entity_reactions";

-- SELECT cron.schedule(
--   'opinion_activity',
--   '* * * * *',
--   $CRON$ REFRESH MATERIALIZED VIEW example.opinion_activity; $CRON$
-- );
-- SELECT cron.schedule(
--   'opinion_activity',
--   '* * * * *',
--   $CRON$ SELECT pg_sleep(30); REFRESH MATERIALIZED VIEW example.opinion_activity; $CRON$
-- );
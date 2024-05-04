CREATE TABLE IF NOT EXISTS "user_reactions" (
	"namespace_id" TEXT NOT NULL,
	"entity_id" TEXT NOT NULL,
	"user_id" TEXT NOT NULL,
	"reaction_ids" TEXT[] NOT NULL,
	"create_ts" BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())::BIGINT,
	"last_update_ts" BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())::BIGINT,
	PRIMARY KEY ("namespace_id", "entity_id", "user_id")
) WITH (fillfactor = 97);

CREATE TABLE IF NOT EXISTS "reactions_count" (
    "namespace_id" TEXT NOT NULL,
	"entity_id" TEXT NOT NULL,
    "reactions_count" JSONB NOT NULL,
	PRIMARY KEY ("namespace_id", "entity_id")
) WITH (fillfactor = 93);

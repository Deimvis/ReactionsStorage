INSERT INTO "reaction" ("id", "type", "code") VALUES
('reaction', 'unicode', '😃'),
('reaction1', 'unicode', '1'),
('reaction2', 'unicode', '2'),
('reaction3', 'unicode', '3') ON CONFLICT ("id") DO NOTHING;

INSERT INTO "reaction_set" ("id", "reaction_ids") VALUES
('reaction_set', array['reaction']) ON CONFLICT ("id") DO NOTHING;

INSERT INTO "namespace" ("id", reaction_set_id, max_uniq_reactions, mutually_exclusive_reactions) VALUES
('namespace', 'reaction_set', 10, array[array['reaction1', 'reaction2'], array['reaction2', 'reaction3']]) ON CONFLICT ("id") DO NOTHING;

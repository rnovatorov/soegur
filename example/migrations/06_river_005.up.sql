-- River migration 005 [up]
--
-- Rebuild the migration table so it's based on `(line, version)`.
--
DO $body$
BEGIN
    -- Tolerate users who may be using their own migration system rather than
    -- River's. If they are, they will have skipped version 001 containing
    -- `CREATE TABLE river_migration`, so this table won't exist.
    IF (
        SELECT
            to_regclass ('river_migration') IS NOT NULL) THEN
        ALTER TABLE river_migration RENAME TO river_migration_old;
        CREATE TABLE river_migration (
            line TEXT NOT NULL,
            version BIGINT NOT NULL,
            created_at TIMESTAMPTZ NOT NULL DEFAULT NOW( ),
            CONSTRAINT line_length CHECK (char_length(LINE ) > 0 AND char_length(LINE ) < 128 ),
            CONSTRAINT version_gte_1 CHECK (version >= 1 ),
            PRIMARY KEY (LINE, version )
        );
    INSERT INTO river_migration (created_at, line, version)
    SELECT
        created_at,
        'main',
        version
    FROM
        river_migration_old;
    DROP TABLE river_migration_old;
END IF;
END;
$body$
LANGUAGE 'plpgsql';

--
-- Add `river_job.unique_key` and bring up an index on it.
--
-- These statements use `IF NOT EXISTS` to allow users with a `river_job` table
-- of non-trivial size to build the index `CONCURRENTLY` out of band of this
-- migration, then follow by completing the migration.
ALTER TABLE river_job
    ADD COLUMN IF NOT EXISTS unique_key bytea;

CREATE UNIQUE INDEX IF NOT EXISTS river_job_kind_unique_key_idx ON river_job (kind, unique_key)
WHERE
    unique_key IS NOT NULL;

--
-- Create `river_client` and derivative.
--
-- This feature hasn't quite yet been implemented, but we're taking advantage of
-- the migration to add the schema early so that we can add it later without an
-- additional migration.
--
CREATE UNLOGGED TABLE river_client (
    id TEXT PRIMARY KEY NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    metadata JSONB NOT NULL DEFAULT '{}',
    paused_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL,
    CONSTRAINT name_length CHECK (char_length(id) > 0 AND char_length(id) < 128)
);

-- Differs from `river_queue` in that it tracks the queue state for a particular
-- active client.
CREATE UNLOGGED TABLE river_client_queue (
    river_client_id TEXT NOT NULL REFERENCES river_client (id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    max_workers BIGINT NOT NULL DEFAULT 0,
    metadata JSONB NOT NULL DEFAULT '{}',
    num_jobs_completed BIGINT NOT NULL DEFAULT 0,
    num_jobs_running BIGINT NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (river_client_id, NAME),
    CONSTRAINT name_length CHECK (char_length(NAME) > 0 AND char_length(NAME) < 128),
    CONSTRAINT num_jobs_completed_zero_or_positive CHECK (num_jobs_completed >= 0),
    CONSTRAINT num_jobs_running_zero_or_positive CHECK (num_jobs_running >= 0)
);

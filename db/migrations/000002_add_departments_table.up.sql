BEGIN;

CREATE TABLE departments (
    id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name        text NOT NULL,
    created_at  timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMIT;
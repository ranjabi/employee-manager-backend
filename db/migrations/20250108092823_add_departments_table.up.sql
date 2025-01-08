BEGIN;

CREATE TABLE departments (
    id      uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name    text NOT NULL
);

COMMIT;
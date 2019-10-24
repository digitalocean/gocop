-- DOWN
DROP TABLE IF EXISTS test;
DROP TABLE IF EXISTS run CASCADE;

-- UP
CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;

CREATE TABLE run (
  created   TIMESTAMPTZ PRIMARY KEY,
  repo      TEXT,
  branch    TEXT,
  sha       TEXT,
  build_id  BIGINT,
  cmd       TEXT,
  benchmark BOOL,
  race      BOOL,
  short     BOOL,
  tags      TEXT,
  hash      TEXT,
  duration  INTEGER
);

SELECT create_hypertable('run', 'created');

CREATE OR REPLACE FUNCTION set_run_hash()
RETURNS trigger AS
$$
BEGIN
  NEW.hash := MD5(concat(NEW.benchmark, NEW.race, NEW.short, NEW.tags));
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS run_hash_insert on run;
CREATE TRIGGER run_hash_insert
BEFORE INSERT ON "run"
FOR EACH ROW
EXECUTE PROCEDURE set_run_hash();

DROP TRIGGER IF EXISTS run_hash on run;
CREATE TRIGGER run_hash
BEFORE UPDATE ON run
FOR EACH ROW
WHEN ((OLD.benchmark, OLD.race, OLD.short, OLD.tags)
  IS DISTINCT FROM
  (NEW.benchmark, NEW.race, NEW.short, NEW.tags))
EXECUTE PROCEDURE set_run_hash();

DROP TABLE IF EXISTS test;
CREATE TABLE test (
  created   TIMESTAMPTZ,
  package   TEXT,
  result    TEXT CHECK (result in ('pass', 'fail', 'flaky', 'skip')),
  duration  INTEGER
);

SELECT create_hypertable('test', 'created');

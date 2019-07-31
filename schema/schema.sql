CREATE TABLE run (
  created     TIMESTAMPTZ,
  build_id    BIGINT,
  duration    INTEGER,
  run_cmd     TEXT,
  repo        TEXT,
  branch      TEXT,
  sha         TEXT,
  race        BOOL,
  short       BOOL,
  tags        TEXT,
  hash        TEXT
);

SELECT create_hypertable('run', 'created');

CREATE OR REPLACE FUNCTION set_run_hash()
RETURNS trigger AS 
$$
BEGIN
  NEW.hash := MD5(concat(NEW.race, NEW.short, NEW.tags));
  RETURN NEW;
END; 
$$ LANGUAGE plpgsql;

CREATE TRIGGER run_hash_insert
BEFORE INSERT ON "run"
FOR EACH ROW 
EXECUTE PROCEDURE set_run_hash();

CREATE TRIGGER run_hash
BEFORE UPDATE ON run
FOR EACH ROW
WHEN ((OLD.race, OLD.short, OLD.tags)
  IS DISTINCT FROM
  (NEW.race, NEW.short, NEW.tags))
EXECUTE PROCEDURE set_run_hash();


CREATE TABLE test (
  created   TIMESTAMPTZ,
  id        BIGSERIAL,
  run_id    BIGINT,
  result    TEXT CHECK (result in ('pass', 'fail', 'flaky', 'skip')),
  name      TEXT,
  duration  INTEGER
);

SELECT create_hypertable('test', 'created');


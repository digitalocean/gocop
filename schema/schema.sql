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
  tags        TEXT
);

SELECT create_hypertable('run', 'created');

CREATE TABLE test (
  created   TIMESTAMPTZ,
  id        BIGSERIAL,
  run_id    BIGINT,
  result    TEXT CHECK (result in ('pass', 'fail', 'flaky', 'skip')),
  name      TEXT,
  duration  INTEGER
);

SELECT create_hypertable('test', 'created');
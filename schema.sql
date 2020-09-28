CREATE TABLE IF NOT EXISTS users (
  user_id             VARCHAR(36) NOT NULL PRIMARY KEY,
  full_name           VARCHAR(1024) NOT NULL,
  identity_number     VARCHAR(11) NOT NULL UNIQUE,
  avatar_url          VARCHAR(1024) NOT NULL,
  access_token        VARCHAR(512) NOT NULL DEFAULT '',
  created_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS user_bots (
  user_id      VARCHAR(36) NOT NULL,
  client_id    VARCHAR(36) NOT NULL,
  session_id   VARCHAR(36) NOT NULL,
  private_key  VARCHAR NOT NULL,
  hash         VARCHAR NOT NULL UNIQUE,
  PRIMARY KEY(user_id, client_id)
);
CREATE INDEX user_bots_user_id_idx on user_bots(user_id)
CREATE INDEX user_bots_client_id_idx on user_bots(client_id)

CREATE TABLE IF NOT EXISTS bots (
  client_id           VARCHAR(36) NOT NULL PRIMARY KEY,
  session_id          VARCHAR(36) NOT NULL,
  private_key         VARCHAR NOT NULL,
  full_name           VARCHAR(1024) NOT NULL,
  identity_number     VARCHAR(11) NOT NULL UNIQUE,
  avatar_url          VARCHAR(1024) NOT NULL,
  hash                VARCHAR NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS daily_data (
  client_id  VARCHAR(36) NOT NULL,
  date       DATE NOT NULL,
  users      INTEGER NOT NULL DEFAULT 0,
  messages   INTEGER NOT NULL DEFAULT 0,
  PRIMARY KEY(client_id, date)
);


CREATE TABLE IF NOT EXISTS daily_data (
  client_id  VARCHAR(36) NOT NULL,
  date       DATE NOT NULL,
  users      INTEGER NOT NULL DEFAULT 0,
  messages   INTEGER NOT NULL DEFAULT 0,
  PRIMARY KEY(client_id, date)
);


CREATE TABLE IF NOT EXISTS forward_messages (
  client_id              VARCHAR(36) NOT NULL,
  message_id             VARCHAR(36) NOT NULL,
  admin_id               VARCHAR(36) NOT NULL,
  recipient_id           VARCHAR(36) NOT NULL,
  origin_message_id      VARCHAR(36) NOT NULL,
  conversation_id        VARCHAR(36) NOT NULL,
  admin_message_id       VARCHAR(36),
  created_at             TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  PRIMARY KEY(client_id, message_id)
);


CREATE TABLE IF NOT EXISTS messages (
  message_id            VARCHAR(36) PRIMARY KEY CHECK (message_id ~* '^[0-9a-f-]{36,36}$'),
  user_id               VARCHAR(36) NOT NULL CHECK (user_id ~* '^[0-9a-f-]{36,36}$'),
  category              VARCHAR(512) NOT NULL,
  quote_message_id      VARCHAR(36) NOT NULL DEFAULT '',
  data                  TEXT NOT NULL,
  created_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  state                 VARCHAR(128) NOT NULL,
  last_distribute_at    TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX IF NOT EXISTS messages_state_updatedx ON messages(state, updated_at);

CREATE TABLE IF NOT EXISTS distributed_messages (
  message_id            VARCHAR(36) PRIMARY KEY CHECK (message_id ~* '^[0-9a-f-]{36,36}$'),
  conversation_id       VARCHAR(36) NOT NULL CHECK (conversation_id ~* '^[0-9a-f-]{36,36}$'),
  recipient_id          VARCHAR(36) NOT NULL CHECK (recipient_id ~* '^[0-9a-f-]{36,36}$'),
  user_id               VARCHAR(36) NOT NULL CHECK (user_id ~* '^[0-9a-f-]{36,36}$'),
  parent_id             VARCHAR(36) NOT NULL CHECK (parent_id ~* '^[0-9a-f-]{36,36}$'),
  quote_message_id      VARCHAR(36) NOT NULL DEFAULT '',
  shard                 VARCHAR(36) NOT NULL,
  category              VARCHAR(512) NOT NULL,
  data                  TEXT NOT NULL,
  status                VARCHAR(512) NOT NULL,
  created_at            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS message_shard_statusx ON distributed_messages(shard, status, created_at);
-- Migration: initial_schema
-- Generated: 2026-01-07T13:36:31+05:45

-- Table: organizations
CREATE TABLE IF NOT EXISTS organizations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT(100) NOT NULL,
    slug TEXT NOT NULL,
    created_at DATETIME,
    updated_at DATETIME,
    is_active INTEGER DEFAULT true,
    metadata TEXT
);

-- Table: organizations
CREATE TABLE IF NOT EXISTS chat_channels (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    organization_id INTEGER NOT NULL,
    platform TEXT NOT NULL,
    name TEXT(100) NOT NULL,
    account_identifier TEXT NOT NULL,
    status TEXT DEFAULT active,
    webhook_secret TEXT,
    access_token TEXT,
    config TEXT,
    created_at DATETIME,
    updated_at DATETIME,
    last_message_at DATETIME,
    is_active INTEGER DEFAULT true
);

-- Table: organizations
CREATE TABLE IF NOT EXISTS external_users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    channel_id INTEGER NOT NULL,
    platform_user_id TEXT NOT NULL,
    platform_username TEXT,
    display_name TEXT,
    phone_number TEXT,
    email TEXT,
    avatar_url TEXT,
    metadata TEXT,
    first_seen_at DATETIME,
    last_seen_at DATETIME,
    is_blocked INTEGER DEFAULT false
);

-- Table: organizations
CREATE TABLE IF NOT EXISTS conversations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    channel_id INTEGER NOT NULL,
    external_user_id INTEGER NOT NULL,
    assigned_to_external_id TEXT,
    status TEXT DEFAULT open,
    priority TEXT DEFAULT normal,
    subject TEXT,
    first_message_at DATETIME,
    last_message_at DATETIME,
    resolved_at DATETIME,
    created_at DATETIME,
    updated_at DATETIME,
    metadata TEXT
);

-- Table: organizations
CREATE TABLE IF NOT EXISTS messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    conversation_id INTEGER NOT NULL,
    platform_message_id TEXT,
    sender_type TEXT NOT NULL,
    sender_id INTEGER,
    content TEXT NOT NULL,
    message_type TEXT DEFAULT text,
    media_url TEXT,
    direction TEXT NOT NULL,
    status TEXT DEFAULT received,
    created_at DATETIME,
    delivered_at DATETIME,
    read_at DATETIME,
    metadata TEXT
);

-- Table: organizations
CREATE TABLE IF NOT EXISTS webhook_events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    channel_id INTEGER NOT NULL,
    event_type TEXT NOT NULL,
    payload TEXT NOT NULL,
    processed INTEGER DEFAULT false,
    created_at DATETIME,
    processed_at DATETIME,
    error TEXT
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_organizations_slug ON organizations(slug);
CREATE INDEX IF NOT EXISTS idx_chat_channels_org ON chat_channels(organization_id);
CREATE INDEX IF NOT EXISTS idx_chat_channels_platform ON chat_channels(platform);
CREATE INDEX IF NOT EXISTS idx_external_users_channel ON external_users(channel_id);
CREATE INDEX IF NOT EXISTS idx_conversations_channel ON conversations(channel_id);
CREATE INDEX IF NOT EXISTS idx_messages_conversation ON messages(conversation_id);
CREATE INDEX IF NOT EXISTS idx_webhook_events_channel ON webhook_events(channel_id);

-- 001_initial_schema.sql
-- Multi-channel chat aggregator database schema
-- Supports WhatsApp Business, Telegram, and other external chat platforms

-- Organizations table (businesses using the service)
CREATE TABLE IF NOT EXISTS organizations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    slug TEXT UNIQUE NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    is_active INTEGER DEFAULT 1,
    metadata TEXT, -- JSON: settings, branding, etc.
    CONSTRAINT slug_format CHECK(slug GLOB '[a-z0-9-]*')
);

-- Chat channels (connected external accounts: WhatsApp, Telegram, etc.)
CREATE TABLE IF NOT EXISTS chat_channels (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    organization_id INTEGER NOT NULL,
    platform TEXT NOT NULL CHECK(platform IN ('whatsapp', 'telegram', 'instagram', 'facebook', 'sms', 'email', 'web')),
    name TEXT NOT NULL, -- Friendly name: "Support WhatsApp", "Sales Telegram"
    account_identifier TEXT NOT NULL, -- Phone number, bot token, page ID, etc.
    status TEXT DEFAULT 'active' CHECK(status IN ('active', 'inactive', 'error', 'pending')),
    webhook_secret TEXT, -- For validating incoming webhooks
    access_token TEXT, -- API access token (encrypted)
    config TEXT, -- JSON: platform-specific settings
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_message_at DATETIME,
    is_active INTEGER DEFAULT 1,
    UNIQUE(organization_id, platform, account_identifier),
    FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE
);

-- External users (customers/people messaging from external platforms)
CREATE TABLE IF NOT EXISTS external_users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    channel_id INTEGER NOT NULL,
    platform_user_id TEXT NOT NULL, -- WhatsApp phone, Telegram user ID, etc.
    platform_username TEXT, -- @username for Telegram, etc.
    display_name TEXT,
    phone_number TEXT,
    email TEXT,
    avatar_url TEXT,
    metadata TEXT, -- JSON: platform-specific user data
    first_seen_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_seen_at DATETIME,
    is_blocked INTEGER DEFAULT 0,
    UNIQUE(channel_id, platform_user_id),
    FOREIGN KEY (channel_id) REFERENCES chat_channels(id) ON DELETE CASCADE
);

-- Conversations (chat sessions with external users)
CREATE TABLE IF NOT EXISTS conversations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    channel_id INTEGER NOT NULL,
    external_user_id INTEGER NOT NULL,
    assigned_to_external_id TEXT, -- CRM user ID (managed by NestJS/external system)
    status TEXT DEFAULT 'open' CHECK(status IN ('open', 'pending', 'resolved', 'closed')),
    priority TEXT DEFAULT 'normal' CHECK(priority IN ('low', 'normal', 'high', 'urgent')),
    subject TEXT, -- Optional conversation title
    first_message_at DATETIME,
    last_message_at DATETIME,
    resolved_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    metadata TEXT, -- JSON: tags, custom fields, CRM data, etc.
    FOREIGN KEY (channel_id) REFERENCES chat_channels(id) ON DELETE CASCADE,
    FOREIGN KEY (external_user_id) REFERENCES external_users(id) ON DELETE CASCADE
);

-- Messages (actual message content)
CREATE TABLE IF NOT EXISTS messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    conversation_id INTEGER NOT NULL,
    platform_message_id TEXT, -- External platform's message ID
    sender_type TEXT NOT NULL CHECK(sender_type IN ('external', 'internal', 'system')),
    sender_id INTEGER, -- external_user_id or internal_user_id based on sender_type
    content TEXT NOT NULL,
    message_type TEXT DEFAULT 'text' CHECK(message_type IN ('text', 'image', 'video', 'audio', 'file', 'location', 'contact', 'sticker', 'system')),
    media_url TEXT, -- For image/video/audio/file messages
    direction TEXT NOT NULL CHECK(direction IN ('inbound', 'outbound')),
    status TEXT DEFAULT 'received' CHECK(status IN ('received', 'sent', 'delivered', 'read', 'failed')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    delivered_at DATETIME,
    read_at DATETIME,
    metadata TEXT, -- JSON: reply_to, forwarded, reactions, etc.
    FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE,
    CONSTRAINT content_not_empty CHECK(length(trim(content)) > 0 OR message_type != 'text')
);

-- Webhook events log (for debugging and replay)
CREATE TABLE IF NOT EXISTS webhook_events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    channel_id INTEGER NOT NULL,
    event_type TEXT NOT NULL, -- message, status_update, etc.
    payload TEXT NOT NULL, -- Full JSON payload
    processed INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    processed_at DATETIME,
    error TEXT,
    FOREIGN KEY (channel_id) REFERENCES chat_channels(id) ON DELETE CASCADE
);

-- Indexes for performance optimization
CREATE INDEX IF NOT EXISTS idx_organizations_slug ON organizations(slug);

CREATE INDEX IF NOT EXISTS idx_chat_channels_org ON chat_channels(organization_id);
CREATE INDEX IF NOT EXISTS idx_chat_channels_platform ON chat_channels(platform);
CREATE INDEX IF NOT EXISTS idx_chat_channels_status ON chat_channels(status);

CREATE INDEX IF NOT EXISTS idx_external_users_channel ON external_users(channel_id);
CREATE INDEX IF NOT EXISTS idx_external_users_platform_id ON external_users(channel_id, platform_user_id);
CREATE INDEX IF NOT EXISTS idx_external_users_phone ON external_users(phone_number);

CREATE INDEX IF NOT EXISTS idx_conversations_channel ON conversations(channel_id);
CREATE INDEX IF NOT EXISTS idx_conversations_external_user ON conversations(external_user_id);
CREATE INDEX IF NOT EXISTS idx_conversations_assigned ON conversations(assigned_to_external_id);
CREATE INDEX IF NOT EXISTS idx_conversations_status ON conversations(status);
CREATE INDEX IF NOT EXISTS idx_conversations_updated ON conversations(updated_at DESC);

CREATE INDEX IF NOT EXISTS idx_messages_conversation ON messages(conversation_id);
CREATE INDEX IF NOT EXISTS idx_messages_created ON messages(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_messages_platform_id ON messages(platform_message_id);
CREATE INDEX IF NOT EXISTS idx_messages_conversation_created ON messages(conversation_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_webhook_events_channel ON webhook_events(channel_id);
CREATE INDEX IF NOT EXISTS idx_webhook_events_processed ON webhook_events(processed, created_at);

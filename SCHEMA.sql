-- ZemenLink Core Schema (PostgreSQL)

-- 1. Global Management Database (centralized)
CREATE DATABASE zemenlink_global;

\c zemenlink_global

CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    domain VARCHAR(255) UNIQUE NOT NULL,
    db_connection_string TEXT NOT NULL, -- Encrypted at rest
    s3_bucket_name VARCHAR(255) NOT NULL,
    compliance_level VARCHAR(50) DEFAULT 'standard', -- high, standard, etc.
    oidc_config JSONB, -- SAML/OIDC provider details
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 2. Tenant Database Schema (one per tenant)
-- This schema is applied to each tenant's dedicated database.

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    external_id VARCHAR(255) UNIQUE NOT NULL, -- ID from IdP (Okta/Azure)
    email VARCHAR(255) UNIQUE NOT NULL,
    full_name VARCHAR(255),
    avatar_url TEXT,
    role VARCHAR(50) DEFAULT 'user', -- user, admin, compliance_officer
    is_active BOOLEAN DEFAULT TRUE,
    last_seen_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE chats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type VARCHAR(50) NOT NULL, -- '1:1', 'group'
    title VARCHAR(255),
    metadata JSONB, -- For group icons, descriptions, etc.
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE chat_participants (
    chat_id UUID REFERENCES chats(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    role VARCHAR(50) DEFAULT 'member', -- member, admin
    PRIMARY KEY (chat_id, user_id)
);

CREATE TABLE messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    chat_id UUID REFERENCES chats(id) ON DELETE CASCADE,
    sender_id UUID REFERENCES users(id) ON DELETE SET NULL,
    content_encrypted TEXT NOT NULL, -- The message content, encrypted once with a symmetric key
    content_type VARCHAR(50) DEFAULT 'text', -- text, file, image
    metadata JSONB, -- For threading (parent_id), file info (s3_key, size), reactions
    status VARCHAR(50) DEFAULT 'sent', -- sent, delivered, read
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Separate table for encrypted symmetric keys to allow multiple recipients + escrow
CREATE TABLE message_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    message_id UUID NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
    recipient_id UUID, -- User ID for standard recipient, NULL for corporate escrow key
    key_encrypted TEXT NOT NULL, -- The symmetric key, encrypted with the recipient's/escrow's public key
    is_escrow BOOLEAN DEFAULT FALSE,
    -- Ensure only one escrow key per message
    CONSTRAINT unique_message_escrow UNIQUE (message_id, is_escrow) WHERE (is_escrow = TRUE),
    -- Ensure one key per recipient per message
    CONSTRAINT unique_message_recipient UNIQUE (message_id, recipient_id) WHERE (is_escrow = FALSE)
);

CREATE INDEX idx_messages_chat_id ON messages(chat_id);
CREATE INDEX idx_messages_created_at ON messages(created_at);

-- 3. Compliance & Audit Tables (per tenant)

CREATE TABLE escrow_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    requester_id UUID REFERENCES users(id),
    target_type VARCHAR(50), -- 'user', 'chat'
    target_id UUID,
    reason TEXT NOT NULL,
    status VARCHAR(50) DEFAULT 'pending', -- pending, approved, rejected, expired
    approvals_required INTEGER DEFAULT 2,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE escrow_approvals (
    request_id UUID REFERENCES escrow_requests(id) ON DELETE CASCADE,
    approver_id UUID REFERENCES users(id),
    approved_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (request_id, approver_id)
);

CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    action VARCHAR(255) NOT NULL,
    resource_type VARCHAR(100),
    resource_id UUID,
    details JSONB,
    ip_address INET,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Enable UUID generation
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Domain Enums
CREATE TYPE idea_category AS ENUM (
    'saas',
    'marketplace',
    'consumer',
    'devtool',
    'hardware',
    'fintech',
    'health',
    'education',
    'social',
    'other'
);

CREATE TYPE death_cause AS ENUM (
    'no_market_need',
    'ran_out_of_money',
    'ran_out_of_time',
    'lost_motivation',
    'competitor',
    'technical_difficulty',
    'cofounder_conflict',
    'legal_regulatory',
    'bad_timing',
    'pivoted_away',
    'other'
);

CREATE TYPE validation_level AS ENUM (
    'none',
    'talked_to_users',
    'prototype',
    'paying_customers'
);

CREATE TYPE idea_status AS ENUM (
    'buried',
    'on_life_support',
    'resurrected'
);

-- Users
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    display_name TEXT NOT NULL,
    is_demo BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Refresh tokens for rotation & family revocation
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL UNIQUE,
    family_id UUID NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_family_id ON refresh_tokens(family_id);

-- Ideas ("the graves")
CREATE TABLE ideas (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    tagline TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    category idea_category NOT NULL,
    born_at DATE NOT NULL,
    died_at DATE NOT NULL CHECK (died_at >= born_at),
    primary_cause death_cause NOT NULL,
    secondary_causes death_cause[] NOT NULL DEFAULT '{}',
    epitaph VARCHAR(140) NOT NULL,
    autopsy TEXT NOT NULL DEFAULT '',
    validation_level validation_level NOT NULL DEFAULT 'none',
    money_spent_cents BIGINT NOT NULL DEFAULT 0 CHECK (money_spent_cents >= 0),
    hours_invested INT NOT NULL DEFAULT 0 CHECK (hours_invested >= 0),
    status idea_status NOT NULL DEFAULT 'buried',
    would_resurrect BOOLEAN NOT NULL DEFAULT FALSE,
    search_vector TSVECTOR GENERATED ALWAYS AS (
        to_tsvector(
            'english',
            coalesce(name, '') || ' ' ||
            coalesce(tagline, '') || ' ' ||
            coalesce(description, '') || ' ' ||
            coalesce(autopsy, '')
        )
    ) STORED,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_ideas_user_id ON ideas(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_ideas_user_died_at ON ideas(user_id, died_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_ideas_user_primary_cause ON ideas(user_id, primary_cause) WHERE deleted_at IS NULL;
CREATE INDEX idx_ideas_search_gin ON ideas USING GIN (search_vector) WHERE deleted_at IS NULL;

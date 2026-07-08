CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    original_filename TEXT NOT NULL,
    stored_filename TEXT NOT NULL,
    object_key TEXT NOT NULL,
    bucket_name TEXT NOT NULL,
    content_type TEXT NOT NULL,
    file_size BIGINT NOT NULL,
    width INTEGER NOT NULL,
    height INTEGER NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'uploaded',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_images_status
        CHECK (status IN ('uploaded', 'queued', 'processing', 'completed', 'failed')),
    CONSTRAINT chk_images_file_size_non_negative
        CHECK (file_size >= 0),
    CONSTRAINT chk_images_width_non_negative
        CHECK (width >= 0),
    CONSTRAINT chk_images_height_non_negative
        CHECK (height >= 0)
);

CREATE INDEX IF NOT EXISTS idx_images_object_key
    ON images (object_key);

CREATE INDEX IF NOT EXISTS idx_images_status
    ON images (status);

CREATE INDEX IF NOT EXISTS idx_images_created_at
    ON images (created_at);

-- +goose Up
-- +goose StatementBegin

-- 1. Permissions table — atomic capability flags
CREATE TABLE IF NOT EXISTS permissions (
    id BIGSERIAL PRIMARY KEY,
    slug VARCHAR(100) NOT NULL UNIQUE,
    category VARCHAR(50) NOT NULL,
    label VARCHAR(100) NOT NULL,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_permissions_category ON permissions(category);
COMMENT ON TABLE permissions IS 'Atomic capability flags; roles are templates of permissions';

-- 2. Roles table
CREATE TABLE IF NOT EXISTS roles (
    id BIGSERIAL PRIMARY KEY,
    slug VARCHAR(50) NOT NULL UNIQUE,
    label VARCHAR(100) NOT NULL,
    description TEXT,
    hierarchy_level INT NOT NULL DEFAULT 0,
    is_system_role BOOLEAN NOT NULL DEFAULT FALSE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE roles IS 'Identity templates; each role maps to a set of permissions';

-- 3. Role ↔ Permission mapping
CREATE TABLE IF NOT EXISTS role_permissions (
    role_id BIGINT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id BIGINT NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    granted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (role_id, permission_id)
);

CREATE INDEX idx_role_permissions_role ON role_permissions(role_id);
CREATE INDEX idx_role_permissions_permission ON role_permissions(permission_id);
COMMENT ON TABLE role_permissions IS 'Many-to-many: which permissions each role grants';

-- 4. User ↔ Role assignments
CREATE TABLE IF NOT EXISTS user_roles (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id BIGINT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    assigned_by UUID REFERENCES users(id),
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    UNIQUE(user_id, role_id)
);

CREATE INDEX idx_user_roles_user ON user_roles(user_id);
CREATE INDEX idx_user_roles_role ON user_roles(role_id);
COMMENT ON TABLE user_roles IS 'Per-user role assignments; supports temporary roles with expiry';

-- 5. User permission overrides
CREATE TABLE IF NOT EXISTS user_permission_overrides (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    permission_id BIGINT NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    granted BOOLEAN NOT NULL DEFAULT TRUE,
    reason TEXT,
    granted_by UUID REFERENCES users(id),
    granted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ,
    UNIQUE(user_id, permission_id)
);

CREATE INDEX idx_user_permission_overrides_user ON user_permission_overrides(user_id);
COMMENT ON TABLE user_permission_overrides IS 'Per-user permission overrides; can grant or revoke beyond role';

-- 6. Seed default roles
INSERT INTO roles (slug, label, description, hierarchy_level, is_system_role) VALUES
    ('listener',   'Listener',   'Default role for registered users', 10, TRUE),
    ('creator',    'Creator',    'Artists and uploaders with creation privileges', 30, TRUE),
    ('curator',    'Curator',    'Playlist curators and tastemakers', 50, TRUE),
    ('moderator',  'Moderator',  'Content and community moderation', 80, TRUE),
    ('admin',      'Admin',      'Full platform access', 100, TRUE)
ON CONFLICT (slug) DO NOTHING;

-- 7. Seed default permissions
INSERT INTO permissions (slug, category, label, description) VALUES
    ('upload_track',            'upload',       'Upload Tracks',            'Submit new tracks for review'),
    ('upload_album',            'upload',       'Upload Albums',            'Submit complete albums'),
    ('upload_artwork',          'upload',       'Upload Artwork',           'Upload cover art and images'),
    ('publish_immediately',     'upload',       'Publish Immediately',      'Skip review queue (trusted creators)'),
    ('edit_metadata',           'contribution', 'Edit Metadata',            'Suggest or apply metadata edits'),
    ('edit_lyrics',             'contribution', 'Edit Lyrics',              'Contribute or correct lyrics'),
    ('submit_to_club',          'contribution', 'Submit to Club',           'Submit tracks to music clubs'),
    ('review_uploads',          'moderation',   'Review Uploads',           'Review and approve/reject uploads'),
    ('review_metadata',         'moderation',   'Review Metadata',          'Review and approve metadata edits'),
    ('moderate_comments',       'moderation',   'Moderate Comments',        'Remove or flag comments'),
    ('moderate_users',          'moderation',   'Moderate Users',           'Warn, mute, or ban users'),
    ('moderate_clubs',          'moderation',   'Moderate Clubs',           'Manage and moderate music clubs'),
    ('create_club',             'social',       'Create Club',              'Create and manage music clubs'),
    ('create_playlist',         'social',       'Create Playlist',          'Create public playlists'),
    ('send_invitations',        'social',       'Send Invitations',         'Invite users to clubs or events'),
    ('access_creator_dashboard','creator',      'Creator Dashboard',        'Access creator analytics and tools'),
    ('manage_releases',         'creator',      'Manage Releases',          'Manage published releases'),
    ('enable_monetization',     'creator',      'Enable Monetization',      'Access monetization features'),
    ('manage_users',            'admin',        'Manage Users',             'Full user management'),
    ('manage_roles',            'admin',        'Manage Roles',             'Create and assign roles'),
    ('manage_system',           'admin',        'Manage System',            'System configuration and maintenance'),
    ('view_audit_log',          'admin',        'View Audit Log',           'Access audit trail')
ON CONFLICT (slug) DO NOTHING;

-- 8. Assign permissions to roles
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE r.slug = 'listener' AND p.slug IN ('create_playlist', 'submit_to_club')
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE r.slug = 'creator' AND p.slug IN (
    'upload_track', 'upload_album', 'upload_artwork',
    'edit_metadata', 'edit_lyrics',
    'access_creator_dashboard', 'manage_releases', 'create_playlist', 'submit_to_club'
)
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE r.slug = 'curator' AND p.slug IN (
    'upload_track', 'upload_album', 'upload_artwork',
    'edit_metadata', 'edit_lyrics', 'publish_immediately',
    'access_creator_dashboard', 'manage_releases', 'create_playlist', 'submit_to_club',
    'create_club', 'send_invitations'
)
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE r.slug = 'moderator' AND p.slug IN (
    'upload_track', 'upload_album', 'upload_artwork',
    'edit_metadata', 'edit_lyrics', 'publish_immediately',
    'access_creator_dashboard', 'manage_releases', 'create_playlist', 'submit_to_club',
    'create_club', 'send_invitations',
    'review_uploads', 'review_metadata',
    'moderate_comments', 'moderate_users', 'moderate_clubs'
)
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE r.slug = 'admin'
ON CONFLICT DO NOTHING;

-- 9. Backfill user_roles from existing users.role column
INSERT INTO user_roles (user_id, role_id, assigned_at)
SELECT u.id, r.id, u.created_at
FROM users u JOIN roles r ON r.slug = u.role
ON CONFLICT (user_id, role_id) DO NOTHING;

-- 10. Updated_at trigger
CREATE OR REPLACE FUNCTION update_rbac_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_permissions_updated_at
    BEFORE UPDATE ON permissions
    FOR EACH ROW EXECUTE FUNCTION update_rbac_updated_at();

CREATE TRIGGER trigger_roles_updated_at
    BEFORE UPDATE ON roles
    FOR EACH ROW EXECUTE FUNCTION update_rbac_updated_at();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS trigger_roles_updated_at ON roles;
DROP TRIGGER IF EXISTS trigger_permissions_updated_at ON permissions;
DROP FUNCTION IF EXISTS update_rbac_updated_at();
DROP TABLE IF EXISTS user_permission_overrides;
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS permissions;

-- +goose StatementEnd

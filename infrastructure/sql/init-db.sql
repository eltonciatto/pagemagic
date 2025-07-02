-- Page Magic Database Schema

-- Users and Authentication
CREATE TABLE users (
    id VARCHAR(255) PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    avatar_url VARCHAR(500),
    provider VARCHAR(50) DEFAULT 'magic_link',
    provider_id VARCHAR(255),
    email_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE magic_links (
    id VARCHAR(255) PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    token VARCHAR(500) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    used BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE user_sessions (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(500) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Projects and Sites
CREATE TABLE projects (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) DEFAULT 'draft',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE sites (
    id VARCHAR(255) PRIMARY KEY,
    project_id VARCHAR(255) REFERENCES projects(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    domain VARCHAR(255),
    subdomain VARCHAR(100),
    ast_data JSONB,
    css_data TEXT,
    js_data TEXT,
    i18n_data JSONB,
    status VARCHAR(50) DEFAULT 'draft',
    published_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE site_builds (
    id VARCHAR(255) PRIMARY KEY,
    site_id VARCHAR(255) REFERENCES sites(id) ON DELETE CASCADE,
    build_number INTEGER NOT NULL,
    image_id VARCHAR(255),
    status VARCHAR(50) DEFAULT 'pending',
    build_logs TEXT,
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP
);

-- AI Generations
CREATE TABLE generations (
    id VARCHAR(255) PRIMARY KEY,
    project_id VARCHAR(255) REFERENCES projects(id) ON DELETE CASCADE,
    prompt TEXT NOT NULL,
    model VARCHAR(100),
    locale VARCHAR(10) DEFAULT 'en',
    status VARCHAR(50) DEFAULT 'queued',
    result_data JSONB,
    tokens_used INTEGER DEFAULT 0,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP
);

-- Domains
CREATE TABLE domains (
    id VARCHAR(255) PRIMARY KEY,
    site_id VARCHAR(255) REFERENCES sites(id) ON DELETE CASCADE,
    domain VARCHAR(255) UNIQUE NOT NULL,
    registrar VARCHAR(100),
    dns_provider VARCHAR(100),
    ssl_status VARCHAR(50) DEFAULT 'pending',
    ssl_expires_at TIMESTAMP,
    status VARCHAR(50) DEFAULT 'pending',
    purchased_at TIMESTAMP,
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Billing and Subscriptions
CREATE TABLE subscriptions (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) REFERENCES users(id) ON DELETE CASCADE,
    stripe_customer_id VARCHAR(255) NOT NULL,
    stripe_subscription_id VARCHAR(255),
    plan_id VARCHAR(255) NOT NULL,
    status VARCHAR(50) DEFAULT 'active',
    current_period_start TIMESTAMP,
    current_period_end TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE usage_records (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) REFERENCES users(id) ON DELETE CASCADE,
    metric_name VARCHAR(100) NOT NULL,
    quantity BIGINT NOT NULL,
    period VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE invoices (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) REFERENCES users(id) ON DELETE CASCADE,
    stripe_invoice_id VARCHAR(255) NOT NULL,
    amount BIGINT,
    status VARCHAR(50),
    due_date TIMESTAMP,
    paid_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Container Hosting
CREATE TABLE container_services (
    id VARCHAR(255) PRIMARY KEY,
    site_id VARCHAR(255) REFERENCES sites(id) ON DELETE CASCADE,
    service_name VARCHAR(255) NOT NULL,
    image_name VARCHAR(255) NOT NULL,
    status VARCHAR(50) DEFAULT 'stopped',
    replica_count INTEGER DEFAULT 1,
    cpu_limit VARCHAR(20),
    memory_limit VARCHAR(20),
    environment_vars JSONB,
    ports JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- I18n Translations
CREATE TABLE translation_bundles (
    id VARCHAR(255) PRIMARY KEY,
    project_id VARCHAR(255) REFERENCES projects(id) ON DELETE CASCADE,
    locale VARCHAR(10) NOT NULL,
    bundle_data JSONB NOT NULL,
    version INTEGER DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Mobile Push Notifications
CREATE TABLE mobile_push_tokens (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(500) NOT NULL,
    platform VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_magic_links_token ON magic_links(token);
CREATE INDEX idx_projects_user_id ON projects(user_id);
CREATE INDEX idx_sites_project_id ON sites(project_id);
CREATE INDEX idx_sites_domain ON sites(domain);
CREATE INDEX idx_generations_project_id ON generations(project_id);
CREATE INDEX idx_subscriptions_user_id ON subscriptions(user_id);
CREATE INDEX idx_usage_records_user_id_period ON usage_records(user_id, period);
CREATE INDEX idx_container_services_site_id ON container_services(site_id);
CREATE INDEX idx_translation_bundles_project_locale ON translation_bundles(project_id, locale);

-- Functions to auto-update timestamps
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Triggers for auto-updating timestamps
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_projects_updated_at BEFORE UPDATE ON projects FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_sites_updated_at BEFORE UPDATE ON sites FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_subscriptions_updated_at BEFORE UPDATE ON subscriptions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_container_services_updated_at BEFORE UPDATE ON container_services FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

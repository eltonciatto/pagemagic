-- Page Magic Database Initialization
-- PostgreSQL 16 Schema

-- Extensões necessárias
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS "btree_gin";

-- Função para atualizar updated_at automaticamente
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = CURRENT_TIMESTAMP;
  RETURN NEW;
END;
$$ language 'plpgsql';

-- ==========================================
-- USERS & AUTHENTICATION
-- ==========================================

-- Enum para tipos de autenticação
CREATE TYPE auth_provider AS ENUM ('email', 'google', 'github', 'apple');

-- Enum para status do usuário
CREATE TYPE user_status AS ENUM ('active', 'inactive', 'suspended', 'deleted');

-- Enum para planos
CREATE TYPE subscription_plan AS ENUM ('free', 'basic', 'pro', 'enterprise');

-- Tabela de usuários
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  email VARCHAR(255) UNIQUE NOT NULL,
  email_verified BOOLEAN DEFAULT FALSE,
  password_hash VARCHAR(255), -- NULL para OAuth users
  first_name VARCHAR(100),
  last_name VARCHAR(100),
  avatar_url TEXT,
  status user_status DEFAULT 'active',
  locale VARCHAR(10) DEFAULT 'en',
  timezone VARCHAR(50) DEFAULT 'UTC',
  last_login_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Trigger para updated_at
CREATE TRIGGER update_users_updated_at 
  BEFORE UPDATE ON users 
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Índices para performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_created_at ON users(created_at);

-- Tabela de autenticação externa (OAuth)
CREATE TABLE user_auth_providers (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  provider auth_provider NOT NULL,
  provider_user_id VARCHAR(255) NOT NULL,
  provider_email VARCHAR(255),
  provider_data JSONB,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  
  UNIQUE(provider, provider_user_id)
);

CREATE INDEX idx_user_auth_providers_user_id ON user_auth_providers(user_id);
CREATE INDEX idx_user_auth_providers_provider ON user_auth_providers(provider, provider_user_id);

-- Tabela de magic links
CREATE TABLE magic_links (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  token VARCHAR(255) UNIQUE NOT NULL,
  expires_at TIMESTAMPTZ NOT NULL,
  used_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_magic_links_token ON magic_links(token);
CREATE INDEX idx_magic_links_expires_at ON magic_links(expires_at);

-- ==========================================
-- SUBSCRIPTIONS & BILLING
-- ==========================================

-- Tabela de assinaturas
CREATE TABLE subscriptions (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  stripe_subscription_id VARCHAR(255) UNIQUE,
  stripe_customer_id VARCHAR(255),
  plan subscription_plan NOT NULL DEFAULT 'free',
  status VARCHAR(50) DEFAULT 'active',
  current_period_start TIMESTAMPTZ,
  current_period_end TIMESTAMPTZ,
  cancel_at_period_end BOOLEAN DEFAULT FALSE,
  trial_start TIMESTAMPTZ,
  trial_end TIMESTAMPTZ,
  metadata JSONB,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER update_subscriptions_updated_at 
  BEFORE UPDATE ON subscriptions 
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE INDEX idx_subscriptions_user_id ON subscriptions(user_id);
CREATE INDEX idx_subscriptions_stripe_subscription_id ON subscriptions(stripe_subscription_id);
CREATE INDEX idx_subscriptions_status ON subscriptions(status);

-- Tabela de uso (Stripe Meters)
CREATE TABLE usage_records (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  subscription_id UUID REFERENCES subscriptions(id) ON DELETE SET NULL,
  meter_type VARCHAR(50) NOT NULL, -- page_generate, ai_token, container_hours, storage_gb
  quantity DECIMAL(10,2) NOT NULL,
  timestamp TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  stripe_usage_record_id VARCHAR(255),
  metadata JSONB,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_usage_records_user_id ON usage_records(user_id);
CREATE INDEX idx_usage_records_meter_type ON usage_records(meter_type);
CREATE INDEX idx_usage_records_timestamp ON usage_records(timestamp);
CREATE INDEX idx_usage_records_stripe_usage_record_id ON usage_records(stripe_usage_record_id);

-- ==========================================
-- PROJECTS & SITES
-- ==========================================

-- Enum para status do projeto
CREATE TYPE project_status AS ENUM ('draft', 'building', 'published', 'archived', 'error');

-- Enum para tipos de site
CREATE TYPE site_type AS ENUM ('landing_page', 'blog', 'portfolio', 'ecommerce', 'custom');

-- Tabela de projetos
CREATE TABLE projects (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  type site_type DEFAULT 'landing_page',
  status project_status DEFAULT 'draft',
  domain VARCHAR(255),
  subdomain VARCHAR(100),
  custom_domain VARCHAR(255),
  ssl_enabled BOOLEAN DEFAULT TRUE,
  password_protected BOOLEAN DEFAULT FALSE,
  password_hash VARCHAR(255),
  seo_title VARCHAR(255),
  seo_description TEXT,
  seo_keywords TEXT[],
  favicon_url TEXT,
  social_image_url TEXT,
  analytics_enabled BOOLEAN DEFAULT TRUE,
  analytics_tracking_id VARCHAR(100),
  settings JSONB DEFAULT '{}',
  metadata JSONB DEFAULT '{}',
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  published_at TIMESTAMPTZ,
  archived_at TIMESTAMPTZ
);

CREATE TRIGGER update_projects_updated_at 
  BEFORE UPDATE ON projects 
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE INDEX idx_projects_user_id ON projects(user_id);
CREATE INDEX idx_projects_status ON projects(status);
CREATE INDEX idx_projects_domain ON projects(domain);
CREATE INDEX idx_projects_subdomain ON projects(subdomain);
CREATE INDEX idx_projects_custom_domain ON projects(custom_domain);
CREATE INDEX idx_projects_created_at ON projects(created_at);

-- Tabela de versões/builds
CREATE TABLE project_builds (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
  version_number INTEGER NOT NULL,
  docker_image_id VARCHAR(255),
  build_log TEXT,
  build_status VARCHAR(50) DEFAULT 'pending', -- pending, building, success, error
  build_started_at TIMESTAMPTZ,
  build_completed_at TIMESTAMPTZ,
  is_current BOOLEAN DEFAULT FALSE,
  size_bytes BIGINT,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_project_builds_project_id ON project_builds(project_id);
CREATE INDEX idx_project_builds_version_number ON project_builds(project_id, version_number);
CREATE INDEX idx_project_builds_is_current ON project_builds(project_id, is_current);

-- ==========================================
-- CONTENT & PAGES
-- ==========================================

-- Enum para tipos de seção
CREATE TYPE section_type AS ENUM ('hero', 'about', 'services', 'portfolio', 'testimonials', 'contact', 'footer', 'custom');

-- Tabela de páginas
CREATE TABLE pages (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
  name VARCHAR(255) NOT NULL,
  slug VARCHAR(255) NOT NULL,
  title VARCHAR(255),
  description TEXT,
  is_home BOOLEAN DEFAULT FALSE,
  is_published BOOLEAN DEFAULT TRUE,
  order_index INTEGER DEFAULT 0,
  seo_title VARCHAR(255),
  seo_description TEXT,
  seo_keywords TEXT[],
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  
  UNIQUE(project_id, slug)
);

CREATE TRIGGER update_pages_updated_at 
  BEFORE UPDATE ON pages 
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE INDEX idx_pages_project_id ON pages(project_id);
CREATE INDEX idx_pages_slug ON pages(project_id, slug);
CREATE INDEX idx_pages_is_home ON pages(project_id, is_home);

-- Tabela de seções
CREATE TABLE sections (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  page_id UUID NOT NULL REFERENCES pages(id) ON DELETE CASCADE,
  type section_type NOT NULL,
  name VARCHAR(255) NOT NULL,
  order_index INTEGER DEFAULT 0,
  html_content TEXT,
  css_styles TEXT,
  js_scripts TEXT,
  data JSONB DEFAULT '{}',
  settings JSONB DEFAULT '{}',
  is_visible BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER update_sections_updated_at 
  BEFORE UPDATE ON sections 
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE INDEX idx_sections_page_id ON sections(page_id);
CREATE INDEX idx_sections_type ON sections(type);
CREATE INDEX idx_sections_order_index ON sections(page_id, order_index);

-- ==========================================
-- AI GENERATIONS
-- ==========================================

-- Enum para status de geração
CREATE TYPE generation_status AS ENUM ('queued', 'processing', 'completed', 'failed', 'cancelled');

-- Enum para modelos de IA
CREATE TYPE ai_model AS ENUM ('llama-3-70b', 'gpt-4', 'claude-3-sonnet', 'gpt-3.5-turbo');

-- Tabela de gerações de IA
CREATE TABLE ai_generations (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
  prompt TEXT NOT NULL,
  model ai_model NOT NULL,
  status generation_status DEFAULT 'queued',
  response_data JSONB,
  tokens_used INTEGER,
  processing_time_ms INTEGER,
  error_message TEXT,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  completed_at TIMESTAMPTZ,
  metadata JSONB DEFAULT '{}'
);

CREATE INDEX idx_ai_generations_user_id ON ai_generations(user_id);
CREATE INDEX idx_ai_generations_project_id ON ai_generations(project_id);
CREATE INDEX idx_ai_generations_status ON ai_generations(status);
CREATE INDEX idx_ai_generations_created_at ON ai_generations(created_at);
CREATE INDEX idx_ai_generations_model ON ai_generations(model);

-- ==========================================
-- MEDIA & ASSETS
-- ==========================================

-- Enum para tipos de mídia
CREATE TYPE media_type AS ENUM ('image', 'video', 'audio', 'document', 'font', 'other');

-- Tabela de arquivos de mídia
CREATE TABLE media_files (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
  filename VARCHAR(255) NOT NULL,
  original_filename VARCHAR(255) NOT NULL,
  mime_type VARCHAR(100) NOT NULL,
  type media_type NOT NULL,
  size_bytes BIGINT NOT NULL,
  width INTEGER,
  height INTEGER,
  duration_seconds INTEGER,
  url TEXT NOT NULL,
  thumbnail_url TEXT,
  alt_text TEXT,
  description TEXT,
  metadata JSONB DEFAULT '{}',
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_media_files_user_id ON media_files(user_id);
CREATE INDEX idx_media_files_project_id ON media_files(project_id);
CREATE INDEX idx_media_files_type ON media_files(type);
CREATE INDEX idx_media_files_created_at ON media_files(created_at);

-- ==========================================
-- DOMAINS & SSL
-- ==========================================

-- Enum para status do domínio
CREATE TYPE domain_status AS ENUM ('pending', 'active', 'expired', 'suspended', 'error');

-- Tabela de domínios customizados
CREATE TABLE custom_domains (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
  domain VARCHAR(255) UNIQUE NOT NULL,
  status domain_status DEFAULT 'pending',
  dns_configured BOOLEAN DEFAULT FALSE,
  ssl_certificate_id VARCHAR(255),
  ssl_issued_at TIMESTAMPTZ,
  ssl_expires_at TIMESTAMPTZ,
  verified_at TIMESTAMPTZ,
  last_check_at TIMESTAMPTZ,
  error_message TEXT,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER update_custom_domains_updated_at 
  BEFORE UPDATE ON custom_domains 
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE INDEX idx_custom_domains_project_id ON custom_domains(project_id);
CREATE INDEX idx_custom_domains_domain ON custom_domains(domain);
CREATE INDEX idx_custom_domains_status ON custom_domains(status);

-- ==========================================
-- ANALYTICS & TRACKING
-- ==========================================

-- Tabela de eventos de analytics (resumido - dados detalhados no TimescaleDB)
CREATE TABLE analytics_events (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
  event_type VARCHAR(50) NOT NULL,
  event_count INTEGER DEFAULT 1,
  event_date DATE NOT NULL,
  metadata JSONB DEFAULT '{}',
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  
  UNIQUE(project_id, event_type, event_date)
);

CREATE INDEX idx_analytics_events_project_id ON analytics_events(project_id);
CREATE INDEX idx_analytics_events_event_date ON analytics_events(event_date);
CREATE INDEX idx_analytics_events_event_type ON analytics_events(event_type);

-- ==========================================
-- MOBILE PUSH NOTIFICATIONS
-- ==========================================

-- Enum para plataformas mobile
CREATE TYPE mobile_platform AS ENUM ('ios', 'android', 'expo');

-- Tabela de tokens push
CREATE TABLE mobile_push_tokens (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  token VARCHAR(512) NOT NULL,
  platform mobile_platform NOT NULL,
  app_version VARCHAR(50),
  device_id VARCHAR(255),
  is_active BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  
  UNIQUE(user_id, token, platform)
);

CREATE TRIGGER update_mobile_push_tokens_updated_at 
  BEFORE UPDATE ON mobile_push_tokens 
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE INDEX idx_mobile_push_tokens_user_id ON mobile_push_tokens(user_id);
CREATE INDEX idx_mobile_push_tokens_platform ON mobile_push_tokens(platform);
CREATE INDEX idx_mobile_push_tokens_is_active ON mobile_push_tokens(is_active);

-- Tabela de notificações enviadas
CREATE TABLE push_notifications (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  title VARCHAR(255) NOT NULL,
  body TEXT NOT NULL,
  data JSONB DEFAULT '{}',
  sent_at TIMESTAMPTZ,
  delivered_at TIMESTAMPTZ,
  clicked_at TIMESTAMPTZ,
  error_message TEXT,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_push_notifications_user_id ON push_notifications(user_id);
CREATE INDEX idx_push_notifications_sent_at ON push_notifications(sent_at);

-- ==========================================
-- INTERNATIONALIZATION
-- ==========================================

-- Tabela de idiomas suportados
CREATE TABLE languages (
  code VARCHAR(10) PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  native_name VARCHAR(100) NOT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  is_rtl BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Inserir idiomas padrão
INSERT INTO languages (code, name, native_name) VALUES
  ('en', 'English', 'English'),
  ('pt', 'Portuguese', 'Português'),
  ('es', 'Spanish', 'Español'),
  ('fr', 'French', 'Français'),
  ('de', 'German', 'Deutsch'),
  ('it', 'Italian', 'Italiano'),
  ('ja', 'Japanese', '日本語'),
  ('ko', 'Korean', '한국어'),
  ('zh', 'Chinese', '中文'),
  ('ar', 'Arabic', 'العربية');

-- Tabela de traduções
CREATE TABLE translations (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  key VARCHAR(255) NOT NULL,
  language_code VARCHAR(10) NOT NULL REFERENCES languages(code) ON DELETE CASCADE,
  value TEXT NOT NULL,
  context VARCHAR(255),
  project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  
  UNIQUE(key, language_code, project_id)
);

CREATE TRIGGER update_translations_updated_at 
  BEFORE UPDATE ON translations 
  FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE INDEX idx_translations_key ON translations(key);
CREATE INDEX idx_translations_language_code ON translations(language_code);
CREATE INDEX idx_translations_project_id ON translations(project_id);

-- ==========================================
-- AUDIT LOG
-- ==========================================

-- Enum para tipos de ação
CREATE TYPE audit_action AS ENUM ('create', 'update', 'delete', 'login', 'logout', 'publish', 'unpublish');

-- Tabela de auditoria
CREATE TABLE audit_logs (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID REFERENCES users(id) ON DELETE SET NULL,
  action audit_action NOT NULL,
  resource_type VARCHAR(100) NOT NULL,
  resource_id UUID,
  old_values JSONB,
  new_values JSONB,
  ip_address INET,
  user_agent TEXT,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_resource ON audit_logs(resource_type, resource_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);

-- ==========================================
-- CONFIGURAÇÕES FINAIS
-- ==========================================

-- Row Level Security (RLS) - Exemplos para tabelas principais
ALTER TABLE projects ENABLE ROW LEVEL SECURITY;
ALTER TABLE pages ENABLE ROW LEVEL SECURITY;
ALTER TABLE sections ENABLE ROW LEVEL SECURITY;
ALTER TABLE media_files ENABLE ROW LEVEL SECURITY;

-- Política RLS para projetos (usuários só veem seus próprios projetos)
CREATE POLICY projects_user_isolation ON projects
  FOR ALL
  TO authenticated
  USING (user_id = current_setting('app.current_user_id')::UUID);

-- Política RLS para páginas (através do projeto)
CREATE POLICY pages_user_isolation ON pages
  FOR ALL
  TO authenticated
  USING (project_id IN (
    SELECT id FROM projects WHERE user_id = current_setting('app.current_user_id')::UUID
  ));

-- Política RLS para seções (através da página)
CREATE POLICY sections_user_isolation ON sections
  FOR ALL
  TO authenticated
  USING (page_id IN (
    SELECT p.id FROM pages p
    JOIN projects pr ON p.project_id = pr.id
    WHERE pr.user_id = current_setting('app.current_user_id')::UUID
  ));

-- Política RLS para arquivos de mídia
CREATE POLICY media_files_user_isolation ON media_files
  FOR ALL
  TO authenticated
  USING (user_id = current_setting('app.current_user_id')::UUID);

-- Views úteis
CREATE VIEW user_project_stats AS
SELECT 
  u.id as user_id,
  u.email,
  COUNT(DISTINCT p.id) as total_projects,
  COUNT(DISTINCT CASE WHEN p.status = 'published' THEN p.id END) as published_projects,
  COUNT(DISTINCT pg.id) as total_pages,
  COUNT(DISTINCT mf.id) as total_media_files,
  SUM(mf.size_bytes) as total_storage_bytes,
  MAX(p.updated_at) as last_project_update
FROM users u
LEFT JOIN projects p ON u.id = p.user_id
LEFT JOIN pages pg ON p.id = pg.project_id
LEFT JOIN media_files mf ON u.id = mf.user_id
GROUP BY u.id, u.email;

-- Índices compostos para performance
CREATE INDEX idx_projects_user_status ON projects(user_id, status);
CREATE INDEX idx_pages_project_published ON pages(project_id, is_published);
CREATE INDEX idx_sections_page_visible_order ON sections(page_id, is_visible, order_index);
CREATE INDEX idx_usage_records_user_meter_timestamp ON usage_records(user_id, meter_type, timestamp);

-- Função para limpar dados antigos (cleanup job)
CREATE OR REPLACE FUNCTION cleanup_old_data()
RETURNS void AS $$
BEGIN
  -- Deletar magic links expirados há mais de 24h
  DELETE FROM magic_links 
  WHERE expires_at < NOW() - INTERVAL '24 hours';
  
  -- Deletar logs de auditoria com mais de 90 dias
  DELETE FROM audit_logs 
  WHERE created_at < NOW() - INTERVAL '90 days';
  
  -- Deletar notificações antigas (mais de 30 dias)
  DELETE FROM push_notifications 
  WHERE created_at < NOW() - INTERVAL '30 days';
  
END;
$$ LANGUAGE plpgsql;

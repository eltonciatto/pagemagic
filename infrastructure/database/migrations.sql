-- Page Magic Database Migrations
-- Versão 1.0.0 - Schema inicial

-- ==========================================
-- MIGRATION 001: Initial Schema
-- ==========================================

-- Criar schema se não existir
CREATE SCHEMA IF NOT EXISTS public;

-- Verificar se a migração já foi aplicada
DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM information_schema.tables 
    WHERE table_schema = 'public' AND table_name = 'schema_migrations'
  ) THEN
    
    -- Criar tabela de controle de migrações
    CREATE TABLE schema_migrations (
      version VARCHAR(50) PRIMARY KEY,
      applied_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
      description TEXT
    );
    
    INSERT INTO schema_migrations (version, description) 
    VALUES ('001', 'Initial schema with users, projects, and basic tables');
    
  END IF;
END
$$;

-- ==========================================
-- MIGRATION 002: Add mobile features
-- ==========================================

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM schema_migrations WHERE version = '002'
  ) THEN
    
    -- Adicionar colunas mobile à tabela users se não existirem
    ALTER TABLE users 
    ADD COLUMN IF NOT EXISTS mobile_verified BOOLEAN DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS mobile_number VARCHAR(20),
    ADD COLUMN IF NOT EXISTS push_notifications_enabled BOOLEAN DEFAULT TRUE;
    
    -- Adicionar índice para número de celular
    CREATE INDEX IF NOT EXISTS idx_users_mobile_number ON users(mobile_number) 
    WHERE mobile_number IS NOT NULL;
    
    INSERT INTO schema_migrations (version, description) 
    VALUES ('002', 'Add mobile features and push notification fields');
    
  END IF;
END
$$;

-- ==========================================
-- MIGRATION 003: Enhanced analytics
-- ==========================================

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM schema_migrations WHERE version = '003'
  ) THEN
    
    -- Adicionar campos de analytics avançados aos projetos
    ALTER TABLE projects 
    ADD COLUMN IF NOT EXISTS analytics_config JSONB DEFAULT '{}',
    ADD COLUMN IF NOT EXISTS conversion_goals JSONB DEFAULT '[]',
    ADD COLUMN IF NOT EXISTS ab_testing_enabled BOOLEAN DEFAULT FALSE;
    
    -- Adicionar campos de performance às seções
    ALTER TABLE sections 
    ADD COLUMN IF NOT EXISTS lazy_loading BOOLEAN DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS priority INTEGER DEFAULT 0,
    ADD COLUMN IF NOT EXISTS cache_duration_seconds INTEGER DEFAULT 3600;
    
    INSERT INTO schema_migrations (version, description) 
    VALUES ('003', 'Enhanced analytics and performance features');
    
  END IF;
END
$$;

-- ==========================================
-- MIGRATION 004: Team collaboration
-- ==========================================

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM schema_migrations WHERE version = '004'
  ) THEN
    
    -- Criar tabela de times
    CREATE TABLE IF NOT EXISTS teams (
      id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
      name VARCHAR(255) NOT NULL,
      slug VARCHAR(100) UNIQUE NOT NULL,
      owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
      plan subscription_plan DEFAULT 'free',
      settings JSONB DEFAULT '{}',
      created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
      updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
    );
    
    -- Criar trigger para updated_at
    CREATE TRIGGER update_teams_updated_at 
      BEFORE UPDATE ON teams 
      FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
    
    -- Criar tabela de membros do time
    CREATE TABLE IF NOT EXISTS team_members (
      id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
      team_id UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
      user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
      role VARCHAR(50) DEFAULT 'member', -- owner, admin, editor, viewer
      invited_by UUID REFERENCES users(id),
      invited_at TIMESTAMPTZ,
      joined_at TIMESTAMPTZ,
      created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
      
      UNIQUE(team_id, user_id)
    );
    
    -- Adicionar team_id aos projetos
    ALTER TABLE projects 
    ADD COLUMN IF NOT EXISTS team_id UUID REFERENCES teams(id) ON DELETE SET NULL;
    
    -- Índices
    CREATE INDEX IF NOT EXISTS idx_teams_owner_id ON teams(owner_id);
    CREATE INDEX IF NOT EXISTS idx_teams_slug ON teams(slug);
    CREATE INDEX IF NOT EXISTS idx_team_members_team_id ON team_members(team_id);
    CREATE INDEX IF NOT EXISTS idx_team_members_user_id ON team_members(user_id);
    CREATE INDEX IF NOT EXISTS idx_projects_team_id ON projects(team_id);
    
    INSERT INTO schema_migrations (version, description) 
    VALUES ('004', 'Team collaboration features');
    
  END IF;
END
$$;

-- ==========================================
-- MIGRATION 005: Enhanced AI features
-- ==========================================

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM schema_migrations WHERE version = '005'
  ) THEN
    
    -- Adicionar mais modelos de IA
    ALTER TYPE ai_model ADD VALUE IF NOT EXISTS 'claude-3-opus';
    ALTER TYPE ai_model ADD VALUE IF NOT EXISTS 'claude-3-haiku';
    ALTER TYPE ai_model ADD VALUE IF NOT EXISTS 'gpt-4-turbo';
    ALTER TYPE ai_model ADD VALUE IF NOT EXISTS 'llama-3-8b';
    
    -- Adicionar campos de contexto para IA
    ALTER TABLE ai_generations 
    ADD COLUMN IF NOT EXISTS context_data JSONB DEFAULT '{}',
    ADD COLUMN IF NOT EXISTS prompt_template VARCHAR(255),
    ADD COLUMN IF NOT EXISTS generation_mode VARCHAR(50) DEFAULT 'standard', -- standard, iterative, collaborative
    ADD COLUMN IF NOT EXISTS parent_generation_id UUID REFERENCES ai_generations(id);
    
    -- Criar tabela de templates de prompt
    CREATE TABLE IF NOT EXISTS prompt_templates (
      id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
      name VARCHAR(255) NOT NULL,
      description TEXT,
      category VARCHAR(100) NOT NULL,
      template TEXT NOT NULL,
      variables JSONB DEFAULT '[]',
      is_public BOOLEAN DEFAULT FALSE,
      created_by UUID REFERENCES users(id) ON DELETE SET NULL,
      usage_count INTEGER DEFAULT 0,
      rating DECIMAL(3,2) DEFAULT 0,
      created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
      updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
    );
    
    CREATE TRIGGER update_prompt_templates_updated_at 
      BEFORE UPDATE ON prompt_templates 
      FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
    
    CREATE INDEX IF NOT EXISTS idx_prompt_templates_category ON prompt_templates(category);
    CREATE INDEX IF NOT EXISTS idx_prompt_templates_is_public ON prompt_templates(is_public);
    CREATE INDEX IF NOT EXISTS idx_prompt_templates_created_by ON prompt_templates(created_by);
    
    INSERT INTO schema_migrations (version, description) 
    VALUES ('005', 'Enhanced AI features with templates and context');
    
  END IF;
END
$$;

-- ==========================================
-- MIGRATION 006: Advanced domains
-- ==========================================

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM schema_migrations WHERE version = '006'
  ) THEN
    
    -- Criar tabela de registradores de domínio
    CREATE TABLE IF NOT EXISTS domain_registrars (
      id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
      name VARCHAR(100) NOT NULL,
      api_endpoint TEXT NOT NULL,
      supported_tlds TEXT[] DEFAULT '{}',
      is_active BOOLEAN DEFAULT TRUE,
      credentials_schema JSONB,
      created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
    );
    
    -- Inserir registradores padrão
    INSERT INTO domain_registrars (name, api_endpoint, supported_tlds) VALUES
    ('Namecheap', 'https://api.namecheap.com/xml.response', ARRAY['.com', '.net', '.org', '.io', '.app']),
    ('GoDaddy', 'https://api.godaddy.com', ARRAY['.com', '.net', '.org', '.co', '.me']),
    ('Cloudflare', 'https://api.cloudflare.com/client/v4', ARRAY['.com', '.net', '.org', '.io']);
    
    -- Adicionar campos ao custom_domains
    ALTER TABLE custom_domains 
    ADD COLUMN IF NOT EXISTS registrar_id UUID REFERENCES domain_registrars(id),
    ADD COLUMN IF NOT EXISTS auto_renew BOOLEAN DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS expires_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS nameservers TEXT[] DEFAULT '{}',
    ADD COLUMN IF NOT EXISTS dns_records JSONB DEFAULT '[]';
    
    -- Criar tabela de histórico DNS
    CREATE TABLE IF NOT EXISTS dns_history (
      id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
      domain_id UUID NOT NULL REFERENCES custom_domains(id) ON DELETE CASCADE,
      record_type VARCHAR(10) NOT NULL,
      name VARCHAR(255) NOT NULL,
      value TEXT NOT NULL,
      ttl INTEGER DEFAULT 300,
      action VARCHAR(20) NOT NULL, -- create, update, delete
      created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
    );
    
    CREATE INDEX IF NOT EXISTS idx_dns_history_domain_id ON dns_history(domain_id);
    CREATE INDEX IF NOT EXISTS idx_dns_history_created_at ON dns_history(created_at);
    
    INSERT INTO schema_migrations (version, description) 
    VALUES ('006', 'Advanced domain management with registrars and DNS history');
    
  END IF;
END
$$;

-- ==========================================
-- MIGRATION 007: E-commerce features
-- ==========================================

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM schema_migrations WHERE version = '007'
  ) THEN
    
    -- Adicionar tipo de site e-commerce
    ALTER TYPE site_type ADD VALUE IF NOT EXISTS 'store';
    
    -- Criar tabela de produtos
    CREATE TABLE IF NOT EXISTS products (
      id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
      project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
      name VARCHAR(255) NOT NULL,
      description TEXT,
      price DECIMAL(10,2) NOT NULL,
      currency VARCHAR(3) DEFAULT 'USD',
      sku VARCHAR(100),
      inventory_quantity INTEGER DEFAULT 0,
      track_inventory BOOLEAN DEFAULT TRUE,
      weight DECIMAL(8,2),
      dimensions JSONB, -- {length, width, height, unit}
      images TEXT[] DEFAULT '{}',
      is_active BOOLEAN DEFAULT TRUE,
      seo_title VARCHAR(255),
      seo_description TEXT,
      metadata JSONB DEFAULT '{}',
      created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
      updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
    );
    
    CREATE TRIGGER update_products_updated_at 
      BEFORE UPDATE ON products 
      FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
    
    -- Criar tabela de variantes de produto
    CREATE TABLE IF NOT EXISTS product_variants (
      id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
      product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
      name VARCHAR(255) NOT NULL,
      price DECIMAL(10,2),
      sku VARCHAR(100),
      inventory_quantity INTEGER DEFAULT 0,
      attributes JSONB DEFAULT '{}', -- {color: "red", size: "L"}
      is_active BOOLEAN DEFAULT TRUE,
      created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
      updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
    );
    
    CREATE TRIGGER update_product_variants_updated_at 
      BEFORE UPDATE ON product_variants 
      FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
    
    -- Índices
    CREATE INDEX IF NOT EXISTS idx_products_project_id ON products(project_id);
    CREATE INDEX IF NOT EXISTS idx_products_sku ON products(sku);
    CREATE INDEX IF NOT EXISTS idx_products_is_active ON products(is_active);
    CREATE INDEX IF NOT EXISTS idx_product_variants_product_id ON product_variants(product_id);
    CREATE INDEX IF NOT EXISTS idx_product_variants_sku ON product_variants(sku);
    
    INSERT INTO schema_migrations (version, description) 
    VALUES ('007', 'E-commerce features with products and variants');
    
  END IF;
END
$$;

-- ==========================================
-- MIGRATION 008: Advanced caching
-- ==========================================

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM schema_migrations WHERE version = '008'
  ) THEN
    
    -- Adicionar configurações de cache aos projetos
    ALTER TABLE projects 
    ADD COLUMN IF NOT EXISTS cache_settings JSONB DEFAULT '{
      "enabled": true,
      "ttl": 3600,
      "strategies": ["browser", "cdn", "server"],
      "vary_headers": ["Accept-Encoding", "Accept-Language"]
    }';
    
    -- Adicionar configurações de CDN
    ALTER TABLE projects 
    ADD COLUMN IF NOT EXISTS cdn_settings JSONB DEFAULT '{
      "enabled": true,
      "provider": "cloudflare",
      "zones": {},
      "purge_on_deploy": true
    }';
    
    -- Criar tabela de cache invalidation
    CREATE TABLE IF NOT EXISTS cache_invalidations (
      id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
      project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
      cache_key VARCHAR(500) NOT NULL,
      invalidation_type VARCHAR(50) NOT NULL, -- path, tag, full
      status VARCHAR(20) DEFAULT 'pending', -- pending, processing, completed, failed
      provider VARCHAR(50) NOT NULL,
      provider_request_id VARCHAR(255),
      error_message TEXT,
      created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
      completed_at TIMESTAMPTZ
    );
    
    CREATE INDEX IF NOT EXISTS idx_cache_invalidations_project_id ON cache_invalidations(project_id);
    CREATE INDEX IF NOT EXISTS idx_cache_invalidations_status ON cache_invalidations(status);
    CREATE INDEX IF NOT EXISTS idx_cache_invalidations_created_at ON cache_invalidations(created_at);
    
    INSERT INTO schema_migrations (version, description) 
    VALUES ('008', 'Advanced caching and CDN features');
    
  END IF;
END
$$;

-- ==========================================
-- MIGRATION 009: API and webhooks
-- ==========================================

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM schema_migrations WHERE version = '009'
  ) THEN
    
    -- Criar tabela de API keys
    CREATE TABLE IF NOT EXISTS api_keys (
      id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
      user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
      name VARCHAR(255) NOT NULL,
      key_hash VARCHAR(255) UNIQUE NOT NULL,
      key_prefix VARCHAR(20) NOT NULL,
      scopes TEXT[] DEFAULT '{}',
      last_used_at TIMESTAMPTZ,
      expires_at TIMESTAMPTZ,
      is_active BOOLEAN DEFAULT TRUE,
      created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
    );
    
    -- Criar tabela de webhooks
    CREATE TABLE IF NOT EXISTS webhooks (
      id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
      user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
      project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
      url TEXT NOT NULL,
      events TEXT[] NOT NULL,
      secret VARCHAR(255) NOT NULL,
      is_active BOOLEAN DEFAULT TRUE,
      last_triggered_at TIMESTAMPTZ,
      created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
      updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
    );
    
    CREATE TRIGGER update_webhooks_updated_at 
      BEFORE UPDATE ON webhooks 
      FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
    
    -- Criar tabela de webhook deliveries
    CREATE TABLE IF NOT EXISTS webhook_deliveries (
      id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
      webhook_id UUID NOT NULL REFERENCES webhooks(id) ON DELETE CASCADE,
      event_type VARCHAR(100) NOT NULL,
      payload JSONB NOT NULL,
      response_status INTEGER,
      response_body TEXT,
      delivery_duration_ms INTEGER,
      attempted_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
      delivered_at TIMESTAMPTZ
    );
    
    -- Índices
    CREATE INDEX IF NOT EXISTS idx_api_keys_user_id ON api_keys(user_id);
    CREATE INDEX IF NOT EXISTS idx_api_keys_key_hash ON api_keys(key_hash);
    CREATE INDEX IF NOT EXISTS idx_api_keys_is_active ON api_keys(is_active);
    CREATE INDEX IF NOT EXISTS idx_webhooks_user_id ON webhooks(user_id);
    CREATE INDEX IF NOT EXISTS idx_webhooks_project_id ON webhooks(project_id);
    CREATE INDEX IF NOT EXISTS idx_webhook_deliveries_webhook_id ON webhook_deliveries(webhook_id);
    CREATE INDEX IF NOT EXISTS idx_webhook_deliveries_attempted_at ON webhook_deliveries(attempted_at);
    
    INSERT INTO schema_migrations (version, description) 
    VALUES ('009', 'API keys and webhooks system');
    
  END IF;
END
$$;

-- ==========================================
-- MIGRATION 010: Performance optimizations
-- ==========================================

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM schema_migrations WHERE version = '010'
  ) THEN
    
    -- Adicionar índices compostos para queries comuns
    CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_projects_user_status_updated 
    ON projects(user_id, status, updated_at DESC);
    
    CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_pages_project_published_order 
    ON pages(project_id, is_published, order_index);
    
    CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_ai_generations_user_status_created 
    ON ai_generations(user_id, status, created_at DESC);
    
    CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_usage_records_user_type_timestamp 
    ON usage_records(user_id, meter_type, timestamp DESC);
    
    -- Adicionar partial indices para dados ativos
    CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_projects_active_user 
    ON projects(user_id, updated_at DESC) 
    WHERE status IN ('draft', 'published');
    
    CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_subscriptions_active_user 
    ON subscriptions(user_id, current_period_end) 
    WHERE status = 'active';
    
    -- Criar função para limpeza automática de dados antigos
    CREATE OR REPLACE FUNCTION cleanup_old_records()
    RETURNS void AS $$
    BEGIN
      -- Deletar magic links expirados há mais de 7 dias
      DELETE FROM magic_links 
      WHERE expires_at < NOW() - INTERVAL '7 days';
      
      -- Deletar logs de auditoria com mais de 6 meses
      DELETE FROM audit_logs 
      WHERE created_at < NOW() - INTERVAL '6 months';
      
      -- Deletar notificações push antigas (mais de 60 dias)
      DELETE FROM push_notifications 
      WHERE created_at < NOW() - INTERVAL '60 days'
      AND clicked_at IS NULL;
      
      -- Deletar cache invalidations antigas (mais de 30 dias)
      DELETE FROM cache_invalidations 
      WHERE created_at < NOW() - INTERVAL '30 days';
      
      -- Deletar webhook deliveries antigas (mais de 90 dias)
      DELETE FROM webhook_deliveries 
      WHERE attempted_at < NOW() - INTERVAL '90 days';
      
    END;
    $$ LANGUAGE plpgsql;
    
    INSERT INTO schema_migrations (version, description) 
    VALUES ('010', 'Performance optimizations and cleanup procedures');
    
  END IF;
END
$$;

-- ==========================================
-- Verificar status das migrações
-- ==========================================

-- View para verificar status das migrações
CREATE OR REPLACE VIEW migration_status AS
SELECT 
  version,
  description,
  applied_at,
  AGE(NOW(), applied_at) AS time_since_applied
FROM schema_migrations
ORDER BY version;

-- Função para verificar se todas as migrações foram aplicadas
CREATE OR REPLACE FUNCTION check_migration_status()
RETURNS TABLE (
  total_migrations INTEGER,
  applied_migrations INTEGER,
  pending_migrations INTEGER,
  last_applied_version VARCHAR(50),
  last_applied_at TIMESTAMPTZ
) AS $$
BEGIN
  RETURN QUERY
  SELECT 
    10 as total_migrations,
    COUNT(*)::INTEGER as applied_migrations,
    (10 - COUNT(*))::INTEGER as pending_migrations,
    MAX(version) as last_applied_version,
    MAX(applied_at) as last_applied_at
  FROM schema_migrations;
END;
$$ LANGUAGE plpgsql;

-- Exibir status atual
SELECT * FROM check_migration_status();

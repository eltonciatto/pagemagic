-- Page Magic TimescaleDB Schema
-- Métricas e Analytics em Time Series

-- Extensões necessárias
CREATE EXTENSION IF NOT EXISTS timescaledb;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ==========================================
-- MÉTRICAS DE USO (STRIPE METERS)
-- ==========================================

-- Tabela principal de métricas
CREATE TABLE metrics (
  time TIMESTAMPTZ NOT NULL,
  user_id UUID NOT NULL,
  project_id UUID,
  meter_type VARCHAR(50) NOT NULL, -- page_generate, ai_token, container_hours, storage_gb
  quantity DECIMAL(10,4) NOT NULL,
  metadata JSONB DEFAULT '{}',
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Converter para hypertable (particionamento automático por tempo)
SELECT create_hypertable('metrics', 'time', 
  chunk_time_interval => INTERVAL '1 day',
  if_not_exists => TRUE
);

-- Índices para performance
CREATE INDEX idx_metrics_user_id_time ON metrics (user_id, time DESC);
CREATE INDEX idx_metrics_meter_type_time ON metrics (meter_type, time DESC);
CREATE INDEX idx_metrics_project_id_time ON metrics (project_id, time DESC);

-- ==========================================
-- ANALYTICS DE SITES
-- ==========================================

-- Tabela de page views
CREATE TABLE page_views (
  time TIMESTAMPTZ NOT NULL,
  project_id UUID NOT NULL,
  page_path VARCHAR(500) NOT NULL,
  visitor_id VARCHAR(100),
  session_id VARCHAR(100),
  ip_address INET,
  user_agent TEXT,
  referrer TEXT,
  country VARCHAR(2),
  city VARCHAR(100),
  device_type VARCHAR(20), -- desktop, mobile, tablet
  browser VARCHAR(50),
  os VARCHAR(50),
  duration_seconds INTEGER,
  metadata JSONB DEFAULT '{}'
);

SELECT create_hypertable('page_views', 'time',
  chunk_time_interval => INTERVAL '1 day',
  if_not_exists => TRUE
);

CREATE INDEX idx_page_views_project_id_time ON page_views (project_id, time DESC);
CREATE INDEX idx_page_views_visitor_id_time ON page_views (visitor_id, time DESC);
CREATE INDEX idx_page_views_session_id ON page_views (session_id);

-- ==========================================
-- EVENTOS PERSONALIZADOS
-- ==========================================

-- Tabela de eventos (cliques, conversões, etc.)
CREATE TABLE events (
  time TIMESTAMPTZ NOT NULL,
  project_id UUID NOT NULL,
  event_name VARCHAR(100) NOT NULL,
  visitor_id VARCHAR(100),
  session_id VARCHAR(100),
  page_path VARCHAR(500),
  event_data JSONB DEFAULT '{}',
  value DECIMAL(10,2),
  currency VARCHAR(3)
);

SELECT create_hypertable('events', 'time',
  chunk_time_interval => INTERVAL '1 day',
  if_not_exists => TRUE
);

CREATE INDEX idx_events_project_id_time ON events (project_id, time DESC);
CREATE INDEX idx_events_event_name_time ON events (event_name, time DESC);
CREATE INDEX idx_events_visitor_id_time ON events (visitor_id, time DESC);

-- ==========================================
-- PERFORMANCE E CORE WEB VITALS
-- ==========================================

-- Tabela de métricas de performance
CREATE TABLE performance_metrics (
  time TIMESTAMPTZ NOT NULL,
  project_id UUID NOT NULL,
  page_path VARCHAR(500) NOT NULL,
  visitor_id VARCHAR(100),
  -- Core Web Vitals
  lcp DECIMAL(8,3), -- Largest Contentful Paint (ms)
  fid DECIMAL(8,3), -- First Input Delay (ms)
  cls DECIMAL(8,6), -- Cumulative Layout Shift
  fcp DECIMAL(8,3), -- First Contentful Paint (ms)
  ttfb DECIMAL(8,3), -- Time to First Byte (ms)
  -- Outras métricas
  dom_load_time DECIMAL(8,3),
  window_load_time DECIMAL(8,3),
  connection_type VARCHAR(20),
  device_memory INTEGER,
  metadata JSONB DEFAULT '{}'
);

SELECT create_hypertable('performance_metrics', 'time',
  chunk_time_interval => INTERVAL '1 day',
  if_not_exists => TRUE
);

CREATE INDEX idx_performance_metrics_project_id_time ON performance_metrics (project_id, time DESC);

-- ==========================================
-- LOGS DE SISTEMA
-- ==========================================

-- Tabela de logs de aplicação
CREATE TABLE application_logs (
  time TIMESTAMPTZ NOT NULL,
  service VARCHAR(50) NOT NULL,
  level VARCHAR(10) NOT NULL, -- debug, info, warn, error, fatal
  message TEXT NOT NULL,
  error_stack TEXT,
  request_id VARCHAR(100),
  user_id UUID,
  metadata JSONB DEFAULT '{}'
);

SELECT create_hypertable('application_logs', 'time',
  chunk_time_interval => INTERVAL '1 hour',
  if_not_exists => TRUE
);

CREATE INDEX idx_application_logs_service_time ON application_logs (service, time DESC);
CREATE INDEX idx_application_logs_level_time ON application_logs (level, time DESC);
CREATE INDEX idx_application_logs_request_id ON application_logs (request_id);

-- ==========================================
-- MÉTRICAS DE INFRAESTRUTURA
-- ==========================================

-- Tabela de métricas de sistema
CREATE TABLE system_metrics (
  time TIMESTAMPTZ NOT NULL,
  service VARCHAR(50) NOT NULL,
  instance VARCHAR(100) NOT NULL,
  metric_name VARCHAR(100) NOT NULL,
  value DECIMAL(15,6) NOT NULL,
  unit VARCHAR(20),
  tags JSONB DEFAULT '{}'
);

SELECT create_hypertable('system_metrics', 'time',
  chunk_time_interval => INTERVAL '1 hour',
  if_not_exists => TRUE
);

CREATE INDEX idx_system_metrics_service_time ON system_metrics (service, time DESC);
CREATE INDEX idx_system_metrics_metric_name_time ON system_metrics (metric_name, time DESC);

-- ==========================================
-- VIEWS AGREGADAS
-- ==========================================

-- View de métricas por hora
CREATE VIEW hourly_metrics AS
SELECT 
  time_bucket('1 hour', time) AS hour,
  user_id,
  project_id,
  meter_type,
  SUM(quantity) AS total_quantity,
  COUNT(*) AS event_count
FROM metrics
GROUP BY hour, user_id, project_id, meter_type
ORDER BY hour DESC;

-- View de analytics diários
CREATE VIEW daily_analytics AS
SELECT 
  time_bucket('1 day', time) AS day,
  project_id,
  COUNT(DISTINCT visitor_id) AS unique_visitors,
  COUNT(*) AS total_views,
  AVG(duration_seconds) AS avg_duration,
  COUNT(DISTINCT session_id) AS sessions
FROM page_views
GROUP BY day, project_id
ORDER BY day DESC;

-- View de performance por página
CREATE VIEW page_performance_summary AS
SELECT 
  project_id,
  page_path,
  time_bucket('1 hour', time) AS hour,
  AVG(lcp) AS avg_lcp,
  AVG(fid) AS avg_fid,
  AVG(cls) AS avg_cls,
  AVG(fcp) AS avg_fcp,
  AVG(ttfb) AS avg_ttfb,
  COUNT(*) AS sample_count
FROM performance_metrics
WHERE time > NOW() - INTERVAL '24 hours'
GROUP BY project_id, page_path, hour
ORDER BY hour DESC;

-- View de top páginas
CREATE VIEW top_pages AS
SELECT 
  project_id,
  page_path,
  COUNT(*) AS views,
  COUNT(DISTINCT visitor_id) AS unique_visitors,
  AVG(duration_seconds) AS avg_duration
FROM page_views
WHERE time > NOW() - INTERVAL '7 days'
GROUP BY project_id, page_path
ORDER BY views DESC;

-- ==========================================
-- POLÍTICAS DE RETENÇÃO
-- ==========================================

-- Política de retenção para métricas (12 meses)
SELECT add_retention_policy('metrics', INTERVAL '12 months');

-- Política de retenção para page views (6 meses)
SELECT add_retention_policy('page_views', INTERVAL '6 months');

-- Política de retenção para eventos (12 meses)
SELECT add_retention_policy('events', INTERVAL '12 months');

-- Política de retenção para performance (3 meses)
SELECT add_retention_policy('performance_metrics', INTERVAL '3 months');

-- Política de retenção para logs (30 dias)
SELECT add_retention_policy('application_logs', INTERVAL '30 days');

-- Política de retenção para métricas de sistema (90 dias)
SELECT add_retention_policy('system_metrics', INTERVAL '90 days');

-- ==========================================
-- COMPRESSÃO
-- ==========================================

-- Habilitar compressão para chunks antigos (7 dias)
ALTER TABLE metrics SET (
  timescaledb.compress,
  timescaledb.compress_segmentby = 'user_id, meter_type',
  timescaledb.compress_orderby = 'time DESC'
);

SELECT add_compression_policy('metrics', INTERVAL '7 days');

ALTER TABLE page_views SET (
  timescaledb.compress,
  timescaledb.compress_segmentby = 'project_id',
  timescaledb.compress_orderby = 'time DESC'
);

SELECT add_compression_policy('page_views', INTERVAL '7 days');

ALTER TABLE events SET (
  timescaledb.compress,
  timescaledb.compress_segmentby = 'project_id, event_name',
  timescaledb.compress_orderby = 'time DESC'
);

SELECT add_compression_policy('events', INTERVAL '7 days');

-- ==========================================
-- FUNÇÕES AUXILIARES
-- ==========================================

-- Função para inserir métrica de uso
CREATE OR REPLACE FUNCTION insert_usage_metric(
  p_user_id UUID,
  p_project_id UUID,
  p_meter_type VARCHAR(50),
  p_quantity DECIMAL(10,4),
  p_metadata JSONB DEFAULT '{}'
) RETURNS UUID AS $$
DECLARE
  metric_id UUID;
BEGIN
  INSERT INTO metrics (time, user_id, project_id, meter_type, quantity, metadata)
  VALUES (NOW(), p_user_id, p_project_id, p_meter_type, p_quantity, p_metadata);
  
  RETURN metric_id;
END;
$$ LANGUAGE plpgsql;

-- Função para registrar page view
CREATE OR REPLACE FUNCTION track_page_view(
  p_project_id UUID,
  p_page_path VARCHAR(500),
  p_visitor_id VARCHAR(100),
  p_session_id VARCHAR(100),
  p_ip_address INET DEFAULT NULL,
  p_user_agent TEXT DEFAULT NULL,
  p_referrer TEXT DEFAULT NULL,
  p_metadata JSONB DEFAULT '{}'
) RETURNS void AS $$
BEGIN
  INSERT INTO page_views (
    time, project_id, page_path, visitor_id, session_id,
    ip_address, user_agent, referrer, metadata
  ) VALUES (
    NOW(), p_project_id, p_page_path, p_visitor_id, p_session_id,
    p_ip_address, p_user_agent, p_referrer, p_metadata
  );
END;
$$ LANGUAGE plpgsql;

-- Função para registrar evento personalizado
CREATE OR REPLACE FUNCTION track_event(
  p_project_id UUID,
  p_event_name VARCHAR(100),
  p_visitor_id VARCHAR(100),
  p_session_id VARCHAR(100),
  p_page_path VARCHAR(500) DEFAULT NULL,
  p_event_data JSONB DEFAULT '{}',
  p_value DECIMAL(10,2) DEFAULT NULL,
  p_currency VARCHAR(3) DEFAULT NULL
) RETURNS void AS $$
BEGIN
  INSERT INTO events (
    time, project_id, event_name, visitor_id, session_id,
    page_path, event_data, value, currency
  ) VALUES (
    NOW(), p_project_id, p_event_name, p_visitor_id, p_session_id,
    p_page_path, p_event_data, p_value, p_currency
  );
END;
$$ LANGUAGE plpgsql;

-- Função para obter estatísticas de uso do usuário
CREATE OR REPLACE FUNCTION get_user_usage_stats(
  p_user_id UUID,
  p_start_date TIMESTAMPTZ DEFAULT NOW() - INTERVAL '30 days',
  p_end_date TIMESTAMPTZ DEFAULT NOW()
) RETURNS TABLE (
  meter_type VARCHAR(50),
  total_quantity DECIMAL(15,4),
  daily_avg DECIMAL(15,4),
  peak_day TIMESTAMPTZ,
  peak_quantity DECIMAL(15,4)
) AS $$
BEGIN
  RETURN QUERY
  WITH daily_usage AS (
    SELECT 
      m.meter_type,
      time_bucket('1 day', m.time) as day,
      SUM(m.quantity) as day_quantity
    FROM metrics m
    WHERE m.user_id = p_user_id 
    AND m.time >= p_start_date 
    AND m.time <= p_end_date
    GROUP BY m.meter_type, day
  ),
  peak_usage AS (
    SELECT DISTINCT ON (meter_type)
      meter_type,
      day as peak_day,
      day_quantity as peak_quantity
    FROM daily_usage
    ORDER BY meter_type, day_quantity DESC
  )
  SELECT 
    du.meter_type,
    SUM(du.day_quantity) as total_quantity,
    AVG(du.day_quantity) as daily_avg,
    pu.peak_day,
    pu.peak_quantity
  FROM daily_usage du
  LEFT JOIN peak_usage pu ON du.meter_type = pu.meter_type
  GROUP BY du.meter_type, pu.peak_day, pu.peak_quantity;
END;
$$ LANGUAGE plpgsql;

-- ==========================================
-- JOBS DE AGREGAÇÃO CONTÍNUA
-- ==========================================

-- Agregação contínua para métricas horárias
CREATE MATERIALIZED VIEW hourly_usage_summary
WITH (timescaledb.continuous) AS
SELECT 
  time_bucket('1 hour', time) AS hour,
  user_id,
  meter_type,
  SUM(quantity) AS total_quantity,
  COUNT(*) AS event_count,
  AVG(quantity) AS avg_quantity,
  MAX(quantity) AS max_quantity
FROM metrics
GROUP BY hour, user_id, meter_type;

SELECT add_continuous_aggregate_policy('hourly_usage_summary',
  start_offset => INTERVAL '2 hours',
  end_offset => INTERVAL '1 hour',
  schedule_interval => INTERVAL '1 hour');

-- Agregação contínua para analytics diários
CREATE MATERIALIZED VIEW daily_project_analytics
WITH (timescaledb.continuous) AS
SELECT 
  time_bucket('1 day', time) AS day,
  project_id,
  COUNT(DISTINCT visitor_id) AS unique_visitors,
  COUNT(*) AS total_views,
  COUNT(DISTINCT session_id) AS sessions,
  AVG(duration_seconds) AS avg_duration,
  COUNT(DISTINCT page_path) AS unique_pages
FROM page_views
GROUP BY day, project_id;

SELECT add_continuous_aggregate_policy('daily_project_analytics',
  start_offset => INTERVAL '1 day',
  end_offset => INTERVAL '1 hour',
  schedule_interval => INTERVAL '1 hour');

-- ==========================================
-- ALERTAS E MONITORAMENTO
-- ==========================================

-- Função para detectar picos de uso anômalos
CREATE OR REPLACE FUNCTION detect_usage_anomalies(
  p_user_id UUID,
  p_threshold_multiplier DECIMAL DEFAULT 3.0
) RETURNS TABLE (
  meter_type VARCHAR(50),
  current_hour_usage DECIMAL(15,4),
  avg_baseline DECIMAL(15,4),
  anomaly_score DECIMAL(10,2)
) AS $$
BEGIN
  RETURN QUERY
  WITH current_hour AS (
    SELECT 
      m.meter_type,
      SUM(m.quantity) as hour_usage
    FROM metrics m
    WHERE m.user_id = p_user_id 
    AND m.time >= date_trunc('hour', NOW())
    GROUP BY m.meter_type
  ),
  baseline AS (
    SELECT 
      meter_type,
      AVG(total_quantity) as avg_usage,
      STDDEV(total_quantity) as stddev_usage
    FROM hourly_usage_summary
    WHERE user_id = p_user_id 
    AND hour >= NOW() - INTERVAL '7 days'
    AND hour < date_trunc('hour', NOW())
    GROUP BY meter_type
  )
  SELECT 
    ch.meter_type,
    ch.hour_usage,
    b.avg_usage,
    CASE 
      WHEN b.stddev_usage > 0 THEN 
        ABS(ch.hour_usage - b.avg_usage) / b.stddev_usage
      ELSE 0
    END as anomaly_score
  FROM current_hour ch
  JOIN baseline b ON ch.meter_type = b.meter_type
  WHERE ch.hour_usage > b.avg_usage * p_threshold_multiplier;
END;
$$ LANGUAGE plpgsql;

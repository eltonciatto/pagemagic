-- Page Magic Seed Data
-- Dados de desenvolvimento e teste

-- ==========================================
-- USUÁRIOS DE TESTE
-- ==========================================

-- Inserir usuários de desenvolvimento (senhas: "password123")
INSERT INTO users (id, email, email_verified, password_hash, first_name, last_name, status, locale) VALUES
  ('550e8400-e29b-41d4-a716-446655440001', 'admin@pagemagic.io', true, '$2b$10$rQUOzWcK6QTLw5lNfzgmWOZtG6TG5qQZw7ZwQwQwQwQwQwQwQwQwQ', 'Admin', 'User', 'active', 'en'),
  ('550e8400-e29b-41d4-a716-446655440002', 'user1@example.com', true, '$2b$10$rQUOzWcK6QTLw5lNfzgmWOZtG6TG5qQZw7ZwQwQwQwQwQwQwQwQwQ', 'John', 'Doe', 'active', 'en'),
  ('550e8400-e29b-41d4-a716-446655440003', 'user2@example.com', true, '$2b$10$rQUOzWcK6QTLw5lNfzgmWOZtG6TG5qQZw7ZwQwQwQwQwQwQwQwQwQ', 'Jane', 'Smith', 'active', 'pt'),
  ('550e8400-e29b-41d4-a716-446655440004', 'designer@example.com', true, '$2b$10$rQUOzWcK6QTLw5lNfzgmWOZtG6TG5qQZw7ZwQwQwQwQwQwQwQwQwQ', 'Sarah', 'Designer', 'active', 'en'),
  ('550e8400-e29b-41d4-a716-446655440005', 'developer@example.com', true, '$2b$10$rQUOzWcK6QTLw5lNfzgmWOZtG6TG5qQZw7ZwQwQwQwQwQwQwQwQwQ', 'Mike', 'Developer', 'active', 'en')
ON CONFLICT (id) DO NOTHING;

-- ==========================================
-- ASSINATURAS
-- ==========================================

-- Inserir assinaturas de teste
INSERT INTO subscriptions (id, user_id, plan, status, current_period_start, current_period_end) VALUES
  ('650e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440001', 'enterprise', 'active', NOW() - INTERVAL '15 days', NOW() + INTERVAL '15 days'),
  ('650e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440002', 'pro', 'active', NOW() - INTERVAL '10 days', NOW() + INTERVAL '20 days'),
  ('650e8400-e29b-41d4-a716-446655440003', '550e8400-e29b-41d4-a716-446655440003', 'basic', 'active', NOW() - INTERVAL '5 days', NOW() + INTERVAL '25 days'),
  ('650e8400-e29b-41d4-a716-446655440004', '550e8400-e29b-41d4-a716-446655440004', 'pro', 'active', NOW() - INTERVAL '20 days', NOW() + INTERVAL '10 days'),
  ('650e8400-e29b-41d4-a716-446655440005', '550e8400-e29b-41d4-a716-446655440005', 'free', 'active', NOW() - INTERVAL '30 days', NOW() + INTERVAL '30 days')
ON CONFLICT (id) DO NOTHING;

-- ==========================================
-- PROJETOS DE EXEMPLO
-- ==========================================

-- Inserir projetos de demonstração
INSERT INTO projects (id, user_id, name, description, type, status, subdomain, seo_title, seo_description, settings) VALUES
  (
    '750e8400-e29b-41d4-a716-446655440001', 
    '550e8400-e29b-41d4-a716-446655440002', 
    'Vegan Bakery Landing', 
    'Landing page for a local vegan bakery with online ordering', 
    'landing_page', 
    'published', 
    'vegan-bakery-demo',
    'Sweet Green Bakery - 100% Vegan Delights',
    'Delicious vegan pastries, cakes and breads made with love and natural ingredients',
    '{"theme": "organic", "colors": {"primary": "#4ade80", "secondary": "#22c55e"}}'
  ),
  (
    '750e8400-e29b-41d4-a716-446655440002', 
    '550e8400-e29b-41d4-a716-446655440002', 
    'Tech Startup Blog', 
    'Company blog for a fintech startup', 
    'blog', 
    'published', 
    'fintech-blog-demo',
    'FinTech Insights - Latest Trends in Financial Technology',
    'Expert insights, industry trends and thought leadership in financial technology',
    '{"theme": "modern", "colors": {"primary": "#3b82f6", "secondary": "#1e40af"}}'
  ),
  (
    '750e8400-e29b-41d4-a716-446655440003', 
    '550e8400-e29b-41d4-a716-446655440003', 
    'Designer Portfolio', 
    'Creative portfolio showcasing UX/UI design work', 
    'portfolio', 
    'published', 
    'sarah-design-portfolio',
    'Sarah Designer - UX/UI Portfolio',
    'Creative UX/UI designer specializing in mobile apps and web experiences',
    '{"theme": "minimal", "colors": {"primary": "#f59e0b", "secondary": "#d97706"}}'
  ),
  (
    '750e8400-e29b-41d4-a716-446655440004', 
    '550e8400-e29b-41d4-a716-446655440004', 
    'Local Restaurant', 
    'Restaurant website with menu and reservations', 
    'custom', 
    'published', 
    'bistro-bella-vita',
    'Bistro Bella Vita - Authentic Italian Cuisine',
    'Experience authentic Italian flavors in the heart of the city. Book your table today.',
    '{"theme": "elegant", "colors": {"primary": "#dc2626", "secondary": "#b91c1c"}}'
  ),
  (
    '750e8400-e29b-41d4-a716-446655440005', 
    '550e8400-e29b-41d4-a716-446655440005', 
    'E-commerce Store', 
    'Online store for handmade crafts', 
    'ecommerce', 
    'draft', 
    'handmade-crafts-store',
    'Artisan Crafts - Handmade with Love',
    'Unique handmade crafts, jewelry and home decor from local artisans',
    '{"theme": "rustic", "colors": {"primary": "#92400e", "secondary": "#78350f"}}'
  )
ON CONFLICT (id) DO NOTHING;

-- ==========================================
-- PÁGINAS DE EXEMPLO
-- ==========================================

-- Páginas para Vegan Bakery
INSERT INTO pages (id, project_id, name, slug, title, description, is_home, order_index) VALUES
  ('850e8400-e29b-41d4-a716-446655440001', '750e8400-e29b-41d4-a716-446655440001', 'Home', '', 'Sweet Green Bakery', 'Welcome to our vegan paradise', true, 0),
  ('850e8400-e29b-41d4-a716-446655440002', '750e8400-e29b-41d4-a716-446655440001', 'Menu', 'menu', 'Our Vegan Menu', 'Explore our delicious vegan offerings', false, 1),
  ('850e8400-e29b-41d4-a716-446655440003', '750e8400-e29b-41d4-a716-446655440001', 'About', 'about', 'Our Story', 'Learn about our vegan journey', false, 2),
  ('850e8400-e29b-41d4-a716-446655440004', '750e8400-e29b-41d4-a716-446655440001', 'Contact', 'contact', 'Get in Touch', 'Visit us or place an order', false, 3)
ON CONFLICT (id) DO NOTHING;

-- Páginas para Tech Blog
INSERT INTO pages (id, project_id, name, slug, title, description, is_home, order_index) VALUES
  ('850e8400-e29b-41d4-a716-446655440005', '750e8400-e29b-41d4-a716-446655440002', 'Home', '', 'FinTech Insights', 'Latest trends in financial technology', true, 0),
  ('850e8400-e29b-41d4-a716-446655440006', '750e8400-e29b-41d4-a716-446655440002', 'Articles', 'articles', 'All Articles', 'Browse our latest articles', false, 1),
  ('850e8400-e29b-41d4-a716-446655440007', '750e8400-e29b-41d4-a716-446655440002', 'Authors', 'authors', 'Our Authors', 'Meet our expert writers', false, 2),
  ('850e8400-e29b-41d4-a716-446655440008', '750e8400-e29b-41d4-a716-446655440002', 'Newsletter', 'newsletter', 'Newsletter', 'Subscribe to our weekly insights', false, 3)
ON CONFLICT (id) DO NOTHING;

-- Páginas para Portfolio
INSERT INTO pages (id, project_id, name, slug, title, description, is_home, order_index) VALUES
  ('850e8400-e29b-41d4-a716-446655440009', '750e8400-e29b-41d4-a716-446655440003', 'Home', '', 'Sarah Designer', 'UX/UI Designer Portfolio', true, 0),
  ('850e8400-e29b-41d4-a716-446655440010', '750e8400-e29b-41d4-a716-446655440003', 'Work', 'work', 'My Work', 'Selected projects and case studies', false, 1),
  ('850e8400-e29b-41d4-a716-446655440011', '750e8400-e29b-41d4-a716-446655440003', 'About', 'about', 'About Me', 'My design philosophy and background', false, 2),
  ('850e8400-e29b-41d4-a716-446655440012', '750e8400-e29b-41d4-a716-446655440003', 'Contact', 'contact', 'Lets Work Together', 'Get in touch for collaborations', false, 3)
ON CONFLICT (id) DO NOTHING;

-- ==========================================
-- SEÇÕES DE EXEMPLO
-- ==========================================

-- Seções para Vegan Bakery Home
INSERT INTO sections (id, page_id, type, name, order_index, html_content, css_styles, data) VALUES
  (
    '950e8400-e29b-41d4-a716-446655440001',
    '850e8400-e29b-41d4-a716-446655440001',
    'hero',
    'Hero Section',
    0,
    '<section class="hero"><div class="container"><h1>{{hero.title}}</h1><p>{{hero.subtitle}}</p><button class="btn-primary">{{hero.cta}}</button></div></section>',
    '.hero { background: linear-gradient(135deg, #4ade80 0%, #22c55e 100%); padding: 4rem 0; text-align: center; } .hero h1 { font-size: 3rem; color: white; margin-bottom: 1rem; } .hero p { font-size: 1.2rem; color: rgba(255,255,255,0.9); margin-bottom: 2rem; }',
    '{"hero": {"title": "Welcome to Sweet Green Bakery", "subtitle": "100% Vegan Delights Made with Love", "cta": "Order Now"}}'
  ),
  (
    '950e8400-e29b-41d4-a716-446655440002',
    '850e8400-e29b-41d4-a716-446655440001',
    'about',
    'About Section',
    1,
    '<section class="about"><div class="container"><div class="row"><div class="col-md-6"><h2>{{about.title}}</h2><p>{{about.description}}</p></div><div class="col-md-6"><img src="{{about.image}}" alt="Bakery"></div></div></div></section>',
    '.about { padding: 4rem 0; } .about h2 { font-size: 2.5rem; margin-bottom: 1.5rem; color: #22c55e; } .about p { font-size: 1.1rem; line-height: 1.6; }',
    '{"about": {"title": "Our Vegan Philosophy", "description": "We believe that delicious food should be kind to animals and the planet. Every item in our bakery is 100% plant-based, made with organic ingredients and lots of love.", "image": "/images/bakery-interior.jpg"}}'
  ),
  (
    '950e8400-e29b-41d4-a716-446655440003',
    '850e8400-e29b-41d4-a716-446655440001',
    'services',
    'Featured Products',
    2,
    '<section class="featured-products"><div class="container"><h2>{{products.title}}</h2><div class="product-grid">{{#each products.items}}<div class="product-card"><img src="{{image}}" alt="{{name}}"><h3>{{name}}</h3><p>{{description}}</p><span class="price">${{price}}</span></div>{{/each}}</div></div></section>',
    '.featured-products { padding: 4rem 0; background: #f8fafc; } .product-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 2rem; margin-top: 2rem; } .product-card { background: white; border-radius: 8px; padding: 1.5rem; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }',
    '{"products": {"title": "Featured Products", "items": [{"name": "Chocolate Chip Cookies", "description": "Classic cookies with vegan chocolate chips", "price": "12.99", "image": "/images/cookies.jpg"}, {"name": "Sourdough Bread", "description": "Artisan sourdough made with ancient grains", "price": "8.99", "image": "/images/bread.jpg"}, {"name": "Carrot Cake", "description": "Moist carrot cake with cashew cream frosting", "price": "24.99", "image": "/images/cake.jpg"}]}}'
  )
ON CONFLICT (id) DO NOTHING;

-- ==========================================
-- GERAÇÕES DE IA DE EXEMPLO
-- ==========================================

INSERT INTO ai_generations (id, user_id, project_id, prompt, model, status, response_data, tokens_used, processing_time_ms) VALUES
  (
    'a50e8400-e29b-41d4-a716-446655440001',
    '550e8400-e29b-41d4-a716-446655440002',
    '750e8400-e29b-41d4-a716-446655440001',
    'Create a landing page for a vegan bakery that emphasizes sustainability and health',
    'llama-3-70b',
    'completed',
    '{"sections": [{"type": "hero", "content": "Welcome to Sweet Green Bakery"}, {"type": "about", "content": "Our story of sustainable baking"}], "css": "/* Generated styles */", "recommendations": ["Add customer testimonials", "Include nutritional information"]}',
    1250,
    5400
  ),
  (
    'a50e8400-e29b-41d4-a716-446655440002',
    '550e8400-e29b-41d4-a716-446655440003',
    '750e8400-e29b-41d4-a716-446655440003',
    'Design a modern portfolio website for a UX/UI designer focusing on mobile app design',
    'gpt-4',
    'completed',
    '{"sections": [{"type": "hero", "content": "Sarah Designer Portfolio"}, {"type": "portfolio", "content": "Featured work showcase"}], "css": "/* Modern portfolio styles */", "recommendations": ["Add case study details", "Include client testimonials"]}',
    890,
    3200
  ),
  (
    'a50e8400-e29b-41d4-a716-446655440003',
    '550e8400-e29b-41d4-a716-446655440004',
    '750e8400-e29b-41d4-a716-446655440004',
    'Create a restaurant website with elegant design and online reservation system',
    'claude-3-sonnet',
    'completed',
    '{"sections": [{"type": "hero", "content": "Bistro Bella Vita"}, {"type": "menu", "content": "Our Italian menu"}], "css": "/* Elegant restaurant styles */", "recommendations": ["Add photo gallery", "Include chef information"]}',
    1050,
    4100
  )
ON CONFLICT (id) DO NOTHING;

-- ==========================================
-- TOKENS PUSH MOBILE DE EXEMPLO
-- ==========================================

INSERT INTO mobile_push_tokens (id, user_id, token, platform, app_version, device_id) VALUES
  ('b50e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440002', 'ExponentPushToken[xxxxxxxxxxxxxxxxxxxxxx]', 'expo', '1.0.0', 'device-001'),
  ('b50e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440003', 'ExponentPushToken[yyyyyyyyyyyyyyyyyyyyyy]', 'expo', '1.0.0', 'device-002'),
  ('b50e8400-e29b-41d4-a716-446655440003', '550e8400-e29b-41d4-a716-446655440004', 'apns-token-zzzzzzzzzzzzzzzzzzzzzzzz', 'ios', '1.0.0', 'device-003'),
  ('b50e8400-e29b-41d4-a716-446655440004', '550e8400-e29b-41d4-a716-446655440005', 'fcm-token-aaaaaaaaaaaaaaaaaaaaaaa', 'android', '1.0.0', 'device-004')
ON CONFLICT (id) DO NOTHING;

-- ==========================================
-- TRADUÇÕES DE EXEMPLO
-- ==========================================

-- Traduções para o projeto Vegan Bakery
INSERT INTO translations (key, language_code, value, project_id) VALUES
  ('hero.title', 'en', 'Welcome to Sweet Green Bakery', '750e8400-e29b-41d4-a716-446655440001'),
  ('hero.title', 'pt', 'Bem-vindos à Padaria Verde Doce', '750e8400-e29b-41d4-a716-446655440001'),
  ('hero.title', 'es', 'Bienvenidos a Panadería Verde Dulce', '750e8400-e29b-41d4-a716-446655440001'),
  ('hero.subtitle', 'en', '100% Vegan Delights Made with Love', '750e8400-e29b-41d4-a716-446655440001'),
  ('hero.subtitle', 'pt', '100% Delícias Veganas Feitas com Amor', '750e8400-e29b-41d4-a716-446655440001'),
  ('hero.subtitle', 'es', '100% Delicias Veganas Hechas con Amor', '750e8400-e29b-41d4-a716-446655440001'),
  ('hero.cta', 'en', 'Order Now', '750e8400-e29b-41d4-a716-446655440001'),
  ('hero.cta', 'pt', 'Peça Agora', '750e8400-e29b-41d4-a716-446655440001'),
  ('hero.cta', 'es', 'Ordenar Ahora', '750e8400-e29b-41d4-a716-446655440001'),
  ('about.title', 'en', 'Our Vegan Philosophy', '750e8400-e29b-41d4-a716-446655440001'),
  ('about.title', 'pt', 'Nossa Filosofia Vegana', '750e8400-e29b-41d4-a716-446655440001'),
  ('about.title', 'es', 'Nuestra Filosofía Vegana', '750e8400-e29b-41d4-a716-446655440001')
ON CONFLICT (key, language_code, project_id) DO NOTHING;

-- ==========================================
-- TEMPLATES DE PROMPT
-- ==========================================

INSERT INTO prompt_templates (id, name, description, category, template, variables, is_public, created_by, usage_count, rating) VALUES
  (
    'c50e8400-e29b-41d4-a716-446655440001',
    'Vegan Restaurant Landing Page',
    'Template for creating landing pages for vegan restaurants and food businesses',
    'food_service',
    'Create a landing page for a {{business_type}} called "{{business_name}}" that specializes in {{cuisine_type}}. The page should emphasize {{key_values}} and include sections for {{sections}}. Use a {{design_style}} design style with {{color_scheme}} colors.',
    '["business_type", "business_name", "cuisine_type", "key_values", "sections", "design_style", "color_scheme"]',
    true,
    '550e8400-e29b-41d4-a716-446655440001',
    45,
    4.2
  ),
  (
    'c50e8400-e29b-41d4-a716-446655440002',
    'Portfolio Website',
    'Template for creating professional portfolio websites',
    'portfolio',
    'Design a portfolio website for a {{profession}} named {{name}} who specializes in {{specialization}}. Include sections for {{portfolio_sections}} and emphasize {{key_skills}}. Use a {{design_approach}} approach with {{visual_style}} styling.',
    '["profession", "name", "specialization", "portfolio_sections", "key_skills", "design_approach", "visual_style"]',
    true,
    '550e8400-e29b-41d4-a716-446655440001',
    32,
    4.5
  ),
  (
    'c50e8400-e29b-41d4-a716-446655440003',
    'Tech Startup Landing',
    'Template for tech startup landing pages',
    'technology',
    'Create a landing page for a {{startup_type}} startup called "{{company_name}}" that offers {{product_service}}. Target audience is {{target_audience}}. Include sections for {{features}} and emphasize {{value_proposition}}. Use a {{tech_style}} design.',
    '["startup_type", "company_name", "product_service", "target_audience", "features", "value_proposition", "tech_style"]',
    true,
    '550e8400-e29b-41d4-a716-446655440001',
    67,
    4.7
  ),
  (
    'c50e8400-e29b-41d4-a716-446655440004',
    'E-commerce Store',
    'Template for creating online store websites',
    'ecommerce',
    'Build an e-commerce website for {{store_name}} selling {{product_category}}. Target market is {{target_market}}. Include {{ecommerce_features}} and emphasize {{unique_selling_points}}. Use {{store_style}} design with {{brand_personality}} personality.',
    '["store_name", "product_category", "target_market", "ecommerce_features", "unique_selling_points", "store_style", "brand_personality"]',
    true,
    '550e8400-e29b-41d4-a716-446655440001',
    28,
    4.3
  ),
  (
    'c50e8400-e29b-41d4-a716-446655440005',
    'Blog Website',
    'Template for creating blog and content websites',
    'content',
    'Create a blog website for {{blog_name}} focusing on {{blog_topic}}. Target audience is {{reader_demographic}}. Include sections for {{blog_sections}} and features like {{blog_features}}. Use a {{content_style}} layout optimized for {{content_type}}.',
    '["blog_name", "blog_topic", "reader_demographic", "blog_sections", "blog_features", "content_style", "content_type"]',
    true,
    '550e8400-e29b-41d4-a716-446655440001',
    53,
    4.4
  )
ON CONFLICT (id) DO NOTHING;

-- ==========================================
-- DADOS DE USO (MÉTRICAS)
-- ==========================================

-- Inserir alguns registros de uso dos últimos 30 dias
INSERT INTO usage_records (user_id, subscription_id, meter_type, quantity, timestamp) 
SELECT 
  u.id,
  s.id,
  meter_type,
  (RANDOM() * 10)::DECIMAL(10,2),
  NOW() - (RANDOM() * INTERVAL '30 days')
FROM users u
JOIN subscriptions s ON u.id = s.user_id
CROSS JOIN (
  VALUES 
    ('page_generate'),
    ('ai_token'),
    ('container_hours'),
    ('storage_gb')
) AS meters(meter_type)
WHERE u.status = 'active'
AND RANDOM() < 0.3; -- 30% chance para cada combinação

-- ==========================================
-- LOGS DE AUDITORIA
-- ==========================================

-- Inserir alguns logs de auditoria
INSERT INTO audit_logs (user_id, action, resource_type, resource_id, new_values, ip_address, user_agent) VALUES
  ('550e8400-e29b-41d4-a716-446655440002', 'create', 'project', '750e8400-e29b-41d4-a716-446655440001', '{"name": "Vegan Bakery Landing"}', '192.168.1.1', 'Mozilla/5.0...'),
  ('550e8400-e29b-41d4-a716-446655440002', 'publish', 'project', '750e8400-e29b-41d4-a716-446655440001', '{"status": "published"}', '192.168.1.1', 'Mozilla/5.0...'),
  ('550e8400-e29b-41d4-a716-446655440003', 'create', 'project', '750e8400-e29b-41d4-a716-446655440003', '{"name": "Designer Portfolio"}', '192.168.1.2', 'Mozilla/5.0...'),
  ('550e8400-e29b-41d4-a716-446655440004', 'create', 'project', '750e8400-e29b-41d4-a716-446655440004', '{"name": "Local Restaurant"}', '192.168.1.3', 'Mozilla/5.0...')
ON CONFLICT DO NOTHING;

-- ==========================================
-- WEBHOOKS DE EXEMPLO
-- ==========================================

INSERT INTO webhooks (id, user_id, project_id, url, events, secret) VALUES
  ('d50e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440002', '750e8400-e29b-41d4-a716-446655440001', 'https://api.example.com/webhooks/pagemagic', ARRAY['site.published', 'site.updated'], 'wh_secret_123'),
  ('d50e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440003', '750e8400-e29b-41d4-a716-446655440003', 'https://designer-tools.com/webhook', ARRAY['site.published'], 'wh_secret_456')
ON CONFLICT (id) DO NOTHING;

-- ==========================================
-- API KEYS DE EXEMPLO
-- ==========================================

INSERT INTO api_keys (id, user_id, name, key_hash, key_prefix, scopes) VALUES
  ('e50e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440002', 'Production API', '$2b$10$hashedapikey...', 'pk_live_', ARRAY['projects:read', 'projects:write', 'analytics:read']),
  ('e50e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440003', 'Dev Environment', '$2b$10$hashedapikey...', 'pk_test_', ARRAY['projects:read', 'analytics:read']),
  ('e50e8400-e29b-41d4-a716-446655440003', '550e8400-e29b-41d4-a716-446655440004', 'Analytics Only', '$2b$10$hashedapikey...', 'pk_analytics_', ARRAY['analytics:read'])
ON CONFLICT (id) DO NOTHING;

-- ==========================================
-- FINALIZAÇÃO
-- ==========================================

-- Atualizar estatísticas das tabelas
ANALYZE users;
ANALYZE projects;
ANALYZE pages;
ANALYZE sections;
ANALYZE ai_generations;
ANALYZE subscriptions;
ANALYZE usage_records;

-- Exibir resumo dos dados inseridos
SELECT 
  'Users' as table_name, COUNT(*) as records FROM users
UNION ALL
SELECT 'Projects', COUNT(*) FROM projects
UNION ALL
SELECT 'Pages', COUNT(*) FROM pages
UNION ALL
SELECT 'Sections', COUNT(*) FROM sections
UNION ALL
SELECT 'AI Generations', COUNT(*) FROM ai_generations
UNION ALL
SELECT 'Subscriptions', COUNT(*) FROM subscriptions
UNION ALL
SELECT 'Usage Records', COUNT(*) FROM usage_records
UNION ALL
SELECT 'Push Tokens', COUNT(*) FROM mobile_push_tokens
UNION ALL
SELECT 'Translations', COUNT(*) FROM translations
UNION ALL
SELECT 'Prompt Templates', COUNT(*) FROM prompt_templates
UNION ALL
SELECT 'Webhooks', COUNT(*) FROM webhooks
UNION ALL
SELECT 'API Keys', COUNT(*) FROM api_keys;

-- Log de conclusão
INSERT INTO audit_logs (action, resource_type, new_values, ip_address) 
VALUES ('create', 'seed_data', '{"message": "Development seed data inserted successfully"}', '127.0.0.1');

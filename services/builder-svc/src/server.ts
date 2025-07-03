import express from 'express';
import cors from 'cors';
import helmet from 'helmet';
import compression from 'compression';
import morgan from 'morgan';
import { register } from 'prom-client';
import { WebSocketServer } from 'ws';
import { createServer } from 'http';

import { ASTService } from './services/ASTService';
import { ReactGenerator } from './generators/ReactGenerator';
import { VueGenerator } from './generators/VueGenerator';
import { AngularGenerator } from './generators/AngularGenerator';
import { ThemeService } from './services/ThemeService';
import { OptimizationService } from './services/OptimizationService';
import { PreviewService } from './services/PreviewService';
import { ComponentLibrary } from './services/ComponentLibrary';
import { buildRoutes } from './routes/build';
import { astRoutes } from './routes/ast';
import { previewRoutes } from './routes/preview';
import { themeRoutes } from './routes/theme';
import { logger } from './utils/logger';
import { errorHandler } from './middleware/errorHandler';
import { requestLogger } from './middleware/requestLogger';
import { metricsMiddleware } from './middleware/metrics';

// Initialize services
const astService = new ASTService();
const reactGenerator = new ReactGenerator();
const vueGenerator = new VueGenerator();
const angularGenerator = new AngularGenerator();
const themeService = new ThemeService();
const optimizationService = new OptimizationService();
const previewService = new PreviewService();
const componentLibrary = new ComponentLibrary();

const app = express();
const server = createServer(app);

// WebSocket for real-time collaboration
const wss = new WebSocketServer({ server });

wss.on('connection', (ws, req) => {
  logger.info('New WebSocket connection established');
  
  ws.on('message', async (data) => {
    try {
      const message = JSON.parse(data.toString());
      
      switch (message.type) {
        case 'ast_update':
          // Broadcast AST updates to other clients
          const updatedAST = await astService.updateAST(message.payload.id, message.payload.changes);
          wss.clients.forEach(client => {
            if (client !== ws && client.readyState === ws.OPEN) {
              client.send(JSON.stringify({
                type: 'ast_updated',
                payload: { id: message.payload.id, ast: updatedAST }
              }));
            }
          });
          break;
          
        case 'preview_request':
          // Generate live preview
          const preview = await previewService.generatePreview(message.payload.ast);
          ws.send(JSON.stringify({
            type: 'preview_ready',
            payload: { id: message.payload.id, preview }
          }));
          break;
      }
    } catch (error) {
      logger.error('WebSocket error:', error);
      ws.send(JSON.stringify({
        type: 'error',
        payload: { message: 'Failed to process message' }
      }));
    }
  });
  
  ws.on('close', () => {
    logger.info('WebSocket connection closed');
  });
});

// Middleware
app.use(helmet({
  contentSecurityPolicy: {
    directives: {
      defaultSrc: ["'self'"],
      styleSrc: ["'self'", "'unsafe-inline'", "https://fonts.googleapis.com"],
      fontSrc: ["'self'", "https://fonts.gstatic.com"],
      imgSrc: ["'self'", "data:", "https:"],
      scriptSrc: ["'self'"],
      connectSrc: ["'self'", "ws:", "wss:"]
    }
  }
}));

app.use(compression());
app.use(cors({
  origin: process.env.CORS_ORIGINS?.split(',') || ['http://localhost:3000'],
  credentials: true
}));

app.use(morgan('combined', { stream: { write: message => logger.info(message.trim()) } }));
app.use(express.json({ limit: '50mb' }));
app.use(express.urlencoded({ extended: true, limit: '50mb' }));

// Custom middleware
app.use(requestLogger);
app.use(metricsMiddleware);

// Health check
app.get('/health', (req, res) => {
  res.json({
    status: 'healthy',
    service: 'builder-svc',
    version: '2.0.0',
    timestamp: new Date().toISOString(),
    uptime: process.uptime(),
    memory: process.memoryUsage()
  });
});

// Metrics endpoint
app.get('/metrics', async (req, res) => {
  res.set('Content-Type', register.contentType);
  res.end(await register.metrics());
});

// API Routes
app.use('/api/v1/build', buildRoutes);
app.use('/api/v1/ast', astRoutes);
app.use('/api/v1/preview', previewRoutes);
app.use('/api/v1/themes', themeRoutes);

// Component library endpoint
app.get('/api/v1/components', async (req, res) => {
  try {
    const category = req.query.category as string;
    const framework = req.query.framework as string || 'react';
    
    const components = await componentLibrary.getComponents({
      category,
      framework: framework as 'react' | 'vue' | 'angular'
    });
    
    res.json({
      success: true,
      data: components,
      total: components.length
    });
  } catch (error) {
    logger.error('Failed to get components:', error);
    res.status(500).json({
      success: false,
      error: 'Failed to get components'
    });
  }
});

// Advanced build endpoint with multiple frameworks
app.post('/api/v1/build/advanced', async (req, res) => {
  try {
    const { 
      site_data, 
      framework = 'react', 
      theme_id, 
      optimization_level = 'standard',
      target_devices = ['desktop', 'mobile'],
      performance_budget = {},
      accessibility_level = 'aa'
    } = req.body;

    if (!site_data) {
      return res.status(400).json({
        success: false,
        error: 'site_data is required'
      });
    }

    logger.info(`Starting advanced build for framework: ${framework}`);

    // Convert to AST
    const ast = await astService.convertToAST(site_data, {
      framework,
      target_devices,
      accessibility_level
    });

    // Apply theme
    let themedAST = ast;
    if (theme_id) {
      themedAST = await themeService.applyTheme(ast, theme_id);
    }

    // Generate code based on framework
    let generatedCode;
    switch (framework) {
      case 'react':
        generatedCode = await reactGenerator.generate(themedAST, {
          typescript: true,
          hooks: true,
          styled_components: true
        });
        break;
      case 'vue':
        generatedCode = await vueGenerator.generate(themedAST, {
          composition_api: true,
          typescript: true,
          scoped_styles: true
        });
        break;
      case 'angular':
        generatedCode = await angularGenerator.generate(themedAST, {
          standalone_components: true,
          signals: true,
          typescript: true
        });
        break;
      default:
        throw new Error(`Unsupported framework: ${framework}`);
    }

    // Apply optimizations
    const optimizedCode = await optimizationService.optimize(generatedCode, {
      level: optimization_level,
      target_devices,
      performance_budget
    });

    // Generate preview
    const preview = await previewService.generatePreview(themedAST, {
      framework,
      devices: target_devices
    });

    res.json({
      success: true,
      data: {
        ast: themedAST,
        code: optimizedCode,
        preview: preview,
        metadata: {
          framework,
          theme_id,
          optimization_level,
          target_devices,
          build_time: Date.now() - req.startTime,
          size_estimate: optimizedCode.estimated_size,
          performance_score: optimizedCode.performance_score
        }
      }
    });

  } catch (error) {
    logger.error('Advanced build failed:', error);
    res.status(500).json({
      success: false,
      error: error instanceof Error ? error.message : 'Advanced build failed'
    });
  }
});

// Error handling
app.use(errorHandler);

// 404 handler
app.use('*', (req, res) => {
  res.status(404).json({
    success: false,
    error: 'Endpoint not found'
  });
});

const port = process.env.PORT || 8080;

server.listen(port, () => {
  logger.info(`🚀 Builder service (Enterprise) running on port ${port}`);
  logger.info(`📡 WebSocket server ready for real-time collaboration`);
  logger.info(`🎨 Supported frameworks: React, Vue, Angular`);
  logger.info(`🔧 Advanced optimizations enabled`);
});

// Graceful shutdown
process.on('SIGTERM', () => {
  logger.info('SIGTERM received, shutting down gracefully');
  server.close(() => {
    process.exit(0);
  });
});

process.on('SIGINT', () => {
  logger.info('SIGINT received, shutting down gracefully');
  server.close(() => {
    process.exit(0);
  });
});

export default app;

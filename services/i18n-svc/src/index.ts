import express from 'express';
import cors from 'cors';
import helmet from 'helmet';
import compression from 'compression';
import morgan from 'morgan';
import { Pool } from 'pg';
import { createClient } from 'redis';
import dotenv from 'dotenv';

import { TranslationService } from './services/translation';
import { TranslationController } from './controllers/translation';
import { translationRoutes } from './routes/translation';
import { logger } from './utils/logger';

// Load environment variables
dotenv.config();

class I18nService {
  private app: express.Application;
  private db: Pool;
  private redis: any;
  private translationService: TranslationService;
  private translationController: TranslationController;

  constructor() {
    this.app = express();
    this.initializeDatabase();
    this.initializeRedis();
    this.initializeMiddleware();
    this.initializeServices();
    this.initializeRoutes();
    this.initializeErrorHandling();
  }

  private initializeDatabase(): void {
    this.db = new Pool({
      connectionString: process.env.DATABASE_URL || 'postgresql://user:pass@localhost:5432/pagemagic',
      max: 20,
      idleTimeoutMillis: 30000,
      connectionTimeoutMillis: 2000,
    });

    this.db.on('error', (err) => {
      logger.error('Unexpected error on idle client', err);
    });
  }

  private async initializeRedis(): Promise<void> {
    this.redis = createClient({
      url: process.env.REDIS_URL || 'redis://localhost:6379'
    });

    this.redis.on('error', (err: Error) => {
      logger.error('Redis client error', err);
    });

    try {
      await this.redis.connect();
      logger.info('Connected to Redis');
    } catch (error) {
      logger.error('Failed to connect to Redis', error);
    }
  }

  private initializeMiddleware(): void {
    // Security middleware
    this.app.use(helmet());
    
    // CORS
    this.app.use(cors({
      origin: process.env.ALLOWED_ORIGINS?.split(',') || '*',
      credentials: true
    }));
    
    // Compression
    this.app.use(compression());
    
    // Request logging
    this.app.use(morgan('combined', {
      stream: {
        write: (message: string) => {
          logger.info(message.trim());
        }
      }
    }));
    
    // Body parsing
    this.app.use(express.json({ limit: '10mb' }));
    this.app.use(express.urlencoded({ extended: true, limit: '10mb' }));
  }

  private initializeServices(): void {
    this.translationService = new TranslationService(this.db, this.redis);
    this.translationController = new TranslationController(this.translationService);
  }

  private initializeRoutes(): void {
    // Health check
    this.app.get('/health', this.translationController.healthCheck.bind(this.translationController));
    
    // API routes
    this.app.use('/v1', translationRoutes(this.translationController));
    
    // 404 handler
    this.app.use('*', (req, res) => {
      res.status(404).json({
        error: 'Route not found',
        path: req.originalUrl
      });
    });
  }

  private initializeErrorHandling(): void {
    this.app.use((error: Error, req: express.Request, res: express.Response, next: express.NextFunction) => {
      logger.error('Unhandled error', {
        error: error.message,
        stack: error.stack,
        path: req.path,
        method: req.method
      });

      res.status(500).json({
        error: 'Internal server error',
        message: process.env.NODE_ENV === 'development' ? error.message : 'Something went wrong'
      });
    });
  }

  public async start(): Promise<void> {
    const port = process.env.PORT || 8080;
    
    try {
      // Test database connection
      await this.db.query('SELECT 1');
      logger.info('Database connection established');
      
      // Start server
      this.app.listen(port, () => {
        logger.info(`I18n service started on port ${port}`);
      });
    } catch (error) {
      logger.error('Failed to start service', error);
      process.exit(1);
    }
  }

  public async shutdown(): Promise<void> {
    logger.info('Shutting down i18n service...');
    
    try {
      await this.db.end();
      await this.redis.quit();
      logger.info('Service shutdown complete');
    } catch (error) {
      logger.error('Error during shutdown', error);
    }
  }
}

// Start the service
const service = new I18nService();

// Graceful shutdown
process.on('SIGTERM', () => {
  service.shutdown();
});

process.on('SIGINT', () => {
  service.shutdown();
});

// Start the service
service.start().catch((error) => {
  logger.error('Failed to start service', error);
  process.exit(1);
});

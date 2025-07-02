import express from 'express';
import cors from 'cors';
import helmet from 'helmet';
import compression from 'compression';
import { generateRoutes } from './routes/generate';

// Simplified config for MVP
const config = {
  port: process.env.PORT || 3001,
  cors: {
    origins: process.env.CORS_ORIGINS?.split(',') || ['http://localhost:3000']
  }
};

const app = express();

// Basic middleware
app.use(helmet());
app.use(compression());
app.use(cors({
  origin: config.cors.origins,
  credentials: true,
}));
app.use(express.json({ limit: '10mb' }));
app.use(express.urlencoded({ extended: true, limit: '10mb' }));

// Routes - MVP focused
app.use('/v1', generateRoutes);

// Catch all - simplified error response  
app.use('*', (req: any, res: any) => {
  res.status(404).json({ error: 'Endpoint not found' });
});

// Start server
app.listen(config.port, () => {
  console.log(`Prompt service running on port ${config.port}`);
});

// Graceful shutdown
process.on('SIGTERM', () => {
  console.log('SIGTERM received, shutting down gracefully');
  process.exit(0);
});

process.on('SIGINT', () => {
  console.log('SIGINT received, shutting down gracefully');  
  process.exit(0);
});

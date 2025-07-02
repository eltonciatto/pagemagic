import Redis from 'ioredis';
import { logger } from '../utils/logger';

export class CacheService {
  private redis: Redis;
  private isConnected: boolean = false;

  constructor() {
    const redisUrl = process.env.REDIS_URL || 'redis://localhost:6379';
    
    this.redis = new Redis(redisUrl, {
      retryDelayOnFailover: 100,
      maxRetriesPerRequest: 3,
      lazyConnect: true,
    });

    this.redis.on('connect', () => {
      this.isConnected = true;
      logger.info('Connected to Redis cache');
    });

    this.redis.on('error', (error) => {
      this.isConnected = false;
      logger.error(`Redis connection error: ${error}`);
    });
  }

  async get<T>(key: string): Promise<T | null> {
    if (!this.isConnected) {
      return null;
    }

    try {
      const value = await this.redis.get(key);
      return value ? JSON.parse(value) : null;
    } catch (error) {
      logger.error(`Cache get error for key ${key}: ${error}`);
      return null;
    }
  }

  async set<T>(key: string, value: T, ttlSeconds?: number): Promise<boolean> {
    if (!this.isConnected) {
      return false;
    }

    try {
      const serialized = JSON.stringify(value);
      
      if (ttlSeconds) {
        await this.redis.setex(key, ttlSeconds, serialized);
      } else {
        await this.redis.set(key, serialized);
      }
      
      return true;
    } catch (error) {
      logger.error(`Cache set error for key ${key}: ${error}`);
      return false;
    }
  }

  async delete(key: string): Promise<boolean> {
    if (!this.isConnected) {
      return false;
    }

    try {
      const result = await this.redis.del(key);
      return result > 0;
    } catch (error) {
      logger.error(`Cache delete error for key ${key}: ${error}`);
      return false;
    }
  }

  async exists(key: string): Promise<boolean> {
    if (!this.isConnected) {
      return false;
    }

    try {
      const result = await this.redis.exists(key);
      return result === 1;
    } catch (error) {
      logger.error(`Cache exists error for key ${key}: ${error}`);
      return false;
    }
  }

  async increment(key: string, ttlSeconds?: number): Promise<number> {
    if (!this.isConnected) {
      return 0;
    }

    try {
      const result = await this.redis.incr(key);
      
      if (ttlSeconds && result === 1) {
        // Set TTL only on first increment
        await this.redis.expire(key, ttlSeconds);
      }
      
      return result;
    } catch (error) {
      logger.error(`Cache increment error for key ${key}: ${error}`);
      return 0;
    }
  }

  async flush(): Promise<boolean> {
    if (!this.isConnected) {
      return false;
    }

    try {
      await this.redis.flushdb();
      return true;
    } catch (error) {
      logger.error(`Cache flush error: ${error}`);
      return false;
    }
  }

  async disconnect(): Promise<void> {
    try {
      await this.redis.disconnect();
      this.isConnected = false;
    } catch (error) {
      logger.error(`Error disconnecting from Redis: ${error}`);
    }
  }

  isHealthy(): boolean {
    return this.isConnected;
  }
}

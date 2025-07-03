"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.CacheService = void 0;
const ioredis_1 = __importDefault(require("ioredis"));
const logger_1 = require("../utils/logger");
class CacheService {
    redis;
    isConnected = false;
    constructor() {
        const redisUrl = process.env.REDIS_URL || 'redis://localhost:6379';
        this.redis = new ioredis_1.default(redisUrl, {
            retryDelayOnFailover: 100,
            maxRetriesPerRequest: 3,
            lazyConnect: true,
        });
        this.redis.on('connect', () => {
            this.isConnected = true;
            logger_1.logger.info('Connected to Redis cache');
        });
        this.redis.on('error', (error) => {
            this.isConnected = false;
            logger_1.logger.error(`Redis connection error: ${error}`);
        });
    }
    async get(key) {
        if (!this.isConnected) {
            return null;
        }
        try {
            const value = await this.redis.get(key);
            return value ? JSON.parse(value) : null;
        }
        catch (error) {
            logger_1.logger.error(`Cache get error for key ${key}: ${error}`);
            return null;
        }
    }
    async set(key, value, ttlSeconds) {
        if (!this.isConnected) {
            return false;
        }
        try {
            const serialized = JSON.stringify(value);
            if (ttlSeconds) {
                await this.redis.setex(key, ttlSeconds, serialized);
            }
            else {
                await this.redis.set(key, serialized);
            }
            return true;
        }
        catch (error) {
            logger_1.logger.error(`Cache set error for key ${key}: ${error}`);
            return false;
        }
    }
    async delete(key) {
        if (!this.isConnected) {
            return false;
        }
        try {
            const result = await this.redis.del(key);
            return result > 0;
        }
        catch (error) {
            logger_1.logger.error(`Cache delete error for key ${key}: ${error}`);
            return false;
        }
    }
    async exists(key) {
        if (!this.isConnected) {
            return false;
        }
        try {
            const result = await this.redis.exists(key);
            return result === 1;
        }
        catch (error) {
            logger_1.logger.error(`Cache exists error for key ${key}: ${error}`);
            return false;
        }
    }
    async increment(key, ttlSeconds) {
        if (!this.isConnected) {
            return 0;
        }
        try {
            const result = await this.redis.incr(key);
            if (ttlSeconds && result === 1) {
                await this.redis.expire(key, ttlSeconds);
            }
            return result;
        }
        catch (error) {
            logger_1.logger.error(`Cache increment error for key ${key}: ${error}`);
            return 0;
        }
    }
    async flush() {
        if (!this.isConnected) {
            return false;
        }
        try {
            await this.redis.flushdb();
            return true;
        }
        catch (error) {
            logger_1.logger.error(`Cache flush error: ${error}`);
            return false;
        }
    }
    async disconnect() {
        try {
            await this.redis.disconnect();
            this.isConnected = false;
        }
        catch (error) {
            logger_1.logger.error(`Error disconnecting from Redis: ${error}`);
        }
    }
    isHealthy() {
        return this.isConnected;
    }
}
exports.CacheService = CacheService;
//# sourceMappingURL=CacheService.js.map
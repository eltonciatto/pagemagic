"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.config = void 0;
const dotenv_1 = __importDefault(require("dotenv"));
dotenv_1.default.config();
exports.config = {
    port: parseInt(process.env.PORT || '3001', 10),
    environment: process.env.NODE_ENV || 'development',
    openai: {
        apiKey: process.env.OPENAI_API_KEY || '',
        model: process.env.OPENAI_MODEL || 'gpt-4-turbo-preview',
        temperature: parseFloat(process.env.OPENAI_TEMPERATURE || '0.7'),
    },
    anthropic: {
        apiKey: process.env.ANTHROPIC_API_KEY || '',
        model: process.env.ANTHROPIC_MODEL || 'claude-3-sonnet-20240229',
        temperature: parseFloat(process.env.ANTHROPIC_TEMPERATURE || '0.7'),
    },
    vllm: {
        baseUrl: process.env.VLLM_BASE_URL || 'http://localhost:8000',
        model: process.env.VLLM_MODEL || 'meta-llama/Llama-2-70b-chat-hf',
        temperature: parseFloat(process.env.VLLM_TEMPERATURE || '0.7'),
    },
    redis: {
        host: process.env.REDIS_HOST || 'localhost',
        port: parseInt(process.env.REDIS_PORT || '6379', 10),
        password: process.env.REDIS_PASSWORD || '',
        db: parseInt(process.env.REDIS_DB || '0', 10),
    },
    nats: {
        url: process.env.NATS_URL || 'nats://localhost:4222',
        user: process.env.NATS_USER || '',
        password: process.env.NATS_PASSWORD || '',
    },
    cors: {
        origins: (process.env.CORS_ORIGINS || 'http://localhost:3000').split(','),
    },
    rateLimit: {
        windowMs: parseInt(process.env.RATE_LIMIT_WINDOW_MS || '900000', 10),
        max: parseInt(process.env.RATE_LIMIT_MAX || '100', 10),
    },
    jwt: {
        secret: process.env.JWT_SECRET || 'your-secret-key',
    },
    logging: {
        level: process.env.LOG_LEVEL || 'info',
        format: process.env.LOG_FORMAT || 'json',
    },
    metrics: {
        enabled: process.env.METRICS_ENABLED === 'true',
        prefix: process.env.METRICS_PREFIX || 'prompt_svc_',
    },
};
//# sourceMappingURL=index.js.map
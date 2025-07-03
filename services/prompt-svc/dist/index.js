"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const express_1 = __importDefault(require("express"));
const cors_1 = __importDefault(require("cors"));
const helmet_1 = __importDefault(require("helmet"));
const compression_1 = __importDefault(require("compression"));
const morgan_1 = __importDefault(require("morgan"));
const config = {
    port: process.env.PORT || 3001,
    cors: {
        origins: process.env.CORS_ORIGINS?.split(',') || ['http://localhost:3000']
    }
};
const app = (0, express_1.default)();
app.use((0, helmet_1.default)());
app.use((0, compression_1.default)());
app.use((0, cors_1.default)({
    origin: config.cors.origins,
    credentials: true,
}));
app.use((0, morgan_1.default)('combined', { stream: { write: (message) => logger.info(message.trim()) } }));
app.use(express_1.default.json({ limit: '10mb' }));
app.use(express_1.default.urlencoded({ extended: true, limit: '10mb' }));
app.use(rateLimiter);
app.use('/health', healthRoutes);
app.use('/metrics', metricsRoutes);
app.use('/api/v1/prompts', promptRoutes);
app.use(errorHandler);
app.use('*', (req, res) => {
    res.status(404).json({ error: 'Route not found' });
});
const server = app.listen(config.port, () => {
    logger.info(`Prompt service running on port ${config.port}`);
});
process.on('SIGTERM', () => {
    logger.info('SIGTERM received, shutting down gracefully');
    server.close(() => {
        logger.info('Process terminated');
        process.exit(0);
    });
});
process.on('SIGINT', () => {
    logger.info('SIGINT received, shutting down gracefully');
    server.close(() => {
        logger.info('Process terminated');
        process.exit(0);
    });
});
exports.default = app;
//# sourceMappingURL=index.js.map
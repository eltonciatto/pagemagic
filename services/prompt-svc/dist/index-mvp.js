"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const express_1 = __importDefault(require("express"));
const cors_1 = __importDefault(require("cors"));
const helmet_1 = __importDefault(require("helmet"));
const compression_1 = __importDefault(require("compression"));
const generate_1 = require("./routes/generate");
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
app.use(express_1.default.json({ limit: '10mb' }));
app.use(express_1.default.urlencoded({ extended: true, limit: '10mb' }));
app.use('/v1', generate_1.generateRoutes);
app.use('*', (req, res) => {
    res.status(404).json({ error: 'Endpoint not found' });
});
app.listen(config.port, () => {
    console.log(`Prompt service running on port ${config.port}`);
});
process.on('SIGTERM', () => {
    console.log('SIGTERM received, shutting down gracefully');
    process.exit(0);
});
process.on('SIGINT', () => {
    console.log('SIGINT received, shutting down gracefully');
    process.exit(0);
});
//# sourceMappingURL=index-mvp.js.map
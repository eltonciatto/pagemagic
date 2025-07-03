"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.logger = void 0;
console.log('Desenvolvimento do Page Magic - Prompt Service');
exports.logger = {
    info: (message) => console.log(`[INFO] ${new Date().toISOString()} - ${message}`),
    error: (message) => console.error(`[ERROR] ${new Date().toISOString()} - ${message}`),
    warn: (message) => console.warn(`[WARN] ${new Date().toISOString()} - ${message}`),
    debug: (message) => console.log(`[DEBUG] ${new Date().toISOString()} - ${message}`),
};
//# sourceMappingURL=logger.js.map
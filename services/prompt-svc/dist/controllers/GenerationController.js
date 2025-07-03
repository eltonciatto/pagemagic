"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GenerationController = void 0;
const express_validator_1 = require("express-validator");
const SiteGeneratorService_1 = require("../services/SiteGeneratorService");
class GenerationController {
    siteGenerator;
    constructor() {
        this.siteGenerator = new SiteGeneratorService_1.SiteGeneratorService();
    }
    static validateGenerateRequest = [
        (0, express_validator_1.body)('description')
            .isString()
            .isLength({ min: 10, max: 2000 })
            .withMessage('Description must be between 10 and 2000 characters'),
        (0, express_validator_1.body)('userId')
            .isString()
            .isLength({ min: 1 })
            .withMessage('User ID is required'),
        (0, express_validator_1.body)('projectId')
            .optional()
            .isString(),
        (0, express_validator_1.body)('preferences')
            .optional()
            .isObject(),
        (0, express_validator_1.body)('preferences.industry')
            .optional()
            .isString(),
        (0, express_validator_1.body)('preferences.style')
            .optional()
            .isString(),
        (0, express_validator_1.body)('preferences.features')
            .optional()
            .isArray()
    ];
    generateSite = async (req, res) => {
        try {
            const errors = (0, express_validator_1.validationResult)(req);
            if (!errors.isEmpty()) {
                return res.status(400).json({
                    success: false,
                    error: 'Validation failed',
                    details: errors.array()
                });
            }
            const request = {
                description: req.body.description,
                userId: req.body.userId,
                projectId: req.body.projectId,
                preferences: req.body.preferences
            };
            const result = await this.siteGenerator.generateSite(request);
            if (result.success) {
                res.status(200).json(result);
            }
            else {
                res.status(500).json(result);
            }
        }
        catch (error) {
            console.error('Generation controller error:', error);
            res.status(500).json({
                success: false,
                error: 'Internal server error',
                metadata: {
                    requestId: `error_${Date.now()}`,
                    tokensUsed: 0,
                    processingTime: 0,
                    model: 'unknown'
                }
            });
        }
    };
    health = async (req, res) => {
        try {
            const isHealthy = !!process.env.OPENAI_API_KEY;
            res.status(isHealthy ? 200 : 503).json({
                status: isHealthy ? 'healthy' : 'unhealthy',
                service: 'prompt-svc',
                timestamp: new Date().toISOString(),
                checks: {
                    openai_api_key: !!process.env.OPENAI_API_KEY
                }
            });
        }
        catch (error) {
            res.status(503).json({
                status: 'unhealthy',
                service: 'prompt-svc',
                timestamp: new Date().toISOString(),
                error: 'Health check failed'
            });
        }
    };
}
exports.GenerationController = GenerationController;
//# sourceMappingURL=GenerationController.js.map
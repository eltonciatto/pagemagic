"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.PromptController = void 0;
const express_validator_1 = require("express-validator");
const logger_1 = require("../utils/logger");
class PromptController {
    promptService;
    templateService;
    analyticsService;
    constructor(promptService, templateService, analyticsService) {
        this.promptService = promptService;
        this.templateService = templateService;
        this.analyticsService = analyticsService;
    }
    static validatePromptRequest = [
        (0, express_validator_1.body)('prompt').isString().isLength({ min: 1, max: 10000 }),
        (0, express_validator_1.body)('userId').isUUID(),
        (0, express_validator_1.body)('projectId').optional().isUUID(),
        (0, express_validator_1.body)('context').optional().isObject(),
        (0, express_validator_1.body)('requestedComponents').optional().isArray(),
        (0, express_validator_1.body)('outputFormat').optional().isIn(['html', 'json', 'markdown']),
        (0, express_validator_1.body)('language').optional().isString().isLength({ min: 2, max: 10 }),
    ];
    static validateBatchRequest = [
        (0, express_validator_1.body)('name').isString().isLength({ min: 1, max: 100 }),
        (0, express_validator_1.body)('userId').isUUID(),
        (0, express_validator_1.body)('requests').isArray().isLength({ min: 1, max: 100 }),
        (0, express_validator_1.body)('requests.*.prompt').isString().isLength({ min: 1, max: 10000 }),
    ];
    generateContent = async (req, res, next) => {
        try {
            const errors = (0, express_validator_1.validationResult)(req);
            if (!errors.isEmpty()) {
                return res.status(400).json({ errors: errors.array() });
            }
            const promptRequest = {
                userId: req.body.userId,
                projectId: req.body.projectId,
                prompt: req.body.prompt,
                context: req.body.context,
                requestedComponents: req.body.requestedComponents,
                outputFormat: req.body.outputFormat || 'html',
                language: req.body.language || 'en',
            };
            const result = await this.promptService.generateContent(promptRequest);
            await this.analyticsService.trackPromptRequest({
                requestId: result.id,
                userId: promptRequest.userId,
                prompt: promptRequest.prompt,
                tokensUsed: result.metadata?.tokensUsed || 0,
                cost: result.metadata?.tokensUsed ? this.calculateCost(result.metadata.tokensUsed) : 0,
                processingTime: result.metadata?.processingTime || 0,
                success: true,
                model: result.metadata?.model || 'unknown',
                timestamp: new Date(),
            });
            res.status(201).json({
                success: true,
                data: result,
            });
        }
        catch (error) {
            logger_1.logger.error('Error generating content:', error);
            if (req.body.userId) {
                await this.analyticsService.trackPromptRequest({
                    requestId: 'failed',
                    userId: req.body.userId,
                    prompt: req.body.prompt || '',
                    tokensUsed: 0,
                    cost: 0,
                    processingTime: 0,
                    success: false,
                    errorType: error instanceof Error ? error.name : 'UnknownError',
                    model: 'unknown',
                    timestamp: new Date(),
                });
            }
            next(error);
        }
    };
    generateBatch = async (req, res, next) => {
        try {
            const errors = (0, express_validator_1.validationResult)(req);
            if (!errors.isEmpty()) {
                return res.status(400).json({ errors: errors.array() });
            }
            const batchRequest = {
                userId: req.body.userId,
                name: req.body.name,
                requests: req.body.requests,
            };
            const batch = await this.promptService.createBatchRequest(batchRequest);
            this.promptService.processBatch(batch.id).catch(error => {
                logger_1.logger.error(`Error processing batch ${batch.id}:`, error);
            });
            res.status(202).json({
                success: true,
                data: batch,
                message: 'Batch request created and processing started',
            });
        }
        catch (error) {
            logger_1.logger.error('Error creating batch request:', error);
            next(error);
        }
    };
    getRequestStatus = async (req, res, next) => {
        try {
            const { requestId } = req.params;
            const request = await this.promptService.getRequest(requestId);
            if (!request) {
                return res.status(404).json({
                    success: false,
                    message: 'Request not found',
                });
            }
            res.json({
                success: true,
                data: request,
            });
        }
        catch (error) {
            logger_1.logger.error('Error getting request status:', error);
            next(error);
        }
    };
    getBatchStatus = async (req, res, next) => {
        try {
            const { batchId } = req.params;
            const batch = await this.promptService.getBatch(batchId);
            if (!batch) {
                return res.status(404).json({
                    success: false,
                    message: 'Batch not found',
                });
            }
            res.json({
                success: true,
                data: batch,
            });
        }
        catch (error) {
            logger_1.logger.error('Error getting batch status:', error);
            next(error);
        }
    };
    getUserRequests = async (req, res, next) => {
        try {
            const { userId } = req.params;
            const page = parseInt(req.query.page) || 1;
            const limit = parseInt(req.query.limit) || 20;
            const status = req.query.status;
            const { requests, total } = await this.promptService.getUserRequests(userId, { page, limit, status });
            res.json({
                success: true,
                data: {
                    requests,
                    pagination: {
                        page,
                        limit,
                        total,
                        pages: Math.ceil(total / limit),
                    },
                },
            });
        }
        catch (error) {
            logger_1.logger.error('Error getting user requests:', error);
            next(error);
        }
    };
    getTemplates = async (req, res, next) => {
        try {
            const category = req.query.category;
            const tags = req.query.tags ? req.query.tags.split(',') : undefined;
            const templates = await this.templateService.getTemplates({ category, tags });
            res.json({
                success: true,
                data: templates,
            });
        }
        catch (error) {
            logger_1.logger.error('Error getting templates:', error);
            next(error);
        }
    };
    applyTemplate = async (req, res, next) => {
        try {
            const { templateId } = req.params;
            const { variables, userId, projectId } = req.body;
            const prompt = await this.templateService.applyTemplate(templateId, variables);
            const promptRequest = {
                userId,
                projectId,
                prompt,
                context: { templateId, variables },
                outputFormat: 'html',
                language: 'en',
            };
            const result = await this.promptService.generateContent(promptRequest);
            res.status(201).json({
                success: true,
                data: result,
            });
        }
        catch (error) {
            logger_1.logger.error('Error applying template:', error);
            next(error);
        }
    };
    getUserAnalytics = async (req, res, next) => {
        try {
            const { userId } = req.params;
            const startDate = req.query.startDate ? new Date(req.query.startDate) : undefined;
            const endDate = req.query.endDate ? new Date(req.query.endDate) : undefined;
            const analytics = await this.analyticsService.getUserAnalytics(userId, {
                startDate,
                endDate,
            });
            res.json({
                success: true,
                data: analytics,
            });
        }
        catch (error) {
            logger_1.logger.error('Error getting user analytics:', error);
            next(error);
        }
    };
    cancelRequest = async (req, res, next) => {
        try {
            const { requestId } = req.params;
            const success = await this.promptService.cancelRequest(requestId);
            if (!success) {
                return res.status(404).json({
                    success: false,
                    message: 'Request not found or cannot be cancelled',
                });
            }
            res.json({
                success: true,
                message: 'Request cancelled successfully',
            });
        }
        catch (error) {
            logger_1.logger.error('Error cancelling request:', error);
            next(error);
        }
    };
    healthCheck = async (req, res) => {
        res.json({
            status: 'healthy',
            service: 'prompt-svc',
            timestamp: new Date().toISOString(),
        });
    };
    calculateCost(tokensUsed) {
        const costPerToken = 0.00002;
        return tokensUsed * costPerToken;
    }
}
exports.PromptController = PromptController;
//# sourceMappingURL=PromptController.js.map
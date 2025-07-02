import { Request, Response, NextFunction } from 'express';
import { body, param, query, validationResult } from 'express-validator';
import { PromptService } from '../services/PromptService';
import { TemplateService } from '../services/TemplateService';
import { AnalyticsService } from '../services/AnalyticsService';
import { PromptRequest, GenerationPreset } from '../types/prompt';
import { logger } from '../utils/logger';

export class PromptController {
  constructor(
    private promptService: PromptService,
    private templateService: TemplateService,
    private analyticsService: AnalyticsService
  ) {}

  // Validation middleware
  static validatePromptRequest = [
    body('prompt').isString().isLength({ min: 1, max: 10000 }),
    body('userId').isUUID(),
    body('projectId').optional().isUUID(),
    body('context').optional().isObject(),
    body('requestedComponents').optional().isArray(),
    body('outputFormat').optional().isIn(['html', 'json', 'markdown']),
    body('language').optional().isString().isLength({ min: 2, max: 10 }),
  ];

  static validateBatchRequest = [
    body('name').isString().isLength({ min: 1, max: 100 }),
    body('userId').isUUID(),
    body('requests').isArray().isLength({ min: 1, max: 100 }),
    body('requests.*.prompt').isString().isLength({ min: 1, max: 10000 }),
  ];

  // Main prompt generation endpoint
  generateContent = async (req: Request, res: Response, next: NextFunction) => {
    try {
      const errors = validationResult(req);
      if (!errors.isEmpty()) {
        return res.status(400).json({ errors: errors.array() });
      }

      const promptRequest: Omit<PromptRequest, 'id' | 'createdAt' | 'status'> = {
        userId: req.body.userId,
        projectId: req.body.projectId,
        prompt: req.body.prompt,
        context: req.body.context,
        requestedComponents: req.body.requestedComponents,
        outputFormat: req.body.outputFormat || 'html',
        language: req.body.language || 'en',
      };

      const result = await this.promptService.generateContent(promptRequest);
      
      // Track analytics
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
    } catch (error) {
      logger.error('Error generating content:', error);
      
      // Track failed request
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

  // Batch processing endpoint
  generateBatch = async (req: Request, res: Response, next: NextFunction) => {
    try {
      const errors = validationResult(req);
      if (!errors.isEmpty()) {
        return res.status(400).json({ errors: errors.array() });
      }

      const batchRequest = {
        userId: req.body.userId,
        name: req.body.name,
        requests: req.body.requests,
      };

      const batch = await this.promptService.createBatchRequest(batchRequest);
      
      // Start processing in background
      this.promptService.processBatch(batch.id).catch(error => {
        logger.error(`Error processing batch ${batch.id}:`, error);
      });

      res.status(202).json({
        success: true,
        data: batch,
        message: 'Batch request created and processing started',
      });
    } catch (error) {
      logger.error('Error creating batch request:', error);
      next(error);
    }
  };

  // Get prompt request status
  getRequestStatus = async (req: Request, res: Response, next: NextFunction) => {
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
    } catch (error) {
      logger.error('Error getting request status:', error);
      next(error);
    }
  };

  // Get batch status
  getBatchStatus = async (req: Request, res: Response, next: NextFunction) => {
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
    } catch (error) {
      logger.error('Error getting batch status:', error);
      next(error);
    }
  };

  // List user requests
  getUserRequests = async (req: Request, res: Response, next: NextFunction) => {
    try {
      const { userId } = req.params;
      const page = parseInt(req.query.page as string) || 1;
      const limit = parseInt(req.query.limit as string) || 20;
      const status = req.query.status as string;

      const { requests, total } = await this.promptService.getUserRequests(
        userId,
        { page, limit, status }
      );

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
    } catch (error) {
      logger.error('Error getting user requests:', error);
      next(error);
    }
  };

  // Get templates
  getTemplates = async (req: Request, res: Response, next: NextFunction) => {
    try {
      const category = req.query.category as string;
      const tags = req.query.tags ? (req.query.tags as string).split(',') : undefined;

      const templates = await this.templateService.getTemplates({ category, tags });

      res.json({
        success: true,
        data: templates,
      });
    } catch (error) {
      logger.error('Error getting templates:', error);
      next(error);
    }
  };

  // Apply template
  applyTemplate = async (req: Request, res: Response, next: NextFunction) => {
    try {
      const { templateId } = req.params;
      const { variables, userId, projectId } = req.body;

      const prompt = await this.templateService.applyTemplate(templateId, variables);
      
      const promptRequest: Omit<PromptRequest, 'id' | 'createdAt' | 'status'> = {
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
    } catch (error) {
      logger.error('Error applying template:', error);
      next(error);
    }
  };

  // Get user analytics
  getUserAnalytics = async (req: Request, res: Response, next: NextFunction) => {
    try {
      const { userId } = req.params;
      const startDate = req.query.startDate ? new Date(req.query.startDate as string) : undefined;
      const endDate = req.query.endDate ? new Date(req.query.endDate as string) : undefined;

      const analytics = await this.analyticsService.getUserAnalytics(userId, {
        startDate,
        endDate,
      });

      res.json({
        success: true,
        data: analytics,
      });
    } catch (error) {
      logger.error('Error getting user analytics:', error);
      next(error);
    }
  };

  // Cancel request
  cancelRequest = async (req: Request, res: Response, next: NextFunction) => {
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
    } catch (error) {
      logger.error('Error cancelling request:', error);
      next(error);
    }
  };

  // Health check
  healthCheck = async (req: Request, res: Response) => {
    res.json({
      status: 'healthy',
      service: 'prompt-svc',
      timestamp: new Date().toISOString(),
    });
  };

  private calculateCost(tokensUsed: number): number {
    // Simple cost calculation - should be configurable
    const costPerToken = 0.00002; // $0.00002 per token (example)
    return tokensUsed * costPerToken;
  }
}

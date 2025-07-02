import { Request, Response } from 'express';
import { body, validationResult } from 'express-validator';
import { SiteGeneratorService } from '../services/SiteGeneratorService';
import { SiteGenerationRequest } from '../types/generation';

export class GenerationController {
  private siteGenerator: SiteGeneratorService;

  constructor() {
    this.siteGenerator = new SiteGeneratorService();
  }

  // Validation middleware for site generation
  static validateGenerateRequest = [
    body('description')
      .isString()
      .isLength({ min: 10, max: 2000 })
      .withMessage('Description must be between 10 and 2000 characters'),
    body('userId')
      .isString()
      .isLength({ min: 1 })
      .withMessage('User ID is required'),
    body('projectId')
      .optional()
      .isString(),
    body('preferences')
      .optional()
      .isObject(),
    body('preferences.industry')
      .optional()
      .isString(),
    body('preferences.style')
      .optional()
      .isString(),
    body('preferences.features')
      .optional()
      .isArray()
  ];

  // Main generation endpoint - simplified for MVP
  generateSite = async (req: Request, res: Response) => {
    try {
      // Validate request
      const errors = validationResult(req);
      if (!errors.isEmpty()) {
        return res.status(400).json({
          success: false,
          error: 'Validation failed',
          details: errors.array()
        });
      }

      // Extract request data
      const request: SiteGenerationRequest = {
        description: req.body.description,
        userId: req.body.userId,
        projectId: req.body.projectId,
        preferences: req.body.preferences
      };

      // Generate site
      const result = await this.siteGenerator.generateSite(request);

      // Return response
      if (result.success) {
        res.status(200).json(result);
      } else {
        res.status(500).json(result);
      }

    } catch (error) {
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

  // Health check endpoint
  health = async (req: Request, res: Response) => {
    try {
      // Simple health check - could be expanded
      const isHealthy = !!process.env.OPENAI_API_KEY;
      
      res.status(isHealthy ? 200 : 503).json({
        status: isHealthy ? 'healthy' : 'unhealthy',
        service: 'prompt-svc',
        timestamp: new Date().toISOString(),
        checks: {
          openai_api_key: !!process.env.OPENAI_API_KEY
        }
      });
    } catch (error) {
      res.status(503).json({
        status: 'unhealthy',
        service: 'prompt-svc',
        timestamp: new Date().toISOString(),
        error: 'Health check failed'
      });
    }
  };
}

import { Request, Response } from 'express';
import { validationResult } from 'express-validator';
import { TranslationService } from '../services/translation';
import { 
  TranslationRequest, 
  BulkTranslationRequest,
  AutoTranslationRequest 
} from '../models/translation';
import { logger } from '../utils/logger';

export class TranslationController {
  constructor(private translationService: TranslationService) {}

  // Create or update a single translation
  async createTranslation(req: Request, res: Response): Promise<void> {
    try {
      const errors = validationResult(req);
      if (!errors.isEmpty()) {
        res.status(400).json({ errors: errors.array() });
        return;
      }

      const translationRequest: TranslationRequest = req.body;
      const translation = await this.translationService.createTranslation(translationRequest);
      
      res.status(201).json(translation);
    } catch (error) {
      logger.error('Failed to create translation', { error, body: req.body });
      res.status(500).json({ error: 'Failed to create translation' });
    }
  }

  // Bulk create/update translations
  async bulkCreateTranslations(req: Request, res: Response): Promise<void> {
    try {
      const errors = validationResult(req);
      if (!errors.isEmpty()) {
        res.status(400).json({ errors: errors.array() });
        return;
      }

      const bulkRequest: BulkTranslationRequest = req.body;
      const translations = await this.translationService.bulkCreateTranslations(bulkRequest);
      
      res.status(201).json({ 
        message: 'Translations created successfully',
        count: translations.length,
        translations 
      });
    } catch (error) {
      logger.error('Failed to bulk create translations', { error, body: req.body });
      res.status(500).json({ error: 'Failed to create translations' });
    }
  }

  // Get all translations for a project and locale
  async getTranslations(req: Request, res: Response): Promise<void> {
    try {
      const { projectId, locale } = req.params;
      
      const translations = await this.translationService.getTranslations(projectId, locale);
      
      res.json({
        project_id: projectId,
        locale,
        translations
      });
    } catch (error) {
      logger.error('Failed to get translations', { error, params: req.params });
      res.status(500).json({ error: 'Failed to get translations' });
    }
  }

  // Get a published translation bundle
  async getBundle(req: Request, res: Response): Promise<void> {
    try {
      const { projectId, locale } = req.params;
      
      const bundle = await this.translationService.getPublishedBundle(projectId, locale);
      
      if (!bundle) {
        res.status(404).json({ error: 'Bundle not found' });
        return;
      }
      
      res.json(bundle);
    } catch (error) {
      logger.error('Failed to get bundle', { error, params: req.params });
      res.status(500).json({ error: 'Failed to get bundle' });
    }
  }

  // Publish a translation bundle
  async publishBundle(req: Request, res: Response): Promise<void> {
    try {
      const { projectId, locale } = req.params;
      const { version } = req.body;
      
      const bundle = await this.translationService.publishBundle(projectId, locale, version);
      
      res.json({
        message: 'Bundle published successfully',
        bundle
      });
    } catch (error) {
      logger.error('Failed to publish bundle', { error, params: req.params, body: req.body });
      res.status(500).json({ error: 'Failed to publish bundle' });
    }
  }

  // Get translation statistics
  async getStats(req: Request, res: Response): Promise<void> {
    try {
      const { projectId, locale } = req.params;
      
      const stats = await this.translationService.getTranslationStats(projectId, locale);
      
      res.json(stats);
    } catch (error) {
      logger.error('Failed to get translation stats', { error, params: req.params });
      res.status(500).json({ error: 'Failed to get translation stats' });
    }
  }

  // Delete a translation
  async deleteTranslation(req: Request, res: Response): Promise<void> {
    try {
      const { projectId, locale, key } = req.params;
      
      const deleted = await this.translationService.deleteTranslation(projectId, locale, key);
      
      if (!deleted) {
        res.status(404).json({ error: 'Translation not found' });
        return;
      }
      
      res.json({ message: 'Translation deleted successfully' });
    } catch (error) {
      logger.error('Failed to delete translation', { error, params: req.params });
      res.status(500).json({ error: 'Failed to delete translation' });
    }
  }

  // Create a new project
  async createProject(req: Request, res: Response): Promise<void> {
    try {
      const errors = validationResult(req);
      if (!errors.isEmpty()) {
        res.status(400).json({ errors: errors.array() });
        return;
      }

      const { name, default_locale, supported_locales } = req.body;
      
      const project = await this.translationService.createProject(name, default_locale, supported_locales);
      
      res.status(201).json(project);
    } catch (error) {
      logger.error('Failed to create project', { error, body: req.body });
      res.status(500).json({ error: 'Failed to create project' });
    }
  }

  // Get project details
  async getProject(req: Request, res: Response): Promise<void> {
    try {
      const { projectId } = req.params;
      
      const project = await this.translationService.getProject(projectId);
      
      if (!project) {
        res.status(404).json({ error: 'Project not found' });
        return;
      }
      
      res.json(project);
    } catch (error) {
      logger.error('Failed to get project', { error, params: req.params });
      res.status(500).json({ error: 'Failed to get project' });
    }
  }

  // Health check
  async healthCheck(req: Request, res: Response): Promise<void> {
    res.json({
      status: 'healthy',
      service: 'i18n-svc',
      timestamp: new Date().toISOString()
    });
  }
}

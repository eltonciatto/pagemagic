import { Router } from 'express';
import { body, param } from 'express-validator';
import { TranslationController } from '../controllers/translation';

export function translationRoutes(controller: TranslationController): Router {
  const router = Router();

  // Project routes
  router.post('/projects', [
    body('name').notEmpty().withMessage('Project name is required'),
    body('default_locale').notEmpty().withMessage('Default locale is required'),
    body('supported_locales').isArray().withMessage('Supported locales must be an array')
  ], controller.createProject.bind(controller));

  router.get('/projects/:projectId', [
    param('projectId').isUUID().withMessage('Invalid project ID')
  ], controller.getProject.bind(controller));

  // Translation routes
  router.post('/translations', [
    body('project_id').isUUID().withMessage('Invalid project ID'),
    body('locale').notEmpty().withMessage('Locale is required'),
    body('key').notEmpty().withMessage('Translation key is required'),
    body('value').notEmpty().withMessage('Translation value is required')
  ], controller.createTranslation.bind(controller));

  router.post('/translations/bulk', [
    body('project_id').isUUID().withMessage('Invalid project ID'),
    body('locale').notEmpty().withMessage('Locale is required'),
    body('translations').isArray().withMessage('Translations must be an array')
  ], controller.bulkCreateTranslations.bind(controller));

  router.get('/projects/:projectId/locales/:locale/translations', [
    param('projectId').isUUID().withMessage('Invalid project ID'),
    param('locale').notEmpty().withMessage('Locale is required')
  ], controller.getTranslations.bind(controller));

  router.delete('/projects/:projectId/locales/:locale/translations/:key', [
    param('projectId').isUUID().withMessage('Invalid project ID'),
    param('locale').notEmpty().withMessage('Locale is required'),
    param('key').notEmpty().withMessage('Translation key is required')
  ], controller.deleteTranslation.bind(controller));

  // Bundle routes
  router.get('/projects/:projectId/locales/:locale/bundle', [
    param('projectId').isUUID().withMessage('Invalid project ID'),
    param('locale').notEmpty().withMessage('Locale is required')
  ], controller.getBundle.bind(controller));

  router.post('/projects/:projectId/locales/:locale/bundle/publish', [
    param('projectId').isUUID().withMessage('Invalid project ID'),
    param('locale').notEmpty().withMessage('Locale is required')
  ], controller.publishBundle.bind(controller));

  // Stats routes
  router.get('/projects/:projectId/locales/:locale/stats', [
    param('projectId').isUUID().withMessage('Invalid project ID'),
    param('locale').notEmpty().withMessage('Locale is required')
  ], controller.getStats.bind(controller));

  return router;
}

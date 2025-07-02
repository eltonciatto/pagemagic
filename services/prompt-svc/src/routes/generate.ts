import { Router } from 'express';
import { GenerationController } from '../controllers/GenerationController';

const router = Router();
const controller = new GenerationController();

// Main generation endpoint - MVP focused
router.post('/generate', 
  GenerationController.validateGenerateRequest,
  controller.generateSite
);

// Health check
router.get('/health', controller.health);

export const generateRoutes = router;

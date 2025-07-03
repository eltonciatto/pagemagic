"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.generateRoutes = void 0;
const express_1 = require("express");
const GenerationController_1 = require("../controllers/GenerationController");
const router = (0, express_1.Router)();
const controller = new GenerationController_1.GenerationController();
router.post('/generate', GenerationController_1.GenerationController.validateGenerateRequest, controller.generateSite);
router.get('/health', controller.health);
exports.generateRoutes = router;
//# sourceMappingURL=generate.js.map
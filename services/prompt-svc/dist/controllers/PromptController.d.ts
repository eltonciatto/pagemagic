import { Request, Response, NextFunction } from 'express';
import { PromptService } from '../services/PromptService';
import { TemplateService } from '../services/TemplateService';
import { AnalyticsService } from '../services/AnalyticsService';
export declare class PromptController {
    private promptService;
    private templateService;
    private analyticsService;
    constructor(promptService: PromptService, templateService: TemplateService, analyticsService: AnalyticsService);
    static validatePromptRequest: import("express-validator").ValidationChain[];
    static validateBatchRequest: import("express-validator").ValidationChain[];
    generateContent: (req: Request, res: Response, next: NextFunction) => Promise<Response<any, Record<string, any>> | undefined>;
    generateBatch: (req: Request, res: Response, next: NextFunction) => Promise<Response<any, Record<string, any>> | undefined>;
    getRequestStatus: (req: Request, res: Response, next: NextFunction) => Promise<Response<any, Record<string, any>> | undefined>;
    getBatchStatus: (req: Request, res: Response, next: NextFunction) => Promise<Response<any, Record<string, any>> | undefined>;
    getUserRequests: (req: Request, res: Response, next: NextFunction) => Promise<void>;
    getTemplates: (req: Request, res: Response, next: NextFunction) => Promise<void>;
    applyTemplate: (req: Request, res: Response, next: NextFunction) => Promise<void>;
    getUserAnalytics: (req: Request, res: Response, next: NextFunction) => Promise<void>;
    cancelRequest: (req: Request, res: Response, next: NextFunction) => Promise<Response<any, Record<string, any>> | undefined>;
    healthCheck: (req: Request, res: Response) => Promise<void>;
    private calculateCost;
}
//# sourceMappingURL=PromptController.d.ts.map
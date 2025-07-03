import { Request, Response } from 'express';
export declare class GenerationController {
    private siteGenerator;
    constructor();
    static validateGenerateRequest: import("express-validator").ValidationChain[];
    generateSite: (req: Request, res: Response) => Promise<Response<any, Record<string, any>> | undefined>;
    health: (req: Request, res: Response) => Promise<void>;
}
//# sourceMappingURL=GenerationController.d.ts.map
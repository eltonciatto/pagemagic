import { SiteGenerationRequest, GenerationApiResponse } from '../types/generation';
export declare class SiteGeneratorService {
    private client;
    constructor();
    generateSite(request: SiteGenerationRequest): Promise<GenerationApiResponse>;
    private buildSystemPrompt;
    private buildUserPrompt;
    private parseAndValidateResponse;
}
//# sourceMappingURL=SiteGeneratorService.d.ts.map
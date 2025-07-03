import { PromptRequest, PromptResponse, BatchRequest } from '../types/prompt';
interface GenerationOptions {
    modelId?: string;
    temperature?: number;
    maxTokens?: number;
    useCache?: boolean;
}
export declare class PromptService {
    private openaiService;
    private anthropicService;
    private templateService;
    private cacheService;
    private requestStore;
    private batchStore;
    constructor();
    generateContent(request: Omit<PromptRequest, 'id' | 'createdAt' | 'status'>, options?: GenerationOptions): Promise<PromptResponse>;
    createBatchRequest(request: {
        userId: string;
        name: string;
        requests: Omit<PromptRequest, 'id' | 'createdAt' | 'status'>[];
    }): Promise<BatchRequest>;
    processBatch(batchId: string): Promise<void>;
    getRequest(requestId: string): Promise<PromptRequest | null>;
    getBatch(batchId: string): Promise<BatchRequest | null>;
    getUserRequests(userId: string, options?: {
        page?: number;
        limit?: number;
        status?: string;
    }): Promise<{
        requests: PromptRequest[];
        total: number;
    }>;
    cancelRequest(requestId: string): Promise<boolean>;
    private buildEnhancedPrompt;
    private selectAIService;
    private parseResponse;
    private extractComponents;
    private getCacheKey;
    private chunkArray;
}
export {};
//# sourceMappingURL=PromptService.d.ts.map
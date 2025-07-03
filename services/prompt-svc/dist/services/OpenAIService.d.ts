interface GenerationOptions {
    temperature?: number;
    maxTokens?: number;
    requestedComponents?: string[];
    outputFormat?: 'html' | 'json' | 'markdown';
}
interface GenerationResponse {
    content: string;
    tokensUsed: number;
    model: string;
    confidence?: number;
}
export declare class OpenAIService {
    private client;
    constructor();
    generateContent(prompt: string, options?: GenerationOptions): Promise<GenerationResponse>;
    generateWithFunctions(prompt: string, functions: any[], options?: GenerationOptions): Promise<GenerationResponse>;
    generateEmbedding(text: string): Promise<number[]>;
    private buildSystemPrompt;
    private calculateConfidence;
}
export {};
//# sourceMappingURL=OpenAIService.d.ts.map
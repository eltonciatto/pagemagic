export interface PromptRequest {
    id: string;
    userId: string;
    projectId?: string;
    prompt: string;
    context?: {
        websiteType?: string;
        industry?: string;
        targetAudience?: string;
        brandColors?: string[];
        existingContent?: string;
        preferences?: Record<string, any>;
    };
    requestedComponents?: string[];
    outputFormat?: 'html' | 'json' | 'markdown';
    language?: string;
    createdAt: Date;
    status: 'pending' | 'processing' | 'completed' | 'failed';
}
export interface PromptResponse {
    id: string;
    requestId: string;
    content: string;
    components?: ComponentOutput[];
    metadata?: {
        tokensUsed: number;
        processingTime: number;
        model: string;
        confidence: number;
    };
    generatedAt: Date;
}
export interface ComponentOutput {
    type: 'header' | 'navigation' | 'hero' | 'section' | 'footer' | 'sidebar' | 'form' | 'gallery';
    name: string;
    html: string;
    css?: string;
    props?: Record<string, any>;
    assets?: string[];
}
export interface PromptTemplate {
    id: string;
    name: string;
    description: string;
    category: string;
    template: string;
    variables: string[];
    examples: string[];
    tags: string[];
    isActive: boolean;
    createdAt: Date;
    updatedAt: Date;
}
export interface AIModel {
    id: string;
    name: string;
    provider: 'openai' | 'anthropic' | 'google' | 'custom';
    model: string;
    maxTokens: number;
    pricing: {
        inputTokens: number;
        outputTokens: number;
    };
    capabilities: string[];
    isDefault: boolean;
    isActive: boolean;
}
export interface GenerationPreset {
    id: string;
    name: string;
    description: string;
    modelId: string;
    temperature: number;
    maxTokens: number;
    topP: number;
    frequencyPenalty: number;
    presencePenalty: number;
    systemPrompt: string;
    isPublic: boolean;
    userId: string;
    createdAt: Date;
    updatedAt: Date;
}
export interface PromptAnalytics {
    requestId: string;
    userId: string;
    prompt: string;
    tokensUsed: number;
    cost: number;
    processingTime: number;
    success: boolean;
    errorType?: string;
    model: string;
    timestamp: Date;
}
export interface BatchRequest {
    id: string;
    userId: string;
    name: string;
    requests: PromptRequest[];
    status: 'pending' | 'processing' | 'completed' | 'failed' | 'cancelled';
    progress: {
        total: number;
        completed: number;
        failed: number;
    };
    createdAt: Date;
    completedAt?: Date;
}
//# sourceMappingURL=prompt.d.ts.map
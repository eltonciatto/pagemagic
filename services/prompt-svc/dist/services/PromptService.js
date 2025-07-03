"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.PromptService = void 0;
const uuid_1 = require("uuid");
const OpenAIService_1 = require("./OpenAIService");
const AnthropicService_1 = require("./AnthropicService");
const TemplateService_1 = require("./TemplateService");
const CacheService_1 = require("./CacheService");
const logger_1 = require("../utils/logger");
class PromptService {
    openaiService;
    anthropicService;
    templateService;
    cacheService;
    requestStore = new Map();
    batchStore = new Map();
    constructor() {
        this.openaiService = new OpenAIService_1.OpenAIService();
        this.anthropicService = new AnthropicService_1.AnthropicService();
        this.templateService = new TemplateService_1.TemplateService();
        this.cacheService = new CacheService_1.CacheService();
    }
    async generateContent(request, options = {}) {
        const requestId = (0, uuid_1.v4)();
        const startTime = Date.now();
        const fullRequest = {
            ...request,
            id: requestId,
            createdAt: new Date(),
            status: 'processing',
        };
        this.requestStore.set(requestId, fullRequest);
        try {
            logger_1.logger.info(`Processing prompt request ${requestId}`);
            if (options.useCache !== false) {
                const cached = await this.cacheService.get(this.getCacheKey(request));
                if (cached) {
                    logger_1.logger.info(`Cache hit for request ${requestId}`);
                    fullRequest.status = 'completed';
                    this.requestStore.set(requestId, fullRequest);
                    return cached;
                }
            }
            const enhancedPrompt = await this.buildEnhancedPrompt(request);
            const aiService = this.selectAIService(options.modelId);
            const response = await aiService.generateContent(enhancedPrompt, {
                temperature: options.temperature || 0.7,
                maxTokens: options.maxTokens || 4000,
                requestedComponents: request.requestedComponents,
                outputFormat: request.outputFormat,
            });
            const structuredResponse = await this.parseResponse(response, request);
            const promptResponse = {
                id: (0, uuid_1.v4)(),
                requestId,
                content: structuredResponse.content,
                components: structuredResponse.components,
                metadata: {
                    tokensUsed: response.tokensUsed,
                    processingTime: Date.now() - startTime,
                    model: response.model,
                    confidence: response.confidence || 0.8,
                },
                generatedAt: new Date(),
            };
            if (options.useCache !== false) {
                await this.cacheService.set(this.getCacheKey(request), promptResponse, 3600);
            }
            fullRequest.status = 'completed';
            this.requestStore.set(requestId, fullRequest);
            logger_1.logger.info(`Completed prompt request ${requestId} in ${Date.now() - startTime}ms`);
            return promptResponse;
        }
        catch (error) {
            logger_1.logger.error(`Error processing prompt request ${requestId}: ${error}`);
            fullRequest.status = 'failed';
            this.requestStore.set(requestId, fullRequest);
            throw error;
        }
    }
    async createBatchRequest(request) {
        const batchId = (0, uuid_1.v4)();
        const batch = {
            id: batchId,
            userId: request.userId,
            name: request.name,
            requests: request.requests.map(req => ({
                ...req,
                id: (0, uuid_1.v4)(),
                createdAt: new Date(),
                status: 'pending',
            })),
            status: 'pending',
            progress: {
                total: request.requests.length,
                completed: 0,
                failed: 0,
            },
            createdAt: new Date(),
        };
        this.batchStore.set(batchId, batch);
        return batch;
    }
    async processBatch(batchId) {
        const batch = this.batchStore.get(batchId);
        if (!batch) {
            throw new Error(`Batch ${batchId} not found`);
        }
        batch.status = 'processing';
        this.batchStore.set(batchId, batch);
        logger_1.logger.info(`Starting batch processing for ${batchId} with ${batch.requests.length} requests`);
        try {
            const concurrency = 5;
            const chunks = this.chunkArray(batch.requests, concurrency);
            for (const chunk of chunks) {
                await Promise.allSettled(chunk.map(async (request) => {
                    try {
                        await this.generateContent(request);
                        batch.progress.completed++;
                    }
                    catch (error) {
                        logger_1.logger.error(`Failed to process request ${request.id} in batch ${batchId}: ${error}`);
                        batch.progress.failed++;
                    }
                    this.batchStore.set(batchId, batch);
                }));
            }
            batch.status = 'completed';
            batch.completedAt = new Date();
            this.batchStore.set(batchId, batch);
            logger_1.logger.info(`Completed batch processing for ${batchId}`);
        }
        catch (error) {
            batch.status = 'failed';
            this.batchStore.set(batchId, batch);
            logger_1.logger.error(`Batch processing failed for ${batchId}: ${error}`);
            throw error;
        }
    }
    async getRequest(requestId) {
        return this.requestStore.get(requestId) || null;
    }
    async getBatch(batchId) {
        return this.batchStore.get(batchId) || null;
    }
    async getUserRequests(userId, options = {}) {
        const { page = 1, limit = 20, status } = options;
        const allRequests = Array.from(this.requestStore.values())
            .filter(req => req.userId === userId)
            .filter(req => !status || req.status === status)
            .sort((a, b) => b.createdAt.getTime() - a.createdAt.getTime());
        const total = allRequests.length;
        const startIndex = (page - 1) * limit;
        const requests = allRequests.slice(startIndex, startIndex + limit);
        return { requests, total };
    }
    async cancelRequest(requestId) {
        const request = this.requestStore.get(requestId);
        if (!request || request.status !== 'pending') {
            return false;
        }
        request.status = 'failed';
        this.requestStore.set(requestId, request);
        return true;
    }
    async buildEnhancedPrompt(request) {
        let enhancedPrompt = request.prompt;
        if (request.context) {
            const contextStr = Object.entries(request.context)
                .map(([key, value]) => `${key}: ${JSON.stringify(value)}`)
                .join('\n');
            enhancedPrompt = `Context:\n${contextStr}\n\nUser Request:\n${request.prompt}`;
        }
        if (request.requestedComponents && request.requestedComponents.length > 0) {
            enhancedPrompt += `\n\nRequired Components: ${request.requestedComponents.join(', ')}`;
        }
        if (request.outputFormat) {
            enhancedPrompt += `\n\nOutput Format: ${request.outputFormat}`;
        }
        if (request.language && request.language !== 'en') {
            enhancedPrompt += `\n\nLanguage: ${request.language}`;
        }
        return enhancedPrompt;
    }
    selectAIService(modelId) {
        if (modelId?.includes('claude') || modelId?.includes('anthropic')) {
            return this.anthropicService;
        }
        return this.openaiService;
    }
    async parseResponse(response, request) {
        if (request.outputFormat === 'json') {
            try {
                const parsed = JSON.parse(response.content);
                return {
                    content: response.content,
                    components: parsed.components || [],
                };
            }
            catch (error) {
                logger_1.logger.error('Failed to parse JSON response:', error);
            }
        }
        if (request.outputFormat === 'html' && request.requestedComponents) {
            const components = await this.extractComponents(response.content, request.requestedComponents);
            return {
                content: response.content,
                components,
            };
        }
        return { content: response.content };
    }
    async extractComponents(html, requestedComponents) {
        const components = [];
        for (const componentType of requestedComponents) {
            const regex = new RegExp(`<!-- ${componentType} -->(.*?)<!-- /${componentType} -->`, 'gis');
            const matches = html.match(regex);
            if (matches) {
                matches.forEach((match, index) => {
                    components.push({
                        type: componentType,
                        name: `${componentType}-${index + 1}`,
                        html: match,
                        props: {},
                        assets: [],
                    });
                });
            }
        }
        return components;
    }
    getCacheKey(request) {
        const key = JSON.stringify({
            prompt: request.prompt,
            context: request.context,
            components: request.requestedComponents,
            format: request.outputFormat,
            language: request.language,
        });
        let hash = 0;
        for (let i = 0; i < key.length; i++) {
            const char = key.charCodeAt(i);
            hash = ((hash << 5) - hash) + char;
            hash = hash & hash;
        }
        return `prompt:${Math.abs(hash)}`;
    }
    chunkArray(array, size) {
        const chunks = [];
        for (let i = 0; i < array.length; i += size) {
            chunks.push(array.slice(i, i + size));
        }
        return chunks;
    }
}
exports.PromptService = PromptService;
//# sourceMappingURL=PromptService.js.map
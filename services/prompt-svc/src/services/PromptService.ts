import { v4 as uuidv4 } from 'uuid';
import { PromptRequest, PromptResponse, BatchRequest, ComponentOutput } from '../types/prompt';
import { OpenAIService } from './OpenAIService';
import { AnthropicService } from './AnthropicService';
import { TemplateService } from './TemplateService';
import { CacheService } from './CacheService';
import { logger } from '../utils/logger';

interface GenerationOptions {
  modelId?: string;
  temperature?: number;
  maxTokens?: number;
  useCache?: boolean;
}

export class PromptService {
  private openaiService: OpenAIService;
  private anthropicService: AnthropicService;
  private templateService: TemplateService;
  private cacheService: CacheService;
  private requestStore: Map<string, PromptRequest> = new Map();
  private batchStore: Map<string, BatchRequest> = new Map();

  constructor() {
    this.openaiService = new OpenAIService();
    this.anthropicService = new AnthropicService();
    this.templateService = new TemplateService();
    this.cacheService = new CacheService();
  }

  async generateContent(
    request: Omit<PromptRequest, 'id' | 'createdAt' | 'status'>,
    options: GenerationOptions = {}
  ): Promise<PromptResponse> {
    const requestId = uuidv4();
    const startTime = Date.now();

    const fullRequest: PromptRequest = {
      ...request,
      id: requestId,
      createdAt: new Date(),
      status: 'processing',
    };

    this.requestStore.set(requestId, fullRequest);

    try {
      logger.info(`Processing prompt request ${requestId}`);

      // Check cache first if enabled
      if (options.useCache !== false) {
        const cached = await this.cacheService.get(this.getCacheKey(request));
        if (cached) {
          logger.info(`Cache hit for request ${requestId}`);
          fullRequest.status = 'completed';
          this.requestStore.set(requestId, fullRequest);
          return cached;
        }
      }

      // Build enhanced prompt with context
      const enhancedPrompt = await this.buildEnhancedPrompt(request);

      // Choose appropriate AI service based on modelId or default
      const aiService = this.selectAIService(options.modelId);
      
      // Generate content
      const response = await aiService.generateContent(enhancedPrompt, {
        temperature: options.temperature || 0.7,
        maxTokens: options.maxTokens || 4000,
        requestedComponents: request.requestedComponents,
        outputFormat: request.outputFormat,
      });

      // Parse and structure the response
      const structuredResponse = await this.parseResponse(response, request);

      const promptResponse: PromptResponse = {
        id: uuidv4(),
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

      // Cache the response
      if (options.useCache !== false) {
        await this.cacheService.set(
          this.getCacheKey(request),
          promptResponse,
          3600 // 1 hour cache
        );
      }

      fullRequest.status = 'completed';
      this.requestStore.set(requestId, fullRequest);

      logger.info(`Completed prompt request ${requestId} in ${Date.now() - startTime}ms`);
      return promptResponse;

    } catch (error) {
      logger.error(`Error processing prompt request ${requestId}: ${error}`);
      fullRequest.status = 'failed';
      this.requestStore.set(requestId, fullRequest);
      throw error;
    }
  }

  async createBatchRequest(request: {
    userId: string;
    name: string;
    requests: Omit<PromptRequest, 'id' | 'createdAt' | 'status'>[];
  }): Promise<BatchRequest> {
    const batchId = uuidv4();
    
    const batch: BatchRequest = {
      id: batchId,
      userId: request.userId,
      name: request.name,
      requests: request.requests.map(req => ({
        ...req,
        id: uuidv4(),
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

  async processBatch(batchId: string): Promise<void> {
    const batch = this.batchStore.get(batchId);
    if (!batch) {
      throw new Error(`Batch ${batchId} not found`);
    }

    batch.status = 'processing';
    this.batchStore.set(batchId, batch);

    logger.info(`Starting batch processing for ${batchId} with ${batch.requests.length} requests`);

    try {
      // Process requests in parallel with concurrency limit
      const concurrency = 5;
      const chunks = this.chunkArray(batch.requests, concurrency);

      for (const chunk of chunks) {
        await Promise.allSettled(
          chunk.map(async (request) => {
            try {
              await this.generateContent(request);
              batch.progress.completed++;
            } catch (error) {
              logger.error(`Failed to process request ${request.id} in batch ${batchId}: ${error}`);
              batch.progress.failed++;
            }
            this.batchStore.set(batchId, batch);
          })
        );
      }

      batch.status = 'completed';
      batch.completedAt = new Date();
      this.batchStore.set(batchId, batch);

      logger.info(`Completed batch processing for ${batchId}`);
    } catch (error) {
      batch.status = 'failed';
      this.batchStore.set(batchId, batch);
      logger.error(`Batch processing failed for ${batchId}: ${error}`);
      throw error;
    }
  }

  async getRequest(requestId: string): Promise<PromptRequest | null> {
    return this.requestStore.get(requestId) || null;
  }

  async getBatch(batchId: string): Promise<BatchRequest | null> {
    return this.batchStore.get(batchId) || null;
  }

  async getUserRequests(
    userId: string,
    options: { page?: number; limit?: number; status?: string } = {}
  ): Promise<{ requests: PromptRequest[]; total: number }> {
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

  async cancelRequest(requestId: string): Promise<boolean> {
    const request = this.requestStore.get(requestId);
    if (!request || request.status !== 'pending') {
      return false;
    }

    request.status = 'failed';
    this.requestStore.set(requestId, request);
    return true;
  }

  private async buildEnhancedPrompt(request: Omit<PromptRequest, 'id' | 'createdAt' | 'status'>): Promise<string> {
    let enhancedPrompt = request.prompt;

    // Add context information
    if (request.context) {
      const contextStr = Object.entries(request.context)
        .map(([key, value]) => `${key}: ${JSON.stringify(value)}`)
        .join('\n');
      
      enhancedPrompt = `Context:\n${contextStr}\n\nUser Request:\n${request.prompt}`;
    }

    // Add component requirements
    if (request.requestedComponents && request.requestedComponents.length > 0) {
      enhancedPrompt += `\n\nRequired Components: ${request.requestedComponents.join(', ')}`;
    }

    // Add format requirements
    if (request.outputFormat) {
      enhancedPrompt += `\n\nOutput Format: ${request.outputFormat}`;
    }

    // Add language requirement
    if (request.language && request.language !== 'en') {
      enhancedPrompt += `\n\nLanguage: ${request.language}`;
    }

    return enhancedPrompt;
  }

  private selectAIService(modelId?: string) {
    // Default to OpenAI for now
    // Could be enhanced to route based on model capabilities
    if (modelId?.includes('claude') || modelId?.includes('anthropic')) {
      return this.anthropicService;
    }
    return this.openaiService;
  }

  private async parseResponse(
    response: any,
    request: Omit<PromptRequest, 'id' | 'createdAt' | 'status'>
  ): Promise<{ content: string; components?: ComponentOutput[] }> {
    if (request.outputFormat === 'json') {
      try {
        const parsed = JSON.parse(response.content);
        return {
          content: response.content,
          components: parsed.components || [],
        };
      } catch (error) {
        logger.error('Failed to parse JSON response:', error);
      }
    }

    // For HTML output, try to extract components
    if (request.outputFormat === 'html' && request.requestedComponents) {
      const components = await this.extractComponents(response.content, request.requestedComponents);
      return {
        content: response.content,
        components,
      };
    }

    return { content: response.content };
  }

  private async extractComponents(html: string, requestedComponents: string[]): Promise<ComponentOutput[]> {
    // Simple component extraction logic
    // This could be enhanced with more sophisticated parsing
    const components: ComponentOutput[] = [];

    for (const componentType of requestedComponents) {
      // Extract sections that might be this component type
      const regex = new RegExp(`<!-- ${componentType} -->(.*?)<!-- /${componentType} -->`, 'gis');
      const matches = html.match(regex);

      if (matches) {
        matches.forEach((match, index) => {
          components.push({
            type: componentType as any,
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

  private getCacheKey(request: Omit<PromptRequest, 'id' | 'createdAt' | 'status'>): string {
    const key = JSON.stringify({
      prompt: request.prompt,
      context: request.context,
      components: request.requestedComponents,
      format: request.outputFormat,
      language: request.language,
    });
    
    // Simple hash function for cache key
    let hash = 0;
    for (let i = 0; i < key.length; i++) {
      const char = key.charCodeAt(i);
      hash = ((hash << 5) - hash) + char;
      hash = hash & hash; // Convert to 32-bit integer
    }
    
    return `prompt:${Math.abs(hash)}`;
  }

  private chunkArray<T>(array: T[], size: number): T[][] {
    const chunks: T[][] = [];
    for (let i = 0; i < array.length; i += size) {
      chunks.push(array.slice(i, i + size));
    }
    return chunks;
  }
}

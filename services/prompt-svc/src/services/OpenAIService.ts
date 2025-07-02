import OpenAI from 'openai';
import { logger } from '../utils/logger';

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

export class OpenAIService {
  private client: OpenAI;

  constructor() {
    const apiKey = process.env.OPENAI_API_KEY;
    if (!apiKey) {
      throw new Error('OPENAI_API_KEY environment variable is required');
    }

    this.client = new OpenAI({
      apiKey: apiKey,
    });
  }

  async generateContent(prompt: string, options: GenerationOptions = {}): Promise<GenerationResponse> {
    try {
      const systemPrompt = this.buildSystemPrompt(options);
      
      const response = await this.client.chat.completions.create({
        model: 'gpt-4-turbo-preview',
        messages: [
          { role: 'system', content: systemPrompt },
          { role: 'user', content: prompt }
        ],
        temperature: options.temperature || 0.7,
        max_tokens: options.maxTokens || 4000,
        top_p: 1,
        frequency_penalty: 0,
        presence_penalty: 0,
      });

      const choice = response.choices[0];
      if (!choice || !choice.message) {
        throw new Error('No response generated from OpenAI');
      }

      return {
        content: choice.message.content || '',
        tokensUsed: response.usage?.total_tokens || 0,
        model: response.model,
        confidence: this.calculateConfidence(choice),
      };
    } catch (error) {
      logger.error(`OpenAI API error: ${error}`);
      throw error;
    }
  }

  async generateWithFunctions(
    prompt: string,
    functions: any[],
    options: GenerationOptions = {}
  ): Promise<GenerationResponse> {
    try {
      const response = await this.client.chat.completions.create({
        model: 'gpt-4-turbo-preview',
        messages: [
          { role: 'user', content: prompt }
        ],
        functions: functions,
        function_call: 'auto',
        temperature: options.temperature || 0.7,
        max_tokens: options.maxTokens || 4000,
      });

      const choice = response.choices[0];
      if (!choice || !choice.message) {
        throw new Error('No response generated from OpenAI');
      }

      let content = choice.message.content || '';
      
      // If function was called, include the function call in the response
      if (choice.message.function_call) {
        content = JSON.stringify({
          type: 'function_call',
          function: choice.message.function_call.name,
          arguments: JSON.parse(choice.message.function_call.arguments || '{}'),
          content: content,
        });
      }

      return {
        content,
        tokensUsed: response.usage?.total_tokens || 0,
        model: response.model,
        confidence: this.calculateConfidence(choice),
      };
    } catch (error) {
      logger.error(`OpenAI Functions API error: ${error}`);
      throw error;
    }
  }

  async generateEmbedding(text: string): Promise<number[]> {
    try {
      const response = await this.client.embeddings.create({
        model: 'text-embedding-ada-002',
        input: text,
      });

      return response.data[0].embedding;
    } catch (error) {
      logger.error(`OpenAI Embedding API error: ${error}`);
      throw error;
    }
  }

  private buildSystemPrompt(options: GenerationOptions): string {
    let systemPrompt = 'You are an expert web developer and designer specializing in creating modern, responsive websites.';

    if (options.outputFormat === 'html') {
      systemPrompt += ' Generate clean, semantic HTML with modern CSS. Use best practices for accessibility and SEO.';
    } else if (options.outputFormat === 'json') {
      systemPrompt += ' Respond with valid JSON containing the structured information requested.';
    } else if (options.outputFormat === 'markdown') {
      systemPrompt += ' Respond with well-formatted Markdown content.';
    }

    if (options.requestedComponents && options.requestedComponents.length > 0) {
      systemPrompt += ` Include these specific components: ${options.requestedComponents.join(', ')}.`;
      systemPrompt += ' Wrap each component with HTML comments like <!-- componentName --> content <!-- /componentName -->';
    }

    systemPrompt += ' Focus on modern design trends, user experience, and mobile responsiveness.';
    systemPrompt += ' Use Tailwind CSS classes for styling when generating HTML.';

    return systemPrompt;
  }

  private calculateConfidence(choice: any): number {
    // Simple confidence calculation based on finish reason and other factors
    if (choice.finish_reason === 'stop') {
      return 0.9;
    } else if (choice.finish_reason === 'length') {
      return 0.7;
    } else {
      return 0.5;
    }
  }
}

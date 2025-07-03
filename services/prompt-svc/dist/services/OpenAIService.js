"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.OpenAIService = void 0;
const openai_1 = __importDefault(require("openai"));
const logger_1 = require("../utils/logger");
class OpenAIService {
    client;
    constructor() {
        const apiKey = process.env.OPENAI_API_KEY;
        if (!apiKey) {
            throw new Error('OPENAI_API_KEY environment variable is required');
        }
        this.client = new openai_1.default({
            apiKey: apiKey,
        });
    }
    async generateContent(prompt, options = {}) {
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
        }
        catch (error) {
            logger_1.logger.error(`OpenAI API error: ${error}`);
            throw error;
        }
    }
    async generateWithFunctions(prompt, functions, options = {}) {
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
        }
        catch (error) {
            logger_1.logger.error(`OpenAI Functions API error: ${error}`);
            throw error;
        }
    }
    async generateEmbedding(text) {
        try {
            const response = await this.client.embeddings.create({
                model: 'text-embedding-ada-002',
                input: text,
            });
            return response.data[0].embedding;
        }
        catch (error) {
            logger_1.logger.error(`OpenAI Embedding API error: ${error}`);
            throw error;
        }
    }
    buildSystemPrompt(options) {
        let systemPrompt = 'You are an expert web developer and designer specializing in creating modern, responsive websites.';
        if (options.outputFormat === 'html') {
            systemPrompt += ' Generate clean, semantic HTML with modern CSS. Use best practices for accessibility and SEO.';
        }
        else if (options.outputFormat === 'json') {
            systemPrompt += ' Respond with valid JSON containing the structured information requested.';
        }
        else if (options.outputFormat === 'markdown') {
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
    calculateConfidence(choice) {
        if (choice.finish_reason === 'stop') {
            return 0.9;
        }
        else if (choice.finish_reason === 'length') {
            return 0.7;
        }
        else {
            return 0.5;
        }
    }
}
exports.OpenAIService = OpenAIService;
//# sourceMappingURL=OpenAIService.js.map
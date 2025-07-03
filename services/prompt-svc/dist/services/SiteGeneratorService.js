"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.SiteGeneratorService = void 0;
const openai_1 = __importDefault(require("openai"));
const logger_1 = require("../utils/logger");
class SiteGeneratorService {
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
    async generateSite(request) {
        const startTime = Date.now();
        const requestId = `gen_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
        try {
            logger_1.logger.info(`Starting site generation for request ${requestId}`);
            const systemPrompt = this.buildSystemPrompt();
            const userPrompt = this.buildUserPrompt(request);
            const response = await this.client.chat.completions.create({
                model: 'gpt-4-turbo-preview',
                messages: [
                    { role: 'system', content: systemPrompt },
                    { role: 'user', content: userPrompt }
                ],
                temperature: 0.7,
                max_tokens: 4000,
                response_format: { type: 'json_object' }
            });
            const choice = response.choices[0];
            if (!choice?.message?.content) {
                throw new Error('No content generated from OpenAI');
            }
            const generatedSite = this.parseAndValidateResponse(choice.message.content);
            const processingTime = Date.now() - startTime;
            const tokensUsed = response.usage?.total_tokens || 0;
            logger_1.logger.info(`Site generation completed for request ${requestId} in ${processingTime}ms using ${tokensUsed} tokens`);
            return {
                success: true,
                data: generatedSite,
                metadata: {
                    requestId,
                    tokensUsed,
                    processingTime,
                    model: response.model
                }
            };
        }
        catch (error) {
            const processingTime = Date.now() - startTime;
            logger_1.logger.error(`Site generation failed for request ${requestId}:`, error);
            return {
                success: false,
                error: error instanceof Error ? error.message : 'Unknown error occurred',
                metadata: {
                    requestId,
                    tokensUsed: 0,
                    processingTime,
                    model: 'gpt-4-turbo-preview'
                }
            };
        }
    }
    buildSystemPrompt() {
        return `You are a professional web designer and developer. Your task is to generate complete, modern, responsive website structures based on user descriptions.

You must respond with a valid JSON object that matches this exact structure:

{
  "site": {
    "title": "Website title",
    "description": "SEO description",
    "sections": [
      {
        "id": "unique_id",
        "type": "hero|features|testimonials|cta|about|contact",
        "title": "Section title",
        "content": "Rich content with HTML tags allowed",
        "order": 1,
        "cta": {
          "text": "Button text",
          "link": "#contact",
          "style": "primary"
        },
        "features": [
          {
            "title": "Feature name",
            "description": "Feature description",
            "icon": "icon-name"
          }
        ],
        "testimonials": [
          {
            "quote": "Testimonial text",
            "author": "Author Name",
            "role": "Job Title",
            "company": "Company Name"
          }
        ]
      }
    ],
    "theme": {
      "primaryColor": "#3B82F6",
      "secondaryColor": "#EF4444",
      "fontFamily": "sans",
      "layout": "modern"
    },
    "metadata": {
      "industry": "detected industry",
      "target_audience": "target audience",
      "tone": "professional",
      "language": "en"
    }
  }
}

Guidelines:
- Always include at least a hero section and a CTA section
- Generate relevant, professional content for each section
- Use appropriate colors and themes for the business type
- Include realistic testimonials if relevant
- Make content engaging and conversion-focused
- Use modern web design principles
- Ensure all fields are properly filled
- Generate 3-6 sections total for a complete site`;
    }
    buildUserPrompt(request) {
        let prompt = `Generate a complete website structure for: "${request.description}"`;
        if (request.preferences?.industry) {
            prompt += `\nIndustry: ${request.preferences.industry}`;
        }
        if (request.preferences?.style) {
            prompt += `\nPreferred style: ${request.preferences.style}`;
        }
        if (request.preferences?.features?.length) {
            prompt += `\nRequested features: ${request.preferences.features.join(', ')}`;
        }
        prompt += '\n\nGenerate a professional, conversion-optimized website structure that would be perfect for this business.';
        return prompt;
    }
    parseAndValidateResponse(content) {
        try {
            const parsed = JSON.parse(content);
            if (!parsed.site) {
                throw new Error('Response missing site object');
            }
            if (!parsed.site.title || !parsed.site.description) {
                throw new Error('Response missing required title or description');
            }
            if (!Array.isArray(parsed.site.sections) || parsed.site.sections.length === 0) {
                throw new Error('Response missing sections array');
            }
            for (const section of parsed.site.sections) {
                if (!section.id || !section.type || !section.title) {
                    throw new Error('Section missing required fields');
                }
            }
            const sectionTypes = parsed.site.sections.map((s) => s.type);
            if (!sectionTypes.includes('hero')) {
                throw new Error('Response missing hero section');
            }
            parsed.site.sections.forEach((section, index) => {
                if (typeof section.order !== 'number') {
                    section.order = index + 1;
                }
            });
            parsed.site.sections.sort((a, b) => a.order - b.order);
            if (!parsed.site.theme) {
                parsed.site.theme = {
                    primaryColor: '#3B82F6',
                    secondaryColor: '#EF4444',
                    fontFamily: 'sans',
                    layout: 'modern'
                };
            }
            if (!parsed.site.metadata) {
                parsed.site.metadata = {
                    tone: 'professional',
                    language: 'en'
                };
            }
            return parsed;
        }
        catch (error) {
            logger_1.logger.error('Failed to parse OpenAI response:', error);
            logger_1.logger.error('Raw response content:', content);
            throw new Error(`Invalid JSON response from AI: ${error instanceof Error ? error.message : 'Unknown parsing error'}`);
        }
    }
}
exports.SiteGeneratorService = SiteGeneratorService;
//# sourceMappingURL=SiteGeneratorService.js.map
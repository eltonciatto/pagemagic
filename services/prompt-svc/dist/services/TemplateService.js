"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.TemplateService = void 0;
class TemplateService {
    templates = new Map();
    constructor() {
        this.initializeDefaultTemplates();
    }
    async getTemplates(filters) {
        let templates = Array.from(this.templates.values()).filter(t => t.isActive);
        if (filters?.category) {
            templates = templates.filter(t => t.category === filters.category);
        }
        if (filters?.tags && filters.tags.length > 0) {
            templates = templates.filter(t => filters.tags.some(tag => t.tags.includes(tag)));
        }
        return templates.sort((a, b) => a.name.localeCompare(b.name));
    }
    async getTemplate(templateId) {
        return this.templates.get(templateId) || null;
    }
    async applyTemplate(templateId, variables) {
        const template = this.templates.get(templateId);
        if (!template) {
            throw new Error(`Template ${templateId} not found`);
        }
        let prompt = template.template;
        template.variables.forEach(variable => {
            const value = variables[variable] || '';
            const regex = new RegExp(`{{${variable}}}`, 'g');
            prompt = prompt.replace(regex, value);
        });
        return prompt;
    }
    async createTemplate(template) {
        const id = `template_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
        const now = new Date();
        const newTemplate = {
            ...template,
            id,
            createdAt: now,
            updatedAt: now,
        };
        this.templates.set(id, newTemplate);
        return newTemplate;
    }
    initializeDefaultTemplates() {
        const defaultTemplates = [
            {
                name: 'Landing Page',
                description: 'Create a modern landing page for a business',
                category: 'website',
                template: `Create a modern, responsive landing page for {{businessName}}.

Business Details:
- Industry: {{industry}}
- Target Audience: {{targetAudience}}
- Key Services: {{services}}
- Brand Colors: {{brandColors}}

Include these sections:
1. Hero section with compelling headline
2. About section
3. Services/Features section
4. Testimonials section
5. Contact section
6. Footer

Use modern design principles, mobile-first approach, and include call-to-action buttons.`,
                variables: ['businessName', 'industry', 'targetAudience', 'services', 'brandColors'],
                examples: [
                    'Create a landing page for "Tech Solutions Inc", a software consulting company targeting small businesses',
                ],
                tags: ['landing-page', 'business', 'modern'],
                isActive: true,
            },
            {
                name: 'E-commerce Product Page',
                description: 'Generate a product page for an e-commerce store',
                category: 'ecommerce',
                template: `Create an e-commerce product page for {{productName}}.

Product Details:
- Category: {{category}}
- Price: {{price}}
- Key Features: {{features}}
- Description: {{description}}

Include:
1. Product image gallery section
2. Product title and price
3. Feature highlights
4. Product description
5. Add to cart functionality
6. Related products section
7. Customer reviews section

Use conversion-optimized design and clear call-to-action buttons.`,
                variables: ['productName', 'category', 'price', 'features', 'description'],
                examples: [
                    'Create a product page for "Wireless Headphones" priced at $99',
                ],
                tags: ['ecommerce', 'product-page', 'conversion'],
                isActive: true,
            },
            {
                name: 'Blog Article',
                description: 'Create a blog article layout with content',
                category: 'blog',
                template: `Create a blog article page about {{topic}}.

Article Details:
- Target Audience: {{audience}}
- Tone: {{tone}}
- Word Count: {{wordCount}}
- Key Points: {{keyPoints}}

Include:
1. Article header with title and meta information
2. Introduction paragraph
3. Main content sections with subheadings
4. Related articles sidebar
5. Author bio section
6. Comments section
7. Social sharing buttons

Make it SEO-friendly and engaging for readers.`,
                variables: ['topic', 'audience', 'tone', 'wordCount', 'keyPoints'],
                examples: [
                    'Create a blog article about "Web Development Best Practices" for developers',
                ],
                tags: ['blog', 'content', 'seo'],
                isActive: true,
            },
            {
                name: 'Portfolio Website',
                description: 'Create a personal portfolio website',
                category: 'portfolio',
                template: `Create a personal portfolio website for {{name}}.

Profile Details:
- Profession: {{profession}}
- Skills: {{skills}}
- Experience: {{experience}}
- Projects: {{projects}}

Include:
1. Hero section with personal introduction
2. About section
3. Skills and expertise section
4. Portfolio/Projects showcase
5. Work experience timeline
6. Contact form
7. Social media links

Use a clean, professional design that showcases work effectively.`,
                variables: ['name', 'profession', 'skills', 'experience', 'projects'],
                examples: [
                    'Create a portfolio for "John Doe", a UX Designer with 5 years experience',
                ],
                tags: ['portfolio', 'personal', 'professional'],
                isActive: true,
            },
            {
                name: 'Restaurant Website',
                description: 'Create a restaurant website with menu',
                category: 'restaurant',
                template: `Create a restaurant website for {{restaurantName}}.

Restaurant Details:
- Cuisine Type: {{cuisineType}}
- Location: {{location}}
- Atmosphere: {{atmosphere}}
- Specialties: {{specialties}}

Include:
1. Hero section with restaurant ambiance
2. About the restaurant
3. Menu sections (appetizers, mains, desserts, drinks)
4. Location and hours
5. Reservation system
6. Photo gallery
7. Contact information

Use warm, inviting design with high-quality food imagery placeholders.`,
                variables: ['restaurantName', 'cuisineType', 'location', 'atmosphere', 'specialties'],
                examples: [
                    'Create a website for "Bella Italia", an authentic Italian restaurant',
                ],
                tags: ['restaurant', 'menu', 'local-business'],
                isActive: true,
            },
        ];
        defaultTemplates.forEach(template => {
            this.createTemplate(template);
        });
    }
}
exports.TemplateService = TemplateService;
//# sourceMappingURL=TemplateService.js.map
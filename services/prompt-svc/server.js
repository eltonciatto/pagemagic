const express = require('express');
const cors = require('cors');
const { Pool } = require('pg');
const { connect, StringCodec } = require('nats');
const OpenAI = require('openai');

class PromptService {
    constructor() {
        this.app = express();
        this.db = null;
        this.nats = null;
        this.openai = null;
        this.sc = StringCodec();
        
        this.setupMiddleware();
        this.setupRoutes();
        this.initializeConnections();
    }

    async initializeConnections() {
        // Database connection
        this.db = new Pool({
            connectionString: process.env.DATABASE_URL || 'postgres://pagemagic:pagemagic_dev@localhost:5432/pagemagic'
        });

        // NATS connection
        try {
            this.nats = await connect({ 
                servers: process.env.NATS_URL || 'nats://localhost:4222' 
            });
            console.log('Connected to NATS');
        } catch (error) {
            console.error('Failed to connect to NATS:', error);
        }

        // OpenAI connection
        this.openai = new OpenAI({
            apiKey: process.env.OPENAI_API_KEY
        });

        console.log('Prompt service initialized');
    }

    setupMiddleware() {
        this.app.use(cors({
            origin: process.env.CORS_ORIGIN || '*',
            credentials: true
        }));
        
        this.app.use(express.json({ limit: '10mb' }));
        this.app.use(express.urlencoded({ extended: true }));
        
        // Request logging
        this.app.use((req, res, next) => {
            console.log(`${new Date().toISOString()} - ${req.method} ${req.path}`);
            next();
        });
    }

    setupRoutes() {
        // Health check
        this.app.get('/v1/health', (req, res) => {
            res.json({
                status: 'healthy',
                timestamp: Date.now(),
                service: 'prompt-svc'
            });
        });

        // Generate website from prompt
        this.app.post('/v1/generate', this.generateWebsite.bind(this));

        // Get generation history
        this.app.get('/v1/generations/:user_id', this.getGenerations.bind(this));

        // Get specific generation
        this.app.get('/v1/generations/:user_id/:generation_id', this.getGeneration.bind(this));

        // Regenerate with modifications
        this.app.post('/v1/regenerate', this.regenerateWebsite.bind(this));

        // Chat with AI about website
        this.app.post('/v1/chat', this.chatAboutWebsite.bind(this));
    }

    async generateWebsite(req, res) {
        try {
            const { prompt, user_id, style_preferences, target_audience } = req.body;
            
            if (!prompt || !user_id) {
                return res.status(400).json({ error: 'Prompt and user_id are required' });
            }

            const generationId = `gen_${Date.now()}`;

            // Create system prompt for website generation
            const systemPrompt = `You are an expert web designer and developer. Generate a complete website structure based on the user's prompt. 

Return a JSON object with this structure:
{
  "title": "Website title",
  "description": "Website description",
  "structure": {
    "header": { "type": "header", "content": "..." },
    "sections": [
      { "type": "hero", "content": "...", "style": "..." },
      { "type": "features", "content": "...", "style": "..." }
    ],
    "footer": { "type": "footer", "content": "..." }
  },
  "colors": { "primary": "#hex", "secondary": "#hex" },
  "fonts": { "heading": "font-name", "body": "font-name" }
}`;

            const userPrompt = `Create a website for: ${prompt}
            
${style_preferences ? `Style preferences: ${style_preferences}` : ''}
${target_audience ? `Target audience: ${target_audience}` : ''}

Please generate a modern, responsive website structure.`;

            // Call OpenAI
            const response = await this.openai.chat.completions.create({
                model: 'gpt-4',
                messages: [
                    { role: 'system', content: systemPrompt },
                    { role: 'user', content: userPrompt }
                ],
                temperature: 0.7,
                max_tokens: 4000
            });

            let websiteData;
            try {
                websiteData = JSON.parse(response.choices[0].message.content);
            } catch (error) {
                console.error('Failed to parse AI response:', error);
                return res.status(500).json({ error: 'Failed to parse AI response' });
            }

            // Save generation to database
            const generation = {
                id: generationId,
                user_id,
                prompt,
                website_data: websiteData,
                style_preferences,
                target_audience,
                status: 'completed',
                created_at: new Date()
            };

            await this.db.query(`
                INSERT INTO generations (id, user_id, prompt, website_data, style_preferences, target_audience, status, created_at)
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
            `, [
                generation.id,
                generation.user_id,
                generation.prompt,
                JSON.stringify(generation.website_data),
                generation.style_preferences,
                generation.target_audience,
                generation.status,
                generation.created_at
            ]);

            // Publish to NATS for builder service
            if (this.nats) {
                await this.nats.publish('website.generated', this.sc.encode(JSON.stringify({
                    generation_id: generationId,
                    user_id,
                    website_data: websiteData
                })));
            }

            res.json({
                generation_id: generationId,
                website_data: websiteData,
                status: 'completed'
            });

        } catch (error) {
            console.error('Generation error:', error);
            res.status(500).json({ error: 'Generation failed' });
        }
    }

    async getGenerations(req, res) {
        try {
            const { user_id } = req.params;
            const { limit = 10, offset = 0 } = req.query;

            const result = await this.db.query(`
                SELECT id, prompt, style_preferences, target_audience, status, created_at
                FROM generations 
                WHERE user_id = $1 
                ORDER BY created_at DESC 
                LIMIT $2 OFFSET $3
            `, [user_id, limit, offset]);

            res.json({
                generations: result.rows,
                total: result.rowCount
            });

        } catch (error) {
            console.error('Get generations error:', error);
            res.status(500).json({ error: 'Failed to fetch generations' });
        }
    }

    async getGeneration(req, res) {
        try {
            const { user_id, generation_id } = req.params;

            const result = await this.db.query(`
                SELECT * FROM generations 
                WHERE id = $1 AND user_id = $2
            `, [generation_id, user_id]);

            if (result.rows.length === 0) {
                return res.status(404).json({ error: 'Generation not found' });
            }

            const generation = result.rows[0];
            generation.website_data = JSON.parse(generation.website_data);

            res.json(generation);

        } catch (error) {
            console.error('Get generation error:', error);
            res.status(500).json({ error: 'Failed to fetch generation' });
        }
    }

    async regenerateWebsite(req, res) {
        try {
            const { generation_id, modifications, user_id } = req.body;
            
            if (!generation_id || !modifications || !user_id) {
                return res.status(400).json({ error: 'generation_id, modifications, and user_id are required' });
            }

            // Get original generation
            const originalResult = await this.db.query(`
                SELECT * FROM generations WHERE id = $1 AND user_id = $2
            `, [generation_id, user_id]);

            if (originalResult.rows.length === 0) {
                return res.status(404).json({ error: 'Original generation not found' });
            }

            const original = originalResult.rows[0];
            const originalData = JSON.parse(original.website_data);

            const newGenerationId = `gen_${Date.now()}`;

            // Create modification prompt
            const modificationPrompt = `Based on this existing website structure:
${JSON.stringify(originalData, null, 2)}

Apply these modifications: ${modifications}

Return the complete updated website structure in the same JSON format.`;

            const response = await this.openai.chat.completions.create({
                model: 'gpt-4',
                messages: [
                    { role: 'user', content: modificationPrompt }
                ],
                temperature: 0.7,
                max_tokens: 4000
            });

            let updatedWebsiteData;
            try {
                updatedWebsiteData = JSON.parse(response.choices[0].message.content);
            } catch (error) {
                console.error('Failed to parse AI response:', error);
                return res.status(500).json({ error: 'Failed to parse AI response' });
            }

            // Save new generation
            await this.db.query(`
                INSERT INTO generations (id, user_id, prompt, website_data, style_preferences, target_audience, status, created_at, parent_id)
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
            `, [
                newGenerationId,
                user_id,
                `${original.prompt} [Modified: ${modifications}]`,
                JSON.stringify(updatedWebsiteData),
                original.style_preferences,
                original.target_audience,
                'completed',
                new Date(),
                generation_id
            ]);

            res.json({
                generation_id: newGenerationId,
                website_data: updatedWebsiteData,
                status: 'completed'
            });

        } catch (error) {
            console.error('Regeneration error:', error);
            res.status(500).json({ error: 'Regeneration failed' });
        }
    }

    async chatAboutWebsite(req, res) {
        try {
            const { message, generation_id, user_id, chat_history = [] } = req.body;
            
            if (!message || !generation_id || !user_id) {
                return res.status(400).json({ error: 'message, generation_id, and user_id are required' });
            }

            // Get website data for context
            const result = await this.db.query(`
                SELECT website_data FROM generations WHERE id = $1 AND user_id = $2
            `, [generation_id, user_id]);

            if (result.rows.length === 0) {
                return res.status(404).json({ error: 'Generation not found' });
            }

            const websiteData = JSON.parse(result.rows[0].website_data);

            // Create chat context
            const systemPrompt = `You are a helpful assistant that can answer questions about this website:
${JSON.stringify(websiteData, null, 2)}

Help the user understand the website structure, suggest improvements, or answer any questions about the design.`;

            const messages = [
                { role: 'system', content: systemPrompt },
                ...chat_history,
                { role: 'user', content: message }
            ];

            const response = await this.openai.chat.completions.create({
                model: 'gpt-3.5-turbo',
                messages: messages,
                temperature: 0.7,
                max_tokens: 1000
            });

            const reply = response.choices[0].message.content;

            res.json({
                reply,
                updated_history: [...chat_history, 
                    { role: 'user', content: message },
                    { role: 'assistant', content: reply }
                ]
            });

        } catch (error) {
            console.error('Chat error:', error);
            res.status(500).json({ error: 'Chat failed' });
        }
    }

    start() {
        const port = process.env.PORT || 3000;
        this.app.listen(port, () => {
            console.log(`Prompt service running on port ${port}`);
        });
    }
}

// Start the service
const service = new PromptService();
service.start();

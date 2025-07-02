const express = require('express');
const cors = require('cors');
const { createClient } = require('redis');
const fs = require('fs').promises;
const path = require('path');

class I18nService {
    constructor() {
        this.app = express();
        this.translations = new Map();
        this.fallbackLanguage = 'en';
        this.supportedLanguages = ['en', 'pt', 'es', 'fr', 'de', 'it', 'zh', 'ja', 'ko', 'ar'];
        this.redis = null;
        
        this.setupMiddleware();
        this.setupRoutes();
        this.initializeRedis();
        this.loadTranslations();
    }

    async initializeRedis() {
        try {
            this.redis = createClient({
                url: process.env.REDIS_URL || 'redis://localhost:6379',
                password: process.env.REDIS_PASSWORD || undefined
            });
            
            await this.redis.connect();
            console.log('Connected to Redis');
        } catch (error) {
            console.error('Failed to connect to Redis:', error);
        }
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
                service: 'i18n-svc',
                supported_languages: this.supportedLanguages
            });
        });

        // Get translations for a specific language
        this.app.get('/v1/translations/:language', this.getTranslations.bind(this));

        // Get specific translation key
        this.app.get('/v1/translations/:language/:namespace', this.getNamespaceTranslations.bind(this));

        // Translate text using AI (for dynamic content)
        this.app.post('/v1/translate', this.translateText.bind(this));

        // Add/update translation
        this.app.post('/v1/translations/:language/:namespace', this.updateTranslations.bind(this));

        // Delete translation
        this.app.delete('/v1/translations/:language/:namespace/:key', this.deleteTranslation.bind(this));

        // Get all supported languages
        this.app.get('/v1/languages', this.getSupportedLanguages.bind(this));

        // Auto-detect language from text
        this.app.post('/v1/detect', this.detectLanguage.bind(this));

        // Bulk translate namespace
        this.app.post('/v1/bulk-translate', this.bulkTranslate.bind(this));

        // Export translations
        this.app.get('/v1/export/:language', this.exportTranslations.bind(this));

        // Import translations
        this.app.post('/v1/import/:language', this.importTranslations.bind(this));
    }

    async loadTranslations() {
        // Load default translations from files
        const translationsDir = path.join(__dirname, 'translations');
        
        try {
            await fs.access(translationsDir);
        } catch {
            await fs.mkdir(translationsDir, { recursive: true });
        }

        for (const lang of this.supportedLanguages) {
            try {
                const langDir = path.join(translationsDir, lang);
                await fs.access(langDir);
                
                const files = await fs.readdir(langDir);
                const langTranslations = {};
                
                for (const file of files) {
                    if (file.endsWith('.json')) {
                        const namespace = file.replace('.json', '');
                        const content = await fs.readFile(path.join(langDir, file), 'utf8');
                        langTranslations[namespace] = JSON.parse(content);
                    }
                }
                
                this.translations.set(lang, langTranslations);
                
                // Cache in Redis
                if (this.redis) {
                    await this.redis.setEx(
                        `i18n:${lang}`,
                        3600,
                        JSON.stringify(langTranslations)
                    );
                }
            } catch (error) {
                console.log(`No translations found for language: ${lang}`);
                this.translations.set(lang, {});
            }
        }

        // Load default English translations if not exist
        if (!this.translations.get('en') || Object.keys(this.translations.get('en')).length === 0) {
            await this.createDefaultTranslations();
        }
    }

    async createDefaultTranslations() {
        const defaultTranslations = {
            common: {
                "hello": "Hello",
                "goodbye": "Goodbye",
                "yes": "Yes",
                "no": "No",
                "save": "Save",
                "cancel": "Cancel",
                "delete": "Delete",
                "edit": "Edit",
                "create": "Create",
                "loading": "Loading...",
                "error": "Error",
                "success": "Success",
                "warning": "Warning",
                "info": "Information"
            },
            auth: {
                "login": "Login",
                "logout": "Logout",
                "register": "Register",
                "forgot_password": "Forgot Password",
                "reset_password": "Reset Password",
                "email": "Email",
                "password": "Password",
                "confirm_password": "Confirm Password",
                "invalid_credentials": "Invalid credentials",
                "account_created": "Account created successfully"
            },
            dashboard: {
                "welcome": "Welcome",
                "my_sites": "My Sites",
                "create_site": "Create Site",
                "usage": "Usage",
                "billing": "Billing",
                "settings": "Settings",
                "analytics": "Analytics"
            },
            editor: {
                "add_element": "Add Element",
                "text": "Text",
                "image": "Image",
                "button": "Button",
                "container": "Container",
                "preview": "Preview",
                "publish": "Publish",
                "draft": "Draft",
                "undo": "Undo",
                "redo": "Redo"
            }
        };

        this.translations.set('en', defaultTranslations);
        
        // Save to file
        const enDir = path.join(__dirname, 'translations', 'en');
        await fs.mkdir(enDir, { recursive: true });
        
        for (const [namespace, translations] of Object.entries(defaultTranslations)) {
            await fs.writeFile(
                path.join(enDir, `${namespace}.json`),
                JSON.stringify(translations, null, 2)
            );
        }

        // Cache in Redis
        if (this.redis) {
            await this.redis.setEx(
                'i18n:en',
                3600,
                JSON.stringify(defaultTranslations)
            );
        }
    }

    async getTranslations(req, res) {
        try {
            const { language } = req.params;
            
            if (!this.supportedLanguages.includes(language)) {
                return res.status(400).json({
                    error: 'Unsupported language',
                    supported: this.supportedLanguages
                });
            }

            // Try Redis first
            let translations = null;
            if (this.redis) {
                const cached = await this.redis.get(`i18n:${language}`);
                if (cached) {
                    translations = JSON.parse(cached);
                }
            }

            // Fallback to memory
            if (!translations) {
                translations = this.translations.get(language) || {};
            }

            // If empty, fallback to English
            if (Object.keys(translations).length === 0 && language !== 'en') {
                translations = this.translations.get('en') || {};
            }

            res.json({
                language,
                translations,
                fallback_used: Object.keys(this.translations.get(language) || {}).length === 0
            });
        } catch (error) {
            console.error('Get translations error:', error);
            res.status(500).json({ error: 'Internal server error' });
        }
    }

    async getNamespaceTranslations(req, res) {
        try {
            const { language, namespace } = req.params;
            
            if (!this.supportedLanguages.includes(language)) {
                return res.status(400).json({
                    error: 'Unsupported language',
                    supported: this.supportedLanguages
                });
            }

            const langTranslations = this.translations.get(language) || {};
            const namespaceTranslations = langTranslations[namespace] || {};

            // Fallback to English if not found
            if (Object.keys(namespaceTranslations).length === 0 && language !== 'en') {
                const enTranslations = this.translations.get('en') || {};
                const enNamespace = enTranslations[namespace] || {};
                return res.json({
                    language,
                    namespace,
                    translations: enNamespace,
                    fallback_used: true
                });
            }

            res.json({
                language,
                namespace,
                translations: namespaceTranslations,
                fallback_used: false
            });
        } catch (error) {
            console.error('Get namespace translations error:', error);
            res.status(500).json({ error: 'Internal server error' });
        }
    }

    async translateText(req, res) {
        try {
            const { text, from, to, context } = req.body;
            
            if (!text || !to) {
                return res.status(400).json({ error: 'Text and target language required' });
            }

            if (!this.supportedLanguages.includes(to)) {
                return res.status(400).json({
                    error: 'Unsupported target language',
                    supported: this.supportedLanguages
                });
            }

            // Simple translation service (in production, use Google Translate, DeepL, etc.)
            const translatedText = await this.performTranslation(text, from || 'auto', to, context);
            
            res.json({
                original_text: text,
                translated_text: translatedText,
                from: from || 'auto',
                to,
                context
            });
        } catch (error) {
            console.error('Translation error:', error);
            res.status(500).json({ error: 'Translation failed' });
        }
    }

    async performTranslation(text, from, to, context) {
        // This is a placeholder - in production you'd integrate with:
        // - Google Translate API
        // - DeepL API
        // - Azure Translator
        // - OpenAI for context-aware translation
        
        // For now, return the original text with a note
        if (to === 'en') {
            return text;
        }
        
        // Simple mock translations for common phrases
        const mockTranslations = {
            pt: {
                "Hello": "Olá",
                "Welcome": "Bem-vindo",
                "Login": "Entrar",
                "Create Site": "Criar Site",
                "Save": "Salvar",
                "Cancel": "Cancelar"
            },
            es: {
                "Hello": "Hola",
                "Welcome": "Bienvenido",
                "Login": "Iniciar sesión",
                "Create Site": "Crear Sitio",
                "Save": "Guardar",
                "Cancel": "Cancelar"
            },
            fr: {
                "Hello": "Bonjour",
                "Welcome": "Bienvenue",
                "Login": "Connexion",
                "Create Site": "Créer un Site",
                "Save": "Enregistrer",
                "Cancel": "Annuler"
            }
        };
        
        return mockTranslations[to]?.[text] || `[${to.toUpperCase()}] ${text}`;
    }

    async updateTranslations(req, res) {
        try {
            const { language, namespace } = req.params;
            const { translations } = req.body;
            
            if (!this.supportedLanguages.includes(language)) {
                return res.status(400).json({
                    error: 'Unsupported language',
                    supported: this.supportedLanguages
                });
            }

            if (!translations || typeof translations !== 'object') {
                return res.status(400).json({ error: 'Invalid translations object' });
            }

            // Update memory
            const langTranslations = this.translations.get(language) || {};
            langTranslations[namespace] = { ...langTranslations[namespace], ...translations };
            this.translations.set(language, langTranslations);

            // Save to file
            const langDir = path.join(__dirname, 'translations', language);
            await fs.mkdir(langDir, { recursive: true });
            
            await fs.writeFile(
                path.join(langDir, `${namespace}.json`),
                JSON.stringify(langTranslations[namespace], null, 2)
            );

            // Update Redis cache
            if (this.redis) {
                await this.redis.setEx(
                    `i18n:${language}`,
                    3600,
                    JSON.stringify(langTranslations)
                );
            }

            res.json({
                status: 'success',
                language,
                namespace,
                updated_keys: Object.keys(translations)
            });
        } catch (error) {
            console.error('Update translations error:', error);
            res.status(500).json({ error: 'Failed to update translations' });
        }
    }

    async deleteTranslation(req, res) {
        try {
            const { language, namespace, key } = req.params;
            
            const langTranslations = this.translations.get(language) || {};
            const namespaceTranslations = langTranslations[namespace] || {};
            
            if (namespaceTranslations[key]) {
                delete namespaceTranslations[key];
                this.translations.set(language, langTranslations);
                
                // Save to file
                const langDir = path.join(__dirname, 'translations', language);
                await fs.writeFile(
                    path.join(langDir, `${namespace}.json`),
                    JSON.stringify(namespaceTranslations, null, 2)
                );

                // Update Redis cache
                if (this.redis) {
                    await this.redis.setEx(
                        `i18n:${language}`,
                        3600,
                        JSON.stringify(langTranslations)
                    );
                }

                res.json({ status: 'deleted', key });
            } else {
                res.status(404).json({ error: 'Translation key not found' });
            }
        } catch (error) {
            console.error('Delete translation error:', error);
            res.status(500).json({ error: 'Failed to delete translation' });
        }
    }

    async getSupportedLanguages(req, res) {
        res.json({
            languages: this.supportedLanguages,
            default: this.fallbackLanguage
        });
    }

    async detectLanguage(req, res) {
        try {
            const { text } = req.body;
            
            if (!text) {
                return res.status(400).json({ error: 'Text required' });
            }

            // Simple language detection (in production, use proper detection service)
            const detectedLanguage = this.simpleLanguageDetection(text);
            
            res.json({
                text,
                detected_language: detectedLanguage,
                confidence: 0.8 // Mock confidence
            });
        } catch (error) {
            console.error('Language detection error:', error);
            res.status(500).json({ error: 'Language detection failed' });
        }
    }

    simpleLanguageDetection(text) {
        // Very simple detection based on common words
        const patterns = {
            en: /\b(the|and|or|but|in|on|at|to|for|of|with|by)\b/gi,
            pt: /\b(o|a|e|ou|mas|em|no|na|para|de|com|por)\b/gi,
            es: /\b(el|la|y|o|pero|en|del|para|de|con|por)\b/gi,
            fr: /\b(le|la|et|ou|mais|dans|sur|pour|de|avec|par)\b/gi
        };

        let maxMatches = 0;
        let detectedLang = 'en';

        for (const [lang, pattern] of Object.entries(patterns)) {
            const matches = (text.match(pattern) || []).length;
            if (matches > maxMatches) {
                maxMatches = matches;
                detectedLang = lang;
            }
        }

        return detectedLang;
    }

    async bulkTranslate(req, res) {
        try {
            const { namespace, from_language, to_languages } = req.body;
            
            if (!namespace || !from_language || !to_languages || !Array.isArray(to_languages)) {
                return res.status(400).json({ 
                    error: 'namespace, from_language, and to_languages array required' 
                });
            }

            const sourceTranslations = this.translations.get(from_language)?.[namespace] || {};
            
            if (Object.keys(sourceTranslations).length === 0) {
                return res.status(404).json({ error: 'Source namespace not found' });
            }

            const results = {};
            
            for (const targetLang of to_languages) {
                if (!this.supportedLanguages.includes(targetLang)) continue;
                
                const translatedNamespace = {};
                
                for (const [key, text] of Object.entries(sourceTranslations)) {
                    translatedNamespace[key] = await this.performTranslation(
                        text, 
                        from_language, 
                        targetLang, 
                        namespace
                    );
                }
                
                // Update translations
                const langTranslations = this.translations.get(targetLang) || {};
                langTranslations[namespace] = translatedNamespace;
                this.translations.set(targetLang, langTranslations);
                
                results[targetLang] = translatedNamespace;
            }

            res.json({
                status: 'success',
                namespace,
                from_language,
                translated_to: Object.keys(results),
                results
            });
        } catch (error) {
            console.error('Bulk translate error:', error);
            res.status(500).json({ error: 'Bulk translation failed' });
        }
    }

    async exportTranslations(req, res) {
        try {
            const { language } = req.params;
            const { format = 'json' } = req.query;
            
            if (!this.supportedLanguages.includes(language)) {
                return res.status(400).json({
                    error: 'Unsupported language',
                    supported: this.supportedLanguages
                });
            }

            const translations = this.translations.get(language) || {};
            
            if (format === 'json') {
                res.setHeader('Content-Type', 'application/json');
                res.setHeader('Content-Disposition', `attachment; filename="${language}.json"`);
                res.json(translations);
            } else {
                res.status(400).json({ error: 'Unsupported format' });
            }
        } catch (error) {
            console.error('Export error:', error);
            res.status(500).json({ error: 'Export failed' });
        }
    }

    async importTranslations(req, res) {
        try {
            const { language } = req.params;
            const { translations, merge = true } = req.body;
            
            if (!this.supportedLanguages.includes(language)) {
                return res.status(400).json({
                    error: 'Unsupported language',
                    supported: this.supportedLanguages
                });
            }

            if (!translations || typeof translations !== 'object') {
                return res.status(400).json({ error: 'Invalid translations object' });
            }

            let langTranslations = {};
            
            if (merge) {
                langTranslations = this.translations.get(language) || {};
            }

            // Merge or replace translations
            for (const [namespace, namespaceTranslations] of Object.entries(translations)) {
                if (merge && langTranslations[namespace]) {
                    langTranslations[namespace] = { 
                        ...langTranslations[namespace], 
                        ...namespaceTranslations 
                    };
                } else {
                    langTranslations[namespace] = namespaceTranslations;
                }
            }

            this.translations.set(language, langTranslations);

            // Save to files
            const langDir = path.join(__dirname, 'translations', language);
            await fs.mkdir(langDir, { recursive: true });
            
            for (const [namespace, namespaceTranslations] of Object.entries(langTranslations)) {
                await fs.writeFile(
                    path.join(langDir, `${namespace}.json`),
                    JSON.stringify(namespaceTranslations, null, 2)
                );
            }

            // Update Redis cache
            if (this.redis) {
                await this.redis.setEx(
                    `i18n:${language}`,
                    3600,
                    JSON.stringify(langTranslations)
                );
            }

            res.json({
                status: 'success',
                language,
                imported_namespaces: Object.keys(translations),
                merge_mode: merge
            });
        } catch (error) {
            console.error('Import error:', error);
            res.status(500).json({ error: 'Import failed' });
        }
    }

    start() {
        const port = process.env.PORT || 8086;
        this.app.listen(port, () => {
            console.log(`I18n service running on port ${port}`);
        });
    }
}

// Start the service
const service = new I18nService();
service.start();

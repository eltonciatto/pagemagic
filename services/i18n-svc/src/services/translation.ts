import { Pool } from 'pg';
import { Redis } from 'redis';
import { 
  Translation, 
  TranslationBundle, 
  Project, 
  TranslationRequest,
  BulkTranslationRequest,
  TranslationStats 
} from '../models/translation';
import { logger } from '../utils/logger';

export class TranslationService {
  private db: Pool;
  private redis: Redis;

  constructor(db: Pool, redis: Redis) {
    this.db = db;
    this.redis = redis;
  }

  async createProject(name: string, defaultLocale: string, supportedLocales: string[]): Promise<Project> {
    const client = await this.db.connect();
    try {
      const query = `
        INSERT INTO i18n_projects (id, name, default_locale, supported_locales, auto_translate, created_at, updated_at)
        VALUES (gen_random_uuid(), $1, $2, $3, false, NOW(), NOW())
        RETURNING *
      `;
      
      const result = await client.query(query, [name, defaultLocale, JSON.stringify(supportedLocales)]);
      const project = result.rows[0];
      
      return {
        id: project.id,
        name: project.name,
        default_locale: project.default_locale,
        supported_locales: JSON.parse(project.supported_locales),
        auto_translate: project.auto_translate,
        created_at: project.created_at,
        updated_at: project.updated_at
      };
    } finally {
      client.release();
    }
  }

  async getProject(projectId: string): Promise<Project | null> {
    const client = await this.db.connect();
    try {
      const query = 'SELECT * FROM i18n_projects WHERE id = $1';
      const result = await client.query(query, [projectId]);
      
      if (result.rows.length === 0) {
        return null;
      }
      
      const project = result.rows[0];
      return {
        id: project.id,
        name: project.name,
        default_locale: project.default_locale,
        supported_locales: JSON.parse(project.supported_locales),
        auto_translate: project.auto_translate,
        created_at: project.created_at,
        updated_at: project.updated_at
      };
    } finally {
      client.release();
    }
  }

  async createTranslation(request: TranslationRequest): Promise<Translation> {
    const client = await this.db.connect();
    try {
      // Check if translation already exists
      const existingQuery = 'SELECT id FROM i18n_translations WHERE project_id = $1 AND locale = $2 AND key = $3';
      const existingResult = await client.query(existingQuery, [request.project_id, request.locale, request.key]);
      
      if (existingResult.rows.length > 0) {
        // Update existing translation
        const updateQuery = `
          UPDATE i18n_translations 
          SET value = $1, context = $2, updated_at = NOW()
          WHERE project_id = $3 AND locale = $4 AND key = $5
          RETURNING *
        `;
        const result = await client.query(updateQuery, [
          request.value, request.context, request.project_id, request.locale, request.key
        ]);
        
        const translation = result.rows[0];
        await this.invalidateCache(request.project_id, request.locale);
        
        return {
          id: translation.id,
          project_id: translation.project_id,
          locale: translation.locale,
          key: translation.key,
          value: translation.value,
          context: translation.context,
          created_at: translation.created_at,
          updated_at: translation.updated_at
        };
      } else {
        // Create new translation
        const insertQuery = `
          INSERT INTO i18n_translations (id, project_id, locale, key, value, context, created_at, updated_at)
          VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, NOW(), NOW())
          RETURNING *
        `;
        const result = await client.query(insertQuery, [
          request.project_id, request.locale, request.key, request.value, request.context
        ]);
        
        const translation = result.rows[0];
        await this.invalidateCache(request.project_id, request.locale);
        
        return {
          id: translation.id,
          project_id: translation.project_id,
          locale: translation.locale,
          key: translation.key,
          value: translation.value,
          context: translation.context,
          created_at: translation.created_at,
          updated_at: translation.updated_at
        };
      }
    } finally {
      client.release();
    }
  }

  async bulkCreateTranslations(request: BulkTranslationRequest): Promise<Translation[]> {
    const client = await this.db.connect();
    try {
      await client.query('BEGIN');
      
      const translations: Translation[] = [];
      
      for (const translation of request.translations) {
        const translationRequest: TranslationRequest = {
          project_id: request.project_id,
          locale: request.locale,
          key: translation.key,
          value: translation.value,
          context: translation.context
        };
        
        const result = await this.createTranslation(translationRequest);
        translations.push(result);
      }
      
      await client.query('COMMIT');
      await this.invalidateCache(request.project_id, request.locale);
      
      return translations;
    } catch (error) {
      await client.query('ROLLBACK');
      throw error;
    } finally {
      client.release();
    }
  }

  async getTranslations(projectId: string, locale: string): Promise<Record<string, string>> {
    // Try cache first
    const cacheKey = `i18n:${projectId}:${locale}`;
    try {
      const cached = await this.redis.get(cacheKey);
      if (cached) {
        return JSON.parse(cached);
      }
    } catch (error) {
      logger.warn('Redis cache miss for translations', { projectId, locale, error });
    }

    // Fetch from database
    const client = await this.db.connect();
    try {
      const query = 'SELECT key, value FROM i18n_translations WHERE project_id = $1 AND locale = $2';
      const result = await client.query(query, [projectId, locale]);
      
      const translations: Record<string, string> = {};
      for (const row of result.rows) {
        translations[row.key] = row.value;
      }
      
      // Cache for 1 hour
      try {
        await this.redis.setex(cacheKey, 3600, JSON.stringify(translations));
      } catch (error) {
        logger.warn('Failed to cache translations', { projectId, locale, error });
      }
      
      return translations;
    } finally {
      client.release();
    }
  }

  async publishBundle(projectId: string, locale: string, version?: string): Promise<TranslationBundle> {
    const translations = await this.getTranslations(projectId, locale);
    const bundleVersion = version || new Date().toISOString();
    
    const client = await this.db.connect();
    try {
      // Mark previous bundles as unpublished
      await client.query(
        'UPDATE i18n_bundles SET published = false WHERE project_id = $1 AND locale = $2',
        [projectId, locale]
      );
      
      // Create new bundle
      const query = `
        INSERT INTO i18n_bundles (id, project_id, locale, version, bundle, published, created_at, updated_at)
        VALUES (gen_random_uuid(), $1, $2, $3, $4, true, NOW(), NOW())
        RETURNING *
      `;
      
      const result = await client.query(query, [
        projectId, locale, bundleVersion, JSON.stringify(translations)
      ]);
      
      const bundle = result.rows[0];
      
      return {
        id: bundle.id,
        project_id: bundle.project_id,
        locale: bundle.locale,
        version: bundle.version,
        bundle: JSON.parse(bundle.bundle),
        published: bundle.published,
        created_at: bundle.created_at,
        updated_at: bundle.updated_at
      };
    } finally {
      client.release();
    }
  }

  async getPublishedBundle(projectId: string, locale: string): Promise<TranslationBundle | null> {
    const client = await this.db.connect();
    try {
      const query = `
        SELECT * FROM i18n_bundles 
        WHERE project_id = $1 AND locale = $2 AND published = true
        ORDER BY created_at DESC
        LIMIT 1
      `;
      
      const result = await client.query(query, [projectId, locale]);
      
      if (result.rows.length === 0) {
        return null;
      }
      
      const bundle = result.rows[0];
      return {
        id: bundle.id,
        project_id: bundle.project_id,
        locale: bundle.locale,
        version: bundle.version,
        bundle: JSON.parse(bundle.bundle),
        published: bundle.published,
        created_at: bundle.created_at,
        updated_at: bundle.updated_at
      };
    } finally {
      client.release();
    }
  }

  async getTranslationStats(projectId: string, locale: string): Promise<TranslationStats> {
    const client = await this.db.connect();
    try {
      // Get all unique keys for the project
      const allKeysQuery = 'SELECT DISTINCT key FROM i18n_translations WHERE project_id = $1';
      const allKeysResult = await client.query(allKeysQuery, [projectId]);
      const totalKeys = allKeysResult.rows.length;
      
      // Get translated keys for the specific locale
      const translatedKeysQuery = 'SELECT COUNT(*) FROM i18n_translations WHERE project_id = $1 AND locale = $2';
      const translatedKeysResult = await client.query(translatedKeysQuery, [projectId, locale]);
      const translatedKeys = parseInt(translatedKeysResult.rows[0].count);
      
      const missingKeys = totalKeys - translatedKeys;
      const completionPercentage = totalKeys > 0 ? (translatedKeys / totalKeys) * 100 : 0;
      
      return {
        project_id: projectId,
        locale,
        total_keys: totalKeys,
        translated_keys: translatedKeys,
        missing_keys: missingKeys,
        completion_percentage: Math.round(completionPercentage * 100) / 100
      };
    } finally {
      client.release();
    }
  }

  async deleteTranslation(projectId: string, locale: string, key: string): Promise<boolean> {
    const client = await this.db.connect();
    try {
      const query = 'DELETE FROM i18n_translations WHERE project_id = $1 AND locale = $2 AND key = $3';
      const result = await client.query(query, [projectId, locale, key]);
      
      if (result.rowCount && result.rowCount > 0) {
        await this.invalidateCache(projectId, locale);
        return true;
      }
      
      return false;
    } finally {
      client.release();
    }
  }

  private async invalidateCache(projectId: string, locale: string): Promise<void> {
    try {
      const cacheKey = `i18n:${projectId}:${locale}`;
      await this.redis.del(cacheKey);
    } catch (error) {
      logger.warn('Failed to invalidate cache', { projectId, locale, error });
    }
  }
}

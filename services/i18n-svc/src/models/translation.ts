export interface Translation {
  id: string;
  project_id: string;
  locale: string;
  key: string;
  value: string;
  context?: string;
  created_at: Date;
  updated_at: Date;
}

export interface TranslationBundle {
  id: string;
  project_id: string;
  locale: string;
  version: string;
  bundle: Record<string, string>;
  published: boolean;
  created_at: Date;
  updated_at: Date;
}

export interface Project {
  id: string;
  name: string;
  default_locale: string;
  supported_locales: string[];
  auto_translate: boolean;
  created_at: Date;
  updated_at: Date;
}

export interface TranslationRequest {
  project_id: string;
  locale: string;
  key: string;
  value: string;
  context?: string;
}

export interface BulkTranslationRequest {
  project_id: string;
  locale: string;
  translations: Array<{
    key: string;
    value: string;
    context?: string;
  }>;
}

export interface TranslationImport {
  project_id: string;
  locale: string;
  format: 'json' | 'csv' | 'yaml';
  data: string | Buffer;
}

export interface AutoTranslationRequest {
  project_id: string;
  source_locale: string;
  target_locales: string[];
  keys?: string[];
}

export interface TranslationStats {
  project_id: string;
  locale: string;
  total_keys: number;
  translated_keys: number;
  missing_keys: number;
  completion_percentage: number;
}

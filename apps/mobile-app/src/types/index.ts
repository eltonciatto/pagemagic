export interface User {
  id: string;
  email: string;
  name?: string;
  avatar?: string;
  subscription_status: 'free' | 'pro' | 'enterprise';
  created_at: string;
}

export interface AuthState {
  isAuthenticated: boolean;
  user: User | null;
  token: string | null;
}

export interface Project {
  id: string;
  name: string;
  description?: string;
  thumbnail?: string;
  template_id?: string;
  status: 'draft' | 'published' | 'building';
  domain?: string;
  created_at: string;
  updated_at: string;
  last_build_at?: string;
}

export interface Site {
  id: string;
  project_id: string;
  name: string;
  domain: string;
  status: 'active' | 'inactive' | 'building' | 'error';
  ssl_enabled: boolean;
  analytics_enabled: boolean;
  created_at: string;
  updated_at: string;
}

export interface Template {
  id: string;
  name: string;
  description: string;
  category: string;
  tags: string[];
  thumbnail: string;
  preview_url: string;
  is_premium: boolean;
  price?: number;
}

export interface GenerationRequest {
  prompt: string;
  template_id?: string;
  locale: string;
  style_preferences?: {
    color_scheme?: string;
    font_family?: string;
    layout_style?: string;
  };
}

export interface GenerationResponse {
  id: string;
  status: 'queued' | 'processing' | 'completed' | 'failed';
  preview_url?: string;
  sections?: PageSection[];
  estimated_completion?: string;
}

export interface PageSection {
  id: string;
  type: 'hero' | 'features' | 'testimonials' | 'pricing' | 'contact' | 'footer';
  title: string;
  content: any;
  order: number;
  editable: boolean;
}

export interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
  message?: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  limit: number;
  has_next: boolean;
  has_prev: boolean;
}

export interface Usage {
  period: string;
  page_generations: number;
  ai_tokens: number;
  storage_gb: number;
  bandwidth_gb: number;
  limits: {
    page_generations: number;
    ai_tokens: number;
    storage_gb: number;
    bandwidth_gb: number;
  };
}

export interface Notification {
  id: string;
  title: string;
  body: string;
  type: 'info' | 'success' | 'warning' | 'error';
  data?: any;
  read: boolean;
  created_at: string;
}

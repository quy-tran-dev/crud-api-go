import type { User } from './auth';

export interface TemplateCategoryItem {
  name: string;
  order_index: number;
  description?: string;
}

export interface TemplateCategory {
  code: string;
  name: string;
  description?: string;
  items: TemplateCategoryItem[];
}

export interface FrameworkTemplate {
  id: number;
  user_id: number;
  author?: User;
  name: string;
  genre: string;
  description: string;
  is_public: boolean;
  usage_count: number;
  data: TemplateCategory[];
  created_at: string;
  updated_at: string;
}

export interface CreateFrameworkRequest {
  name: string;
  genre: string;
  description?: string;
  is_public: boolean;
  data: TemplateCategory[];
}

export interface SaveFrameworkFromStoryRequest {
  name: string;
  description?: string;
  is_public: boolean;
}

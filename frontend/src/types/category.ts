export interface SystemCategory {
  id: number;
  story_id: number;
  code: string;
  name: string;
  description: string;
  items?: SystemCategoryItem[];
  created_at: string;
  updated_at: string;
}

export interface SystemCategoryItem {
  id: number;
  category_id: number;
  code?: string;
  name: string;
  order_index: number;
  description?: string;
  metadata?: Record<string, any>;
  created_at: string;
  updated_at: string;
}

export interface CategoryAttributeValue {
  category_code: string;
  value: string;
  note?: string;
}

export interface DetectedNewItem {
  category_code: string;
  new_value: string;
  message: string;
}

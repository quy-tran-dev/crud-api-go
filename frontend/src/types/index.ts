// ============================================================================
// AUTH & USER TYPES
// ============================================================================

export interface User {
  id: number;
  username: string;
  email: string;
  role: string;
  created_at: string;
}

export interface AuthResponse {
  user: User;
  token: string;
}

// ============================================================================
// WORLD & STORY TYPES
// ============================================================================

export interface Story {
  id: number;
  user_id: number;
  title: string;
  genre: string;
  synopsis: string;
  world_setting: string;
  status: 'draft' | 'official';
  is_draft: boolean;
  categories?: SystemCategory[];
  locations?: Location[];
  factions?: Faction[];
  characters?: Character[];
  plots?: Plot[];
  items?: Item[];
  created_at: string;
  updated_at: string;
}

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

export interface Location {
  id: number;
  story_id: number;
  parent_id?: number | null;
  name: string;
  type: string;
  danger_level?: string;
  description?: string;
  children?: Location[];
  created_at: string;
  updated_at: string;
}

export interface DiplomacyRecord {
  target_faction_id: number;
  target_faction_name?: string;
  status: 'ally' | 'enemy' | 'neutral';
  reason?: string;
}

export interface Faction {
  id: number;
  story_id: number;
  name: string;
  type: string;
  headquarter_id?: number | null;
  leader_id?: number | null;
  power_tier?: string;
  description?: string;
  diplomacy?: DiplomacyRecord[];
  created_at: string;
  updated_at: string;
}

export interface RelationshipItem {
  target_character_id: number;
  target_character_name?: string;
  relation_type: string;
  description?: string;
}

export interface InventoryRecord {
  item_id: number;
  item_name?: string;
  acquired_at_plot_id?: number | null;
  note?: string;
}

export interface Character {
  id: number;
  story_id: number;
  name: string;
  alias?: string;
  role?: string;
  gender?: string;
  age?: string;
  faction_id?: number | null;
  faction_rank?: string;
  birthplace_id?: number | null;
  current_level?: string;
  abilities?: string;
  origin_story?: string;
  character_arc?: string;
  relationships?: RelationshipItem[];
  inventory?: InventoryRecord[];
  category_attrs?: CategoryAttributeValue[];
  is_draft: boolean;
  created_at: string;
  updated_at: string;
}

export interface Plot {
  id: number;
  story_id: number;
  title: string;
  type: string;
  order_index: number;
  location_id?: number | null;
  summary?: string;
  involved_faction_ids?: number[];
  participating_character_ids?: number[];
  featured_item_ids?: number[];
  consequences?: string;
  created_at: string;
  updated_at: string;
}

export interface Item {
  id: number;
  story_id: number;
  name: string;
  type: string;
  grade?: string;
  function?: string;
  usage_condition?: string;
  origin_location_id?: number | null;
  owner_character_id?: number | null;
  owner_faction_id?: number | null;
  created_at: string;
  updated_at: string;
}

export interface PublishResult {
  success: boolean;
  message: string;
  issues?: string[];
  story?: Story;
}

export interface APIResponse<T = any> {
  success: boolean;
  message: string;
  data: T;
  error?: any;
}

export interface PaginatedResponse<T = any> {
  success: boolean;
  message: string;
  data: T[];
  pagination: {
    current_page: number;
    limit: number;
    total_items: number;
    total_pages: number;
  };
}

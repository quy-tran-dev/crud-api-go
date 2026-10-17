import type { CategoryAttributeValue } from './category';

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

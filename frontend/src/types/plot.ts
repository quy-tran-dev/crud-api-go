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

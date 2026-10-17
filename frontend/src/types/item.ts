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

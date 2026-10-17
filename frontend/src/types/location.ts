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

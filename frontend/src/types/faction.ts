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

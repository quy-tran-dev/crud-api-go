import type { SystemCategory } from './category';
import type { Location } from './location';
import type { Faction } from './faction';
import type { Character } from './character';
import type { Plot } from './plot';
import type { Item } from './item';

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

export interface PublishResult {
  success: boolean;
  message: string;
  issues?: string[];
  story?: Story;
}

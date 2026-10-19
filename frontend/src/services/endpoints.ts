// Cây định nghĩa toàn bộ RESTful API Endpoints tập trung của hệ thống
export const ENDPOINTS = {
  AUTH: {
    LOGIN: '/auth/login',
    REGISTER: '/auth/register',
    PROFILE: '/auth/profile',
  },
  STORIES: {
    BASE: '/stories',
    DETAIL: (storyId: number) => `/stories/${storyId}`,
    PUBLISH: (storyId: number) => `/stories/${storyId}/publish`,
    RECHECK: (storyId: number) => `/stories/${storyId}/recheck`,
    APPLY_FRAMEWORK: (storyId: number, frameworkId: number) =>
      `/stories/${storyId}/apply-framework/${frameworkId}`,
    SAVE_FRAMEWORK: (storyId: number) => `/stories/${storyId}/frameworks`,
  },
  CATEGORIES: {
    BASE: (storyId: number) => `/stories/${storyId}/categories`,
    DETAIL: (storyId: number, categoryId: number) =>
      `/stories/${storyId}/categories/${categoryId}`,
    ITEMS: {
      BASE: (storyId: number, categoryId: number) =>
        `/stories/${storyId}/categories/${categoryId}/items`,
      DETAIL: (storyId: number, categoryId: number, itemId: number) =>
        `/stories/${storyId}/categories/${categoryId}/items/${itemId}`,
    },
  },
  CHARACTERS: {
    BASE: (storyId: number) => `/stories/${storyId}/characters`,
    DETAIL: (storyId: number, characterId: number) =>
      `/stories/${storyId}/characters/${characterId}`,
  },
  LOCATIONS: {
    BASE: (storyId: number) => `/stories/${storyId}/locations`,
    DETAIL: (storyId: number, locationId: number) =>
      `/stories/${storyId}/locations/${locationId}`,
  },
  FACTIONS: {
    BASE: (storyId: number) => `/stories/${storyId}/factions`,
    DETAIL: (storyId: number, factionId: number) =>
      `/stories/${storyId}/factions/${factionId}`,
  },
  PLOTS: {
    BASE: (storyId: number) => `/stories/${storyId}/plots`,
    DETAIL: (storyId: number, plotId: number) =>
      `/stories/${storyId}/plots/${plotId}`,
  },
  ITEMS: {
    BASE: (storyId: number) => `/stories/${storyId}/items`,
    DETAIL: (storyId: number, itemId: number) =>
      `/stories/${storyId}/items/${itemId}`,
  },
  FRAMEWORKS: {
    BASE: '/frameworks',
    MY: '/frameworks/my',
    PUBLIC: '/frameworks/public',
    DETAIL: (id: number) => `/frameworks/${id}`,
  },
} as const;

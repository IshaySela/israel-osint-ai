import { createSlice } from '@reduxjs/toolkit';
import type { PayloadAction } from '@reduxjs/toolkit';
import { z } from 'zod';

export interface Location {
  name: string;
  lat: string;
  lon: string;
}

export interface OsintEvent<T = unknown> {
  source: string;
  timestamp_epoch: number;
  summary: string;
  raw_message: string;
  locations: Location[];
  data: T;
}

export const TelegramEventDataSchema = z.object({
  channel_title: z.string(),
  channel_main_lang: z.string(),
});

export type TelegramEventData = z.infer<typeof TelegramEventDataSchema>;

export function isTelegramEventData(data: unknown): data is TelegramEventData {
  return TelegramEventDataSchema.safeParse(data).success;
}

export interface TelegramEvent extends OsintEvent<TelegramEventData> {
  source: 'telegram';
}

export type AnyOsintEvent = TelegramEvent;

interface TimeRange {
  fromMinutesAgo: number;
  toMinutesAgo: number;
}

interface EventState {
  events: AnyOsintEvent[];
  selectedEvent: AnyOsintEvent | null;
  searchQuery: string;
  timeRange: TimeRange;
}

const initialState: EventState = {
  events: [],
  selectedEvent: null,
  searchQuery: '',
  timeRange: { fromMinutesAgo: 1440, toMinutesAgo: 0 },
};

export const eventSlice = createSlice({
  name: 'event',
  initialState,
  reducers: {
    setEvents: (state, action: PayloadAction<AnyOsintEvent[]>) => {
      state.events = action.payload;
    },
    addEvent: (state, action: PayloadAction<AnyOsintEvent>) => {
      // Check if event already exists to avoid duplicates
      const exists = state.events.some(
        (e) => e.timestamp_epoch === action.payload.timestamp_epoch && e.summary === action.payload.summary
      );
      if (!exists) {
        state.events = [action.payload, ...state.events].slice(0, 100);
      }
    },
    setSelectedEvent: (state, action: PayloadAction<AnyOsintEvent | null>) => {
      state.selectedEvent = action.payload;
    },
    setSearchQuery: (state, action: PayloadAction<string>) => {
      state.searchQuery = action.payload;
    },
    setTimeRange: (state, action: PayloadAction<TimeRange>) => {
      state.timeRange = action.payload;
    },
  },
});

export const { setEvents, addEvent, setSelectedEvent, setSearchQuery, setTimeRange } = eventSlice.actions;

export default eventSlice.reducer;

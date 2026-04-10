import { createSlice } from '@reduxjs/toolkit';
import type { PayloadAction } from '@reduxjs/toolkit';

export interface Location {
  name: string;
  lat: string;
  lon: string;
}

export interface Event {
  raw_message: string;
  summary: string;
  timestamp_epoch: number;
  channel_title: string;
  channel_main_lang: string;
  source: string;
  locations: Location[];
}

interface TimeRange {
  fromMinutesAgo: number;
  toMinutesAgo: number;
}

interface EventState {
  events: Event[];
  selectedEvent: Event | null;
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
    setEvents: (state, action: PayloadAction<Event[]>) => {
      state.events = action.payload;
    },
    addEvent: (state, action: PayloadAction<Event>) => {
      // Check if event already exists to avoid duplicates
      const exists = state.events.some(
        (e) => e.timestamp_epoch === action.payload.timestamp_epoch && e.summary === action.payload.summary
      );
      if (!exists) {
        state.events = [action.payload, ...state.events].slice(0, 100);
      }
    },
    setSelectedEvent: (state, action: PayloadAction<Event | null>) => {
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

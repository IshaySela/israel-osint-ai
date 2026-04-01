import { createSlice } from '@reduxjs/toolkit';
import type { PayloadAction } from '@reduxjs/toolkit';

interface Location {
  name: string;
  lat: string;
  lon: string;
}

interface Event {
  raw_message: string;
  summary: string;
  timestamp: string;
  locations: Location[];
}

interface EventState {
  events: Event[];
  selectedEvent: Event | null;
  searchQuery: string;
}

const initialState: EventState = {
  events: [],
  selectedEvent: null,
  searchQuery: '',
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
        (e) => e.timestamp === action.payload.timestamp && e.summary === action.payload.summary
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
  },
});

export const { setEvents, addEvent, setSelectedEvent, setSearchQuery } = eventSlice.actions;

export default eventSlice.reducer;

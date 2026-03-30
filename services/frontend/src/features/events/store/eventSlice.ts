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
    setSelectedEvent: (state, action: PayloadAction<Event | null>) => {
      state.selectedEvent = action.payload;
    },
    setSearchQuery: (state, action: PayloadAction<string>) => {
      state.searchQuery = action.payload;
    },
  },
});

export const { setEvents, setSelectedEvent, setSearchQuery } = eventSlice.actions;

export default eventSlice.reducer;

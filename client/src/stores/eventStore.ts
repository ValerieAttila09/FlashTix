import { create } from 'zustand';
import type { EventState, Event } from '../types';

interface EventStore extends EventState {
  setEvents: (events: Event[]) => void;
  setSelectedEvent: (event: Event | null) => void;
  setLoading: (loading: boolean) => void;
}

export const useEventStore = create<EventStore>((set) => ({
  events: [],
  selectedEvent: null,
  loading: false,

  setEvents: (events: Event[]) => set({ events }),
  setSelectedEvent: (selectedEvent: Event | null) => set({ selectedEvent }),
  setLoading: (loading: boolean) => set({ loading }),
}));
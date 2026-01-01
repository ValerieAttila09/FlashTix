import { create } from 'zustand';
import type { TicketState, CartItem } from '../types';

interface TicketStore extends TicketState {
  addToCart: (item: CartItem) => void;
  removeFromCart: (eventId: string, seat: string) => void;
  clearCart: () => void;
  setTickets: (tickets: any[]) => void;
  setLoading: (loading: boolean) => void;
}

export const useTicketStore = create<TicketStore>((set, get) => ({
  tickets: [],
  cart: [],
  loading: false,

  addToCart: (item: CartItem) => {
    const cart = get().cart;
    set({ cart: [...cart, item] });
  },

  removeFromCart: (eventId: string, seat: string) => {
    const cart = get().cart.filter(
      (item) => !(item.eventId === eventId && item.seat === seat)
    );
    set({ cart });
  },

  clearCart: () => set({ cart: [] }),

  setTickets: (tickets: any[]) => set({ tickets }),
  setLoading: (loading: boolean) => set({ loading }),
}));
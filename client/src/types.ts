// Global types for FlashTix

export interface Event {
  id: string;
  name: string;
  description: string;
  date: string;
  venue: string;
  capacity: number;
  created_at: string;
  updated_at: string;
}

export interface Ticket {
  id: string;
  event_id: string;
  user_id: string;
  seat: string;
  status: 'available' | 'reserved' | 'sold';
  price: number;
  reserved_until?: string;
  created_at: string;
  updated_at: string;
}

export interface User {
  id: string;
  email: string;
  name: string;
  created_at: string;
  updated_at: string;
}

export interface Seat {
  id: string;
  row: string;
  number: string;
  status: 'available' | 'reserved' | 'sold';
}

export interface CartItem {
  eventId: string;
  seat: string;
  price: number;
}

export interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
}

export interface EventState {
  events: Event[];
  selectedEvent: Event | null;
  loading: boolean;
}

export interface TicketState {
  tickets: Ticket[];
  cart: CartItem[];
  loading: boolean;
}
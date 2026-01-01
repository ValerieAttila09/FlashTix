import { useEffect } from 'react';
import { useEventStore } from './stores/eventStore';
import { useAuthStore } from './stores/authStore';
import { Button } from './components/ui/Button';
import { Card, CardHeader, CardContent } from './components/ui/Card';
import './App.css';

function App() {
  const { events, loading, setEvents } = useEventStore();
  const { user, isAuthenticated } = useAuthStore();

  useEffect(() => {
    // Fetch events from API
    const fetchEvents = async () => {
      try {
        const response = await fetch('/api/events');
        const data = await response.json();
        setEvents(data);
      } catch (error) {
        console.error('Failed to fetch events:', error);
      }
    };

    fetchEvents();
  }, [setEvents]);

  return (
    <div className="min-h-screen bg-gray-100 p-4">
      <header className="mb-8">
        <h1 className="text-3xl font-bold text-center">FlashTix</h1>
        <p className="text-center text-gray-600">Platform Pemesanan Tiket Event</p>
      </header>

      <main className="max-w-6xl mx-auto">
        {isAuthenticated ? (
          <div className="mb-4">
            <p>Welcome, {user?.name}!</p>
            <Button variant="secondary" size="sm">Logout</Button>
          </div>
        ) : (
          <div className="mb-4">
            <Button variant="primary">Login</Button>
          </div>
        )}

        <section>
          <h2 className="text-2xl font-semibold mb-4">Available Events</h2>
          {loading ? (
            <p>Loading events...</p>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {events.map((event) => (
                <Card key={event.id}>
                  <CardHeader>
                    <h3 className="text-lg font-semibold">{event.name}</h3>
                  </CardHeader>
                  <CardContent>
                    <p className="text-gray-600 mb-2">{event.description}</p>
                    <p className="text-sm text-gray-500">Venue: {event.venue}</p>
                    <p className="text-sm text-gray-500">Date: {new Date(event.date).toLocaleDateString()}</p>
                    <p className="text-sm text-gray-500">Capacity: {event.capacity}</p>
                    <Button className="mt-4">Book Tickets</Button>
                  </CardContent>
                </Card>
              ))}
            </div>
          )}
        </section>
      </main>
    </div>
  );
}

export default App;

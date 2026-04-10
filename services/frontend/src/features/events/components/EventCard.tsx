import React from 'react';
import type { AnyOsintEvent, TelegramEvent } from '../store/eventSlice';
import TelegramEventCard from './TelegramEventCard';

interface EventCardProps {
  event: AnyOsintEvent;
}

const EventCard: React.FC<EventCardProps> = ({ event }) => {
  switch (event.source) {
    case 'telegram':
      return <TelegramEventCard event={event as TelegramEvent} />;
    default:
      console.error(`[EventCard] Unsupported event source: "${event.source}"`, event);
      if (import.meta.env.DEV) {
        return (
          <div className="p-4 mb-3 border border-red-500/50 rounded-lg bg-red-900/20 text-xs font-mono text-red-400">
            Unsupported source: {event.source}
          </div>
        );
      }
      return null;
  }
};

export default EventCard;

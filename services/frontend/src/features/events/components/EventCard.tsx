import React from 'react';
import type { AnyOsintEvent, TelegramEvent } from '../store/eventSlice';
import TelegramEventCard from './TelegramEventCard';
import GenericEventCard from './GenericEventCard';

interface EventCardProps {
  event: AnyOsintEvent;
}

const EventCard: React.FC<EventCardProps> = ({ event }) => {
  switch (event.source) {
    case 'telegram':
      return <TelegramEventCard event={event as TelegramEvent} />;
    default:
      return <GenericEventCard event={event} />;
  }
};

export default EventCard;

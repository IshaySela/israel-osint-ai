import React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { setSelectedEvent } from '../store/eventSlice';
import type { Event } from '../store/eventSlice';
import type { RootState } from '../../../store';

interface EventCardProps {
  event: Event;
}

const EventCard: React.FC<EventCardProps> = ({ event }) => {
  const dispatch = useDispatch();
  const selectedEvent = useSelector((state: RootState) => state.event.selectedEvent);
  const isSelected = selectedEvent?.timestamp_epoch === event.timestamp_epoch && selectedEvent?.raw_message === event.raw_message;

  const formattedDate = new Date(event.timestamp_epoch * 1000).toLocaleTimeString('he-IL', {
    hour: '2-digit',
    minute: '2-digit',
  });

  return (
    <div
      onClick={() => dispatch(setSelectedEvent(event))}
      className={`p-4 mb-3 cursor-pointer transition-all duration-300 border rounded-lg hover:shadow-[0_0_10px_rgba(34,211,238,0.3)] 
        ${isSelected 
          ? 'bg-cyan-500/10 border-cyan-400 border-l-4' 
          : 'bg-slate-900/50 border-slate-800 border-l-4 border-l-slate-700'
        }`}
    >
      <div className="flex justify-between items-start mb-2">
        <span className="text-xs font-mono text-cyan-400">{formattedDate}</span>
        {event.channel_title && (
          <span className="text-xs text-slate-400 font-mono">{event.channel_title}</span>
        )}
      </div>
      <p className="text-sm text-slate-200 leading-relaxed text-right dir-rtl font-medium">
        {event.summary}
      </p>
    </div>
  );
};

export default EventCard;

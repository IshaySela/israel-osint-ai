import React, { useEffect } from 'react';
import { useQuery } from '@apollo/client/react';
import { useDispatch, useSelector } from 'react-redux';
import { GET_LATEST_EVENTS } from '../graphql/queries';
import { setEvents } from '../store/eventSlice';
import type { RootState } from '../../../store';
import EventCard from './EventCard';
import GlassPanel from '../../../components/common/GlassPanel';
import { Radio } from 'lucide-react';

interface GetLatestEventsData {
  latestEvents: any[];
}

const EventSidebar: React.FC = () => {
  const dispatch = useDispatch();
  const events = useSelector((state: RootState) => state.event.events);
  const { loading, error, data } = useQuery<GetLatestEventsData>(GET_LATEST_EVENTS, {
    pollInterval: 30000, // Poll every 30 seconds
  });

  useEffect(() => {
    if (data && data.latestEvents) {
      dispatch(setEvents(data.latestEvents));
    }
  }, [data, dispatch]);

  if (loading && events.length === 0) return (
    <div className="w-full h-full flex items-center justify-center">
      <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-cyan-400"></div>
    </div>
  );
  
  if (error) return <div className="p-4 text-red-400">Error loading events...</div>;

  return (
    <GlassPanel className="h-full w-96 flex flex-col overflow-hidden">
      <div className="p-4 border-b border-cyan-500/20 flex items-center justify-between bg-slate-900/40">
        <h2 className="text-xl font-bold text-cyan-400 flex items-center gap-2">
          <Radio className="w-5 h-5 animate-pulse text-red-500" />
          Live Feed
        </h2>
        <span className="text-xs font-mono text-cyan-600 bg-cyan-950/30 px-2 py-1 rounded border border-cyan-900">
          {events.length} EVENTS
        </span>
      </div>
      
      <div className="flex-1 overflow-y-auto p-4 custom-scrollbar">
        {events.length === 0 ? (
          <p className="text-slate-500 text-center py-10">No events found...</p>
        ) : (
          events.map((event, idx) => (
            <EventCard key={`${event.timestamp}-${idx}`} event={event} />
          ))
        )}
      </div>
    </GlassPanel>
  );
};

export default EventSidebar;

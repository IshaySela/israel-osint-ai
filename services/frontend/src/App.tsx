import React from 'react';
import EventSidebar from './features/events/components/EventSidebar';
import EventMap from './features/map/components/EventMap';

const App: React.FC = () => {
  return (
    <div className="h-full w-full bg-slate-950 flex relative overflow-hidden">
      {/* Background Decorative Elements */}
      <div className="absolute inset-0 bg-[radial-gradient(circle_at_50%_50%,rgba(34,211,238,0.05),transparent_70%)] pointer-events-none"></div>
      <div className="absolute top-0 left-0 w-full h-1 bg-gradient-to-r from-transparent via-cyan-500/50 to-transparent"></div>

      {/* Sidebar - Highest Z-Index to stay above map */}
      <div className="z-20 relative pointer-events-auto h-full">
        <EventSidebar />
      </div>

      {/* Map Implementation */}
      <div className="flex-1 relative z-10 h-full">
        <EventMap />
      </div>
    </div>
  );
};


export default App;


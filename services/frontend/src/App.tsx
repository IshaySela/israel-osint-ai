import React from 'react';
import EventSidebar from './features/events/components/EventSidebar';

const App: React.FC = () => {
  return (
    <div className="min-h-screen bg-slate-950 flex relative overflow-hidden">
      {/* Background Decorative Elements */}
      <div className="absolute inset-0 bg-[radial-gradient(circle_at_50%_50%,rgba(34,211,238,0.05),transparent_70%)] pointer-events-none"></div>
      <div className="absolute top-0 left-0 w-full h-1 bg-gradient-to-r from-transparent via-cyan-500/50 to-transparent"></div>
      
      {/* Sidebar - Highest Z-Index to stay above map */}
      <div className="z-20">
        <EventSidebar />
      </div>

      {/* Map Placeholder */}
      <div className="flex-1 relative z-10 bg-slate-900 flex items-center justify-center">
        <div className="text-center">
          <div className="w-16 h-16 border-4 border-cyan-500/20 border-t-cyan-500 rounded-full animate-spin mx-auto mb-4"></div>
          <p className="text-cyan-600 font-mono tracking-widest text-sm uppercase">Initializing Geospatial Data...</p>
        </div>
        
        {/* Overlaying instructions/info for Map development */}
        <div className="absolute bottom-10 right-10 p-4 border border-cyan-500/20 bg-slate-950/40 text-cyan-400 font-mono text-xs rounded">
          MAP_LAYER: CARTODB_DARK_MATTER
        </div>
      </div>
    </div>
  );
};

export default App;

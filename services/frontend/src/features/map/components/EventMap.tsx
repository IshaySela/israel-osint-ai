import React, { useEffect } from 'react';
import { MapContainer, TileLayer, Marker, Popup, useMap } from 'react-leaflet';
import L from 'leaflet';
import { useSelector, useDispatch } from 'react-redux';
import type { RootState } from '../../../store';
import { setSelectedEvent } from '../../events/store/eventSlice';
import 'leaflet/dist/leaflet.css';

// Fix for default Leaflet icon
import markerIcon from 'leaflet/dist/images/marker-icon.png';
import markerShadow from 'leaflet/dist/images/marker-shadow.png';

const DefaultIcon = L.icon({
  iconUrl: markerIcon,
  shadowUrl: markerShadow,
  iconSize: [25, 41],
  iconAnchor: [12, 41],
});

L.Marker.prototype.options.icon = DefaultIcon;

// Custom Marker Creator
const createPulseIcon = (isSelected: boolean) => {
  return L.divIcon({
    className: 'custom-pulse-icon',
    html: `
      <div class="relative flex items-center justify-center">
        <div class="absolute w-4 h-4 rounded-full bg-cyan-500 ${isSelected ? 'animate-ping scale-150' : 'animate-pulse opacity-50'}"></div>
        <div class="relative w-3 h-3 rounded-full bg-cyan-400 border border-white shadow-[0_0_10px_rgba(34,211,238,0.8)]"></div>
      </div>
    `,
    iconSize: [20, 20],
    iconAnchor: [10, 10],
  });
};

// Map Recenter Component
const MapRecenter: React.FC = () => {
  const map = useMap();
  const selectedEvent = useSelector((state: RootState) => state.event.selectedEvent);

  useEffect(() => {
    if (selectedEvent && selectedEvent.locations.length > 0) {
      const { lat, lon } = selectedEvent.locations[0];
      map.flyTo([parseFloat(lat), parseFloat(lon)], 13, {
        duration: 1.5,
      });
    }
  }, [selectedEvent, map]);

  return null;
};

const EventMap: React.FC = () => {
  const dispatch = useDispatch();
  const events = useSelector((state: RootState) => state.event.events);
  const selectedEvent = useSelector((state: RootState) => state.event.selectedEvent);

  return (
    <div className="w-full h-full relative z-0">
      <MapContainer
        center={[31.0461, 34.8516]} // Center on Israel
        zoom={8}
        className="w-full h-full"
        zoomControl={false}
      >
        <TileLayer
          attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors &copy; <a href="https://carto.com/attributions">CARTO</a>'
          url="https://{s}.basemaps.cartocdn.com/dark_all/{z}/{x}/{y}{r}.png"
        />
        
        <MapRecenter />

        {events.map((event, eventIdx) => 
          event.locations.map((loc, locIdx) => {
            const isSelected = selectedEvent?.timestamp === event.timestamp && selectedEvent?.raw_message === event.raw_message;
            return (
              <Marker
                key={`${event.timestamp}-${eventIdx}-${locIdx}`}
                position={[parseFloat(loc.lat), parseFloat(loc.lon)]}
                icon={createPulseIcon(isSelected)}
                eventHandlers={{
                  click: () => dispatch(setSelectedEvent(event)),
                }}
              >
                <Popup className="custom-popup">
                  <div className="p-2 dir-rtl text-right">
                    <h3 className="font-bold text-slate-900 mb-1">{loc.name}</h3>
                    <p className="text-xs text-slate-700">{event.summary}</p>
                    <span className="text-[10px] text-slate-500 font-mono mt-2 block">
                      {new Date(event.timestamp).toLocaleString('he-IL')}
                    </span>
                  </div>
                </Popup>
              </Marker>
            );
          })
        )}
      </MapContainer>
      
      {/* Legend / Overlay Controls can go here */}
      <div className="absolute bottom-6 right-6 z-1000 pointer-events-none">
        <div className="bg-slate-950/80 backdrop-blur-md border border-cyan-500/20 p-3 rounded text-[10px] font-mono text-cyan-500 space-y-1">
          <div className="flex items-center gap-2">
            <span className="w-2 h-2 rounded-full bg-cyan-400 shadow-[0_0_5px_rgba(34,211,238,0.8)]"></span>
            ACTIVE OSINT EVENT
          </div>
          <div className="opacity-50">SRCE: TELEGRAM_SCRAPER</div>
          <div className="opacity-50">LYR: CARTODB_DARK</div>
        </div>
      </div>
    </div>
  );
};

export default EventMap;

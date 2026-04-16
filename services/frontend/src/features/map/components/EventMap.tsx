import React, { useEffect, useState, useMemo } from 'react';
import { MapContainer, TileLayer, Marker, Popup, GeoJSON } from 'react-leaflet';
import L from 'leaflet';
import type { FeatureCollection, Polygon, MultiPolygon } from 'geojson';
import { createPulseIcon } from './CreatePulseIcon'
import booleanPointInPolygon from '@turf/boolean-point-in-polygon';
import { point } from '@turf/helpers';
import { useSelector, useDispatch } from 'react-redux';
import type { RootState } from '../../../store';
import { setSelectedEvent } from '../../events/store/eventSlice';
import 'leaflet/dist/leaflet.css';
import MapRecenter from './MapRecenter';
// Fix for default Leaflet icon
import markerIcon from 'leaflet/dist/images/marker-icon.png';
import markerShadow from 'leaflet/dist/images/marker-shadow.png';

interface MuniProperties {
  CR_PNIM: string;
  [key: string]: unknown;
}


const DefaultIcon = L.icon({
  iconUrl: markerIcon,
  shadowUrl: markerShadow,
  iconSize: [25, 41],
  iconAnchor: [12, 41],
});

L.Marker.prototype.options.icon = DefaultIcon;


const EventMap: React.FC = () => {
  const dispatch = useDispatch();
  const events = useSelector((state: RootState) => state.event.events);
  const selectedEvent = useSelector((state: RootState) => state.event.selectedEvent);
  const [muniData, setMuniData] = useState<FeatureCollection<Polygon | MultiPolygon, MuniProperties> | null>(null);

  useEffect(() => {
    fetch('/muni_il_compressed.json')
      .then(r => r.json())
      .then(setMuniData);
  }, []);

  const litMuniIds = useMemo<Set<string>>(() => {
    if (!muniData) return new Set();
    const pts = events.flatMap(e => e.locations.map(l => point([parseFloat(l.lon), parseFloat(l.lat)])));
    const lit = new Set<string>();
    for (const feature of muniData.features) {
      if (pts.some(pt => booleanPointInPolygon(pt, feature)))
        lit.add(feature.properties.CR_PNIM);
    }
    return lit;
  }, [muniData, events]);

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

        {muniData && (
          <GeoJSON
            data={muniData}
            style={(feature) => {
              const id = (feature?.properties as MuniProperties | null)?.CR_PNIM;
              return id && litMuniIds.has(id)
                ? { color: '#22d3ee', weight: 1.5, fillOpacity: 0.2, fillColor: '#22d3ee' }
                : { color: 'transparent', weight: 0, fillOpacity: 0 };
            }}
          />
        )}

        {events.map((event, eventIdx) => 
          event.locations.map((loc, locIdx) => {
            const isSelected = selectedEvent?.timestamp_epoch === event.timestamp_epoch && selectedEvent?.raw_message === event.raw_message;
            return (
              <Marker
                key={`${event.timestamp_epoch}-${eventIdx}-${locIdx}`}
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
                      {new Date(event.timestamp_epoch * 1000).toLocaleString('he-IL')}
                    </span>
                  </div>
                </Popup>
              </Marker>
            );
          })
        )}
      </MapContainer>
      
    </div>
  );
};

export default EventMap;

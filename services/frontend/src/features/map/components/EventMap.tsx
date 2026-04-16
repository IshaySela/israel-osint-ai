import React, { useEffect, useState, useMemo } from 'react';
import { MapContainer, TileLayer, GeoJSON } from 'react-leaflet';
import L from 'leaflet';
import type { FeatureCollection, Polygon, MultiPolygon } from 'geojson';
import booleanPointInPolygon from '@turf/boolean-point-in-polygon';
import { point } from '@turf/helpers';
import { useSelector } from 'react-redux';
import type { RootState } from '../../../store';
import 'leaflet/dist/leaflet.css';
import MapRecenter from './MapRecenter';
// Fix for default Leaflet icon
import markerIcon from 'leaflet/dist/images/marker-icon.png';
import markerShadow from 'leaflet/dist/images/marker-shadow.png';
import { EventPopup } from './EventPopup';

interface MuniProperties {
  CR_PNIM: string;
  Muni_Eng: string;
  Muni_Heb: string;
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
  const events = useSelector((state: RootState) => state.event.events);
  const [muniData, setMuniData] = useState<FeatureCollection<Polygon | MultiPolygon, MuniProperties> | null>(null);

  useEffect(() => {
    fetch('/muni_il_compressed.json')
      .then(r => r.json())
      .then(setMuniData);
  }, []);

  const muniEventCounts = useMemo<Map<string, number>>(() => {
    if (!muniData) return new Map();
    const pts = events.flatMap(e => e.locations.map(l => point([parseFloat(l.lon), parseFloat(l.lat)])));
    const counts = new Map<string, number>();
    for (const feature of muniData.features) {
      const count = pts.filter(pt => booleanPointInPolygon(pt, feature)).length;
      if (count > 0) counts.set(feature.properties.CR_PNIM, count);
    }
    return counts;
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
              return id && muniEventCounts.has(id)
                ? { color: '#22d3ee', weight: 1.5, fillOpacity: 0.2, fillColor: '#22d3ee' }
                : { color: 'transparent', weight: 0, fillOpacity: 0 };
            }}
            onEachFeature={(feature, layer) => {
              const props = feature.properties as MuniProperties;
              const count = muniEventCounts.get(props.CR_PNIM);
              if (count) {
                layer.bindTooltip(
                  `<div style="font-family:monospace;font-size:11px"><b>${props.Muni_Eng}</b><br/>${count} event${count > 1 ? 's' : ''}</div>`,
                  { sticky: true }
                );
              }
            }}
          />
        )}

        {events.map((event, eventIdx) => 
          <div key={`event-popup-${eventIdx}`}><EventPopup event={event} /></div>
        )}
      </MapContainer>
      
    </div>
  );
};

export default EventMap;

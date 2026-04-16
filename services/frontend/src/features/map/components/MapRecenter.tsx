import type { RootState } from '../../../store';
import { useEffect } from "react";
import { useMap } from "react-leaflet";
import { useSelector } from "react-redux";

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


export default MapRecenter;
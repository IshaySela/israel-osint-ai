import type { RootState } from "../../../store";
import { setSelectedEvent, type AnyOsintEvent } from "../../events/store/eventSlice";
import { createPulseIcon } from "./CreatePulseIcon";
import { Marker, Popup } from 'react-leaflet';
import { useDispatch, useSelector } from 'react-redux';

export function EventPopup({ event }: { event: AnyOsintEvent; }) {
    const selectedEvent = useSelector((state: RootState) => state.event.selectedEvent);

    const dispatch = useDispatch();
    return <>
        {
            event.locations.map((loc, locIdx) => {
                const isSelected = selectedEvent?.timestamp_epoch === event.timestamp_epoch && selectedEvent?.raw_message === event.raw_message;
                return (
                    <Marker
                        key={`${event.timestamp_epoch}-${locIdx}`}
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
        }
    </>
}
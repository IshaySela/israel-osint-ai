import React, { useCallback, useEffect, useRef, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { setTimeRange } from '../store/eventSlice';
import type { RootState } from '../../../store';

const MAX_HOURS = 72;

const TimeRangeSlider: React.FC = () => {
  const dispatch = useDispatch();
  const { fromHoursAgo, toHoursAgo } = useSelector(
    (state: RootState) => state.event.timeRange
  );

  const [localFrom, setLocalFrom] = useState(fromHoursAgo);
  const [localTo, setLocalTo] = useState(toHoursAgo);

  // Sync local state if Redux changes externally
  useEffect(() => { setLocalFrom(fromHoursAgo); }, [fromHoursAgo]);
  useEffect(() => { setLocalTo(toHoursAgo); }, [toHoursAgo]);

  const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  const debounced = useCallback(
    (from: number, to: number) => {
      if (debounceRef.current) clearTimeout(debounceRef.current);
      debounceRef.current = setTimeout(() => {
        dispatch(setTimeRange({ fromHoursAgo: from, toHoursAgo: to }));
      }, 150);
    },
    [dispatch]
  );

  const handleFromChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const val = Number(e.target.value);
      if (val > localTo) {
        setLocalFrom(val);
        debounced(val, localTo);
      }
    },
    [localTo, debounced]
  );

  const handleToChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const val = Number(e.target.value);
      if (val < localFrom) {
        setLocalTo(val);
        debounced(localFrom, val);
      }
    },
    [localFrom, debounced]
  );

  const label =
    localTo === 0
      ? `Last ${localFrom}h`
      : `${localTo}h\u2013${localFrom}h ago`;

  // Track fill: slider goes 0 (left=now) to 72 (right=72h ago)
  const leftPct = (localTo / MAX_HOURS) * 100;
  const rightPct = ((MAX_HOURS - localFrom) / MAX_HOURS) * 100;

  const thumbClass =
    'absolute w-full appearance-none bg-transparent pointer-events-none ' +
    '[&::-webkit-slider-thumb]:pointer-events-auto ' +
    '[&::-webkit-slider-thumb]:appearance-none ' +
    '[&::-webkit-slider-thumb]:w-4 [&::-webkit-slider-thumb]:h-4 ' +
    '[&::-webkit-slider-thumb]:rounded-full ' +
    '[&::-webkit-slider-thumb]:border [&::-webkit-slider-thumb]:border-white/40 ' +
    '[&::-webkit-slider-thumb]:cursor-pointer ';

  return (
    <div className="px-4 py-3 border-b border-cyan-500/20 bg-slate-900/40">
      <div className="flex items-center justify-between mb-2">
        <span className="text-xs font-mono text-cyan-500 uppercase tracking-wider">
          Time Window
        </span>
        <span className="text-xs font-mono text-cyan-300 bg-cyan-950/40 px-2 py-0.5 rounded border border-cyan-900">
          {label}
        </span>
      </div>

      <div className="relative h-5 flex items-center">
        <div className="absolute w-full h-1 bg-slate-700 rounded-full" />
        <div
          className="absolute h-1 bg-cyan-500/60 rounded-full"
          style={{ left: `${leftPct}%`, right: `${rightPct}%` }}
        />
        {/* "to" handle — closer to now (lower value) */}
        <input
          type="range"
          min={0}
          max={MAX_HOURS}
          value={localTo}
          onChange={handleToChange}
          step={1}
          className={
            thumbClass +
            '[&::-webkit-slider-thumb]:bg-slate-300 ' +
            '[&::-webkit-slider-thumb]:shadow-[0_0_4px_rgba(255,255,255,0.4)]'
          }
        />
        {/* "from" handle — further back (higher value) */}
        <input
          type="range"
          min={0}
          max={MAX_HOURS}
          value={localFrom}
          onChange={handleFromChange}
          step={1}
          className={
            thumbClass +
            '[&::-webkit-slider-thumb]:bg-cyan-400 ' +
            '[&::-webkit-slider-thumb]:shadow-[0_0_6px_rgba(34,211,238,0.8)]'
          }
        />
      </div>

      <div className="flex justify-between mt-1">
        <span className="text-[10px] font-mono text-slate-500">Now</span>
        <span className="text-[10px] font-mono text-slate-500">72h ago</span>
      </div>
    </div>
  );
};

export default TimeRangeSlider;

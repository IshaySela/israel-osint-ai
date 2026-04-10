import React, { useCallback, useEffect, useRef, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { setTimeRange } from '../store/eventSlice';
import type { RootState } from '../../../store';

// Logarithmic tick scale in minutes: finer steps near "now", coarser for large ranges
const TICKS = [0, 15, 30, 45, 60, 90, 120, 180, 240, 360, 480, 720, 1080, 1440, 2160, 2880, 4320];
const MAX_IDX = TICKS.length - 1;

function minutesToIdx(minutes: number): number {
  let best = 0;
  let bestDist = Infinity;
  for (let i = 0; i < TICKS.length; i++) {
    const d = Math.abs(TICKS[i] - minutes);
    if (d < bestDist) { bestDist = d; best = i; }
  }
  return best;
}

function formatMinutes(m: number): string {
  if (m === 0) return 'now';
  if (m < 60) return `${m}m`;
  const h = m / 60;
  return Number.isInteger(h) ? `${h}h` : `${Math.floor(h)}h${m % 60}m`;
}

const TimeRangeSlider: React.FC = () => {
  const dispatch = useDispatch();
  const { fromMinutesAgo, toMinutesAgo } = useSelector(
    (state: RootState) => state.event.timeRange
  );

  const [localFromIdx, setLocalFromIdx] = useState(() => minutesToIdx(fromMinutesAgo));
  const [localToIdx, setLocalToIdx] = useState(() => minutesToIdx(toMinutesAgo));

  useEffect(() => { setLocalFromIdx(minutesToIdx(fromMinutesAgo)); }, [fromMinutesAgo]);
  useEffect(() => { setLocalToIdx(minutesToIdx(toMinutesAgo)); }, [toMinutesAgo]);

  const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  const debounced = useCallback(
    (fromIdx: number, toIdx: number) => {
      if (debounceRef.current) clearTimeout(debounceRef.current);
      debounceRef.current = setTimeout(() => {
        dispatch(setTimeRange({ fromMinutesAgo: TICKS[fromIdx], toMinutesAgo: TICKS[toIdx] }));
      }, 150);
    },
    [dispatch]
  );

  const handleFromChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const idx = Number(e.target.value);
      if (idx > localToIdx) {
        setLocalFromIdx(idx);
        debounced(idx, localToIdx);
      }
    },
    [localToIdx, debounced]
  );

  const handleToChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const idx = Number(e.target.value);
      if (idx < localFromIdx) {
        setLocalToIdx(idx);
        debounced(localFromIdx, idx);
      }
    },
    [localFromIdx, debounced]
  );

  const fromMin = TICKS[localFromIdx];
  const toMin = TICKS[localToIdx];
  const label = toMin === 0
    ? `Last ${formatMinutes(fromMin)}`
    : `${formatMinutes(toMin)}\u2013${formatMinutes(fromMin)} ago`;

  const leftPct = (localToIdx / MAX_IDX) * 100;
  const rightPct = ((MAX_IDX - localFromIdx) / MAX_IDX) * 100;

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
        {/* "to" handle — closer to now */}
        <input
          type="range"
          min={0}
          max={MAX_IDX}
          step={1}
          value={localToIdx}
          onChange={handleToChange}
          className={
            thumbClass +
            '[&::-webkit-slider-thumb]:bg-slate-300 ' +
            '[&::-webkit-slider-thumb]:shadow-[0_0_4px_rgba(255,255,255,0.4)]'
          }
        />
        {/* "from" handle — further back */}
        <input
          type="range"
          min={0}
          max={MAX_IDX}
          step={1}
          value={localFromIdx}
          onChange={handleFromChange}
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

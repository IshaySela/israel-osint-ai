import React from 'react';

interface GlassPanelProps {
  children: React.ReactNode;
  className?: string;
}

const GlassPanel: React.FC<GlassPanelProps> = ({ children, className = '' }) => {
  return (
    <div className={`bg-slate-950/80 backdrop-blur-md border border-cyan-500/30 rounded-lg shadow-[0_0_15px_rgba(34,211,238,0.1)] ${className}`}>
      {children}
    </div>
  );
};

export default GlassPanel;

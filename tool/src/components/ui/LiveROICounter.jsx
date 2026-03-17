import { useState, useEffect, useRef } from 'react';

export default function LiveROICounter({ value }) {
  const [displayValue, setDisplayValue] = useState(0);
  const prevValue = useRef(0);
  const animRef = useRef(null);

  useEffect(() => {
    if (animRef.current) cancelAnimationFrame(animRef.current);
    const start = prevValue.current;
    const end = value;
    const duration = 800;
    const startTime = performance.now();

    const animate = (now) => {
      const elapsed = now - startTime;
      const progress = Math.min(elapsed / duration, 1);
      const eased = 1 - Math.pow(1 - progress, 3);
      setDisplayValue(Math.round(start + (end - start) * eased));
      if (progress < 1) animRef.current = requestAnimationFrame(animate);
    };

    animRef.current = requestAnimationFrame(animate);
    prevValue.current = end;
    return () => { if (animRef.current) cancelAnimationFrame(animRef.current); };
  }, [value]);

  return (
    <div className="live-roi-counter">
      <div className="live-roi-label">Potential Annual Savings Identified</div>
      <div className="live-roi-value">${displayValue.toLocaleString()}</div>
    </div>
  );
}

import React from 'react';

export default function PartProgressBar({ totalParts, partsPassed, currentPart, size = 'md' }) {
  if (!totalParts || totalParts <= 1) return null;

  const dots = [];
  for (let i = 1; i <= totalParts; i++) {
    let color, title;
    if (i <= partsPassed) {
      color = '#01b328'; // green — passed
      title = `Part ${i}: passed`;
    } else if (i === currentPart) {
      color = '#1a90ff'; // blue — current
      title = `Part ${i}: current`;
    } else {
      color = '#d1d5db'; // gray — locked
      title = `Part ${i}: locked`;
    }
    dots.push({ i, color, title });
  }

  const dotSize  = size === 'sm' ? 6 : 8;
  const fontSize = size === 'sm' ? '10px' : '11px';

  return (
    <div
      className="flex items-center gap-1.5"
      title={`${partsPassed} of ${totalParts} parts completed`}
    >
      {dots.map(({ i, color, title }) => (
        <span
          key={i}
          title={title}
          style={{
            width:  dotSize,
            height: dotSize,
            borderRadius: '50%',
            background: color,
            display: 'inline-block',
            flexShrink: 0,
          }}
        />
      ))}
      <span className="text-text-tertiary" style={{ fontSize }}>
        {partsPassed}/{totalParts}
      </span>
    </div>
  );
}

import React from 'react';

const ControlsCard = ({ isWatching, onStart, error }) => (
  <div className="section-card">
    <div className="section-title">Controls</div>
    <button 
      className={`main-btn ${isWatching ? 'active' : ''}`}
      onClick={onStart}
      disabled={isWatching}
    >
      {isWatching ? 'Monitoring Active' : 'Start Monitoring'}
    </button>
    <div className="status-row">
      <span>Status:</span>
      <span className={isWatching ? 'status-active' : 'status-idle'}>{isWatching ? 'Active' : 'Idle'}</span>
    </div>
    {error && <div className="error-msg">{error}</div>}
  </div>
);

export default ControlsCard; 
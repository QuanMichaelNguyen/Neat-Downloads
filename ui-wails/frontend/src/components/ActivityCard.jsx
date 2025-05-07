import React from 'react';

const ActivityCard = ({ recentActivity }) => (
  <div className="section-card">
    <div className="section-title blue">Recent Activity</div>
    <ul className="activity-list">
      {recentActivity.length === 0 && <li>No files categorized yet.</li>}
      {recentActivity.map((item, idx) => (
        <li key={idx}>
          <strong>{item.fileName}</strong> → <span className="blue">{item.category}</span> <span className="activity-time">({item.time})</span>
        </li>
      ))}
    </ul>
  </div>
);

export default ActivityCard; 
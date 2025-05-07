import React from 'react';

const SettingsCard = ({ watchDir, dropboxEnabled, dropboxFolder, onDropboxToggle }) => (
  <div className="section-card">
    <div className="section-title">Settings</div>
    <div className="setting-row"><span>Watch Directory:</span><span>{watchDir}</span></div>
    <div className="setting-row">
      <span>Dropbox Sync:</span>
      <label className="switch">
        <input type="checkbox" checked={dropboxEnabled} onChange={onDropboxToggle} />
        <span className="slider round"></span>
      </label>
      <span style={{marginLeft:8, color: dropboxEnabled ? '#2563eb' : '#888'}}>{dropboxEnabled ? 'Enabled' : 'Disabled'}</span>
    </div>
    {dropboxEnabled && dropboxFolder && (
      <div className="setting-row"><span>Dropbox Folder:</span><span>{dropboxFolder}</span></div>
    )}
  </div>
);

export default SettingsCard; 
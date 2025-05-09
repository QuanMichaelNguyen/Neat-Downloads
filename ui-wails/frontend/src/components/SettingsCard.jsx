import React, { useState, useEffect } from 'react';

const SettingsCard = ({ watchDir, dropboxEnabled, dropboxFolder, syncToDropbox, onDropboxToggle, onSaveDropboxSettings }) => {
    const [localSyncToDropbox, setLocalSyncToDropbox] = useState(syncToDropbox);
    const [dirty, setDirty] = useState(false);

    useEffect(() => {
        setLocalSyncToDropbox(syncToDropbox);
        setDirty(false);
    }, [syncToDropbox]);

    const handleToggle = () => {
        setLocalSyncToDropbox((prev) => !prev);
        setDirty(true);
    };

    const handleSave = () => {
        if (onSaveDropboxSettings) {
            onSaveDropboxSettings({ syncToDropbox: localSyncToDropbox });
            setDirty(false);
        }
    };

    return (
        <div className="card">
            <h2 className="text-xl font-bold mb-4">Settings</h2>
            <div className="space-y-4">
                <div>
                    <label className="block text-sm font-medium text-gray-700">Watch Directory</label>
                    <p className="mt-1 text-sm text-gray-500">{watchDir}</p>
                </div>
                
                <div className="flex items-center justify-between">
                    <div>
                        <label className="block text-sm font-medium text-gray-700">Dropbox Integration</label>
                        <p className="text-sm text-gray-500">
                            {dropboxEnabled ? 'Connected to Dropbox' : 'Not connected'}
                        </p>
                        {dropboxEnabled && dropboxFolder && (
                            <p className="text-sm text-gray-500">Folder: {dropboxFolder}</p>
                        )}
                    </div>
                    {dropboxEnabled && (
                        <button
                            onClick={onDropboxToggle}
                            className="main-btn"
                            style={{ boxShadow: 'none' }}
                        >
                            Disconnect
                        </button>
                    )}
                </div>
                {dropboxEnabled && (
                    <div className="flex items-center justify-between mt-2">
                        <label className="text-sm font-medium text-gray-700" htmlFor="syncToDropboxToggle">Upload to Dropbox</label>
                        <input
                            id="syncToDropboxToggle"
                            type="checkbox"
                            checked={localSyncToDropbox}
                            onChange={handleToggle}
                            className="accent-blue-500 h-4 w-4 rounded focus:ring-2 focus:ring-blue-400"
                        />
                    </div>
                )}
                {dropboxEnabled && (
                    <button
                        onClick={handleSave}
                        disabled={!dirty}
                        className="main-btn"
                        style={{ minWidth: 80 }}
                    >
                        Save
                    </button>
                )}
            </div>
        </div>
    );
};

export default SettingsCard; 
import {useState, useEffect} from 'react';
import {GetSettings, StartWatching} from '../wailsjs/go/main/App';
import {EventsOn} from '../wailsjs/runtime/runtime';
import Header from './components/Header';
import SettingsCard from './components/SettingsCard';
import ControlsCard from './components/ControlsCard';
import ActivityCard from './components/ActivityCard';
import Toast from './components/Toast';
import './App.css';

function App() {
    const [settings, setSettings] = useState({});
    const [isWatching, setIsWatching] = useState(false);
    const [status, setStatus] = useState("Idle");
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState(null);
    const [recentActivity, setRecentActivity] = useState([]);
    const [toast, setToast] = useState("");
    const [dropboxEnabled, setDropboxEnabled] = useState(false);

    useEffect(() => {
        loadSettings();
        const unsubscribe = EventsOn("fileCategorized", (data) => {
            if (data && data.fileName && data.category) {
                handleFileCategorized(data.fileName, data.category);
            }
        });
        return () => unsubscribe();
    }, []);

    const loadSettings = async () => {
        try {
            setIsLoading(true);
            setError(null);
            const result = await GetSettings();
            setSettings(result);
            setDropboxEnabled(result.dropboxEnabled || false);
        } catch (err) {
            setError('Failed to load settings. Please try again.');
        } finally {
            setIsLoading(false);
        }
    };

    const handleStartWatching = async () => {
        try {
            setError(null);
            const result = await StartWatching();
            setStatus(result);
            setIsWatching(true);
        } catch (err) {
            setError('Failed to start monitoring. Please try again.');
        }
    };

    const handleFileCategorized = (fileName, category) => {
        setRecentActivity(prev => [
            { fileName, category, time: new Date().toLocaleTimeString() },
            ...prev.slice(0, 7)
        ]);
        setToast(`Categorized ${fileName} to ${category}`);
        setTimeout(() => setToast(""), 3500);
    };

    // Simulate Dropbox toggle (backend integration can be added)
    const handleDropboxToggle = () => {
        setDropboxEnabled((prev) => !prev);
        // TODO: Call backend to update Dropbox setting if needed
    };

    return (
        <div className="app-bg-light">
            <Header />
            <div className="main-content split-layout">
                <div className="left-panel">
                    <SettingsCard
                        watchDir={settings.watchDir}
                        dropboxEnabled={dropboxEnabled}
                        dropboxFolder={settings.dropboxFolder}
                        onDropboxToggle={handleDropboxToggle}
                    />
                    <ControlsCard
                        isWatching={isWatching}
                        onStart={handleStartWatching}
                        error={error}
                    />
                </div>
                <div className="right-panel">
                    <ActivityCard recentActivity={recentActivity} />
                </div>
            </div>
            <Toast message={toast} />
        </div>
    );
}

export default App;
import React, { useEffect, useState } from "react";

const DropboxCallback: React.FC = () => {
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const handleCallback = () => {
      const urlParams = new URLSearchParams(window.location.search);
      const code = urlParams.get("code");
      const error = urlParams.get("error");

      if (error) {
        setError(`Authentication failed: ${error}`);
        return;
      }

      if (code) {
        try {
          // Send the code back to the main window
          window.opener.postMessage(
            {
              type: "dropbox-auth-callback",
              code: code,
            },
            "*"
          );
          // Close this window after a short delay
          setTimeout(() => window.close(), 1000);
        } catch (err) {
          setError("Failed to complete authentication. Please try again.");
        }
      } else {
        setError("No authentication code received. Please try again.");
      }
    };

    handleCallback();
  }, []);

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-50">
      <div className="text-center p-8 bg-white rounded-lg shadow-lg max-w-md w-full">
        {error ? (
          <>
            <div className="text-red-600 mb-4">
              <svg
                className="w-12 h-12 mx-auto"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
            </div>
            <h1 className="text-xl font-bold text-gray-900 mb-2">
              Authentication Failed
            </h1>
            <p className="text-gray-600">{error}</p>
            <button
              onClick={() => window.close()}
              className="mt-4 px-4 py-2 bg-gray-100 hover:bg-gray-200 text-gray-800 rounded-md transition-colors"
            >
              Close Window
            </button>
          </>
        ) : (
          <>
            <div className="text-green-600 mb-4">
              <svg
                className="w-12 h-12 mx-auto animate-bounce"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M5 13l4 4L19 7"
                />
              </svg>
            </div>
            <h1 className="text-xl font-bold text-gray-900 mb-2">
              Successfully Connected!
            </h1>
            <p className="text-gray-600">You can close this window now.</p>
          </>
        )}
      </div>
    </div>
  );
};

export default DropboxCallback;

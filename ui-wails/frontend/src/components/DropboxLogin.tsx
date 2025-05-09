import React from "react";

interface DropboxLoginProps {
  onLoginSuccess: (token: string) => void;
}

const DropboxLogin: React.FC<DropboxLoginProps> = ({ onLoginSuccess }) => {
  const handleLogin = async () => {
    try {
      const response = await window.go.main.App.GetDropboxAuthURL();
      window.open(response, "_blank", "width=800,height=600");

      window.addEventListener("message", async (event) => {
        if (event.data.type === "dropbox-auth-callback") {
          const code = event.data.code;
          const result = await window.go.main.App.ExchangeDropboxCode(code);
          onLoginSuccess(result.accessToken);
        }
      });
    } catch (error) {
      console.error("Dropbox login failed:", error);
    }
  };

  return (
    <button onClick={handleLogin} className="main-btn">
      Login with Dropbox
    </button>
  );
};

export default DropboxLogin;

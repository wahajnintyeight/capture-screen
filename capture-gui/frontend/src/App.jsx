import { useState } from "react";
import "./App.css";

function App() {
  const [isEnabled, setIsEnabled] = useState(false);
  const [status, setStatus] = useState("Stopped");
  const [isLoading, setIsLoading] = useState(false);

  const handleToggle = async () => {
    try {
      setIsLoading(true);
      if (!isEnabled) {
        console.log("Starting capture service");
        await window.go.main.App.StartCaptureService();
        setStatus("Running");
      } else {
        console.log("Stopping capture service");
        await window.go.main.App.StopCaptureService();
        setStatus("Stopped");
      }
      setIsEnabled(!isEnabled);
    } catch (error) {
      console.error("Error toggling service:", error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleMinimize = () => window.runtime.WindowMinimise();
  const handleClose = () => window.runtime.Quit();

  return (
    <div className="min-h-screen relative overflow-hidden flex flex-col font-inter" style={{ "--wails-draggable": "drag" }}>
      {/* Draggable titlebar */}
      <div className="titlebar relative z-20 flex items-center justify-between px-4 py-3 bg-black/10 backdrop-blur-md border-b border-white/10">
        <div className="flex items-center gap-3">
          <div className="w-2 h-2 rounded-full bg-white/20" />
          <span className="text-white/80 text-sm font-medium">
            Screen Capture
          </span>
        </div>
        <div className="flex gap-4">
          <button
            onClick={handleMinimize}
            className="no-drag w-6 h-6 flex items-center justify-center hover:bg-gray-400/30 transition-colors"
            title="Minimize"
          >
            <div className="w-3 h-0.5 bg-white hover:bg-white/90 transition-colors" />
          </button>
          <button
            onClick={handleClose}
            className="w-6 h-6 flex items-center justify-center hover:bg-gray-400/30 transition-colors"
            title="Close"
          >
            <svg
              className="w-3.5 h-3.5 text-white hover:text-white/90 transition-colors"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M6 18L18 6M6 6l12 12"
              />
            </svg>
          </button>
        </div>
      </div>

      {/* Side drag handles */}
      <div className="absolute inset-y-0 left-0 w-1 titlebar cursor-ew-resize z-50" />
      <div className="absolute inset-y-0 right-0 w-1 titlebar cursor-ew-resize z-50" />
      <div className="absolute inset-x-0 bottom-0 h-1 titlebar cursor-ns-resize z-50" />
      <div className="absolute bottom-0 left-0 w-2 h-2 titlebar cursor-nesw-resize z-50" />
      <div className="absolute bottom-0 right-0 w-2 h-2 titlebar cursor-nwse-resize z-50" />

      {/* Main Content */}
      <div className="flex-1 flex flex-col items-center justify-center px-6">
        {/* Animated background */}
        <div className="absolute inset-0 bg-gradient-to-br from-purple-500 via-pink-500 to-yellow-500 animate-gradient-xy"></div>
        <div className="absolute inset-0 bg-black/20 backdrop-blur-[2px]"></div>

        {/* Animated circles */}
        <div className="absolute inset-0 overflow-hidden">
          <div className="circle-1"></div>
          <div className="circle-2"></div>
          <div className="circle-3"></div>
        </div>

        {/* Content */}
        <div className="relative z-10 space-y-12 text-center">
          <h1 className="text-4xl font-bold text-white tracking-tight">
            Screen Capture Service
          </h1>

          <div className="flex flex-col items-center gap-8">
            <span
              className={`inline-flex items-center px-6 py-2 rounded-full text-sm font-semibold tracking-wide transition-all duration-300 ${
                status === "Running"
                  ? "bg-green-400/20 text-green-300 ring-1 ring-green-400/30"
                  : "bg-red-400/20 text-red-300 ring-1 ring-red-400/30"
              }`}
            >
              {status}
              {isLoading && (
                <svg className="animate-spin ml-2 h-4 w-4" viewBox="0 0 24 24">
                  <circle
                    className="opacity-25"
                    cx="12"
                    cy="12"
                    r="10"
                    stroke="currentColor"
                    strokeWidth="4"
                  ></circle>
                  <path
                    className="opacity-75"
                    fill="currentColor"
                    d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                  ></path>
                </svg>
              )}
            </span>

            <button
              onClick={handleToggle}
              disabled={isLoading}
              className={`px-8 py-3 rounded-lg font-medium text-white transition-all duration-300 transform hover:scale-105 disabled:opacity-50 disabled:cursor-not-allowed ${
                isEnabled
                  ? "bg-red-500/30 hover:bg-red-500/40 ring-1 ring-red-400/50"
                  : "bg-blue-500/30 hover:bg-blue-500/40 ring-1 ring-blue-400/50"
              }`}
            >
              {isEnabled ? "Stop Capture" : "Start Capture"}
            </button>

            <p className="text-sm text-white/60 font-medium max-w-sm">
              {isEnabled
                ? "Screen capture service is running in the background"
                : "Click Start Capture to begin monitoring"}
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}

export default App;

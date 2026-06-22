using System;
using System.Threading;

namespace WarpCsharp
{
    public class ViewportConfig
    {
        public uint Width { get; set; }
        public uint Height { get; set; }
        public double DeviceScaleFactor { get; set; }
    }

    public class BrowserConfig
    {
        public bool Headless { get; set; }
        public ViewportConfig Viewport { get; set; } = new ViewportConfig();
    }

    public class CursorState
    {
        public double X { get; set; }
        public double Y { get; set; }
    }

    public class CdpPage
    {
        public string Url { get; private set; } = "about:blank";
        public CursorState CursorState { get; private set; } = new CursorState { X = 0.0, Y = 0.0 };
        private readonly object _lock = new object();

        public void Navigate(string url)
        {
            Console.WriteLine($"[Browser:Page] Navigating to: {url}");
            Url = url;
            Thread.Sleep(600); // Simulate load time
            Console.WriteLine($"[Browser:Page] Successfully loaded {url}");
        }

        public void DispatchMouseEvent(double x, double y, string eventType)
        {
            lock (_lock)
            {
                Console.WriteLine($"[Browser:Page] Dispatching mouse {eventType} at ({x}, {y})");
                CursorState.X = x;
                CursorState.Y = y;
            }
        }

        public void DispatchKeyEvent(string text)
        {
            Console.WriteLine($"[Browser:Page] Dispatching keystroke: '{text}'");
        }

        public string CaptureScreenshot()
        {
            Console.WriteLine($"[Browser:Page] Capturing screenshot for {Url}");
            return $"screenshot_data_base64_for_{Url}";
        }
    }

    public class BrowserManager
    {
        public BrowserConfig Config { get; private set; }
        public CdpPage? ActivePage { get; private set; }
        private readonly object _lock = new object();

        public BrowserManager(BrowserConfig config)
        {
            Config = config;
        }

        public void Launch()
        {
            lock (_lock)
            {
                string mode = Config.Headless ? "Headless" : "Windowed";
                Console.WriteLine($"[Browser:Manager] Launching {mode} Chrome via CDP...");
                ActivePage = new CdpPage();
            }
        }
    }
}

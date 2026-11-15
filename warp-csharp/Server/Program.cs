using System;
using System.IO;
using System.Net;
using System.Text.Json;
using System.Threading.Tasks;

namespace Server
{
    class Program
    {
        static async Task Main(string[] args)
        {
            HttpListener listener = new HttpListener();
            listener.Prefixes.Add("http://localhost:8081/");
            listener.Start();
            Console.WriteLine("[Init] Warp C# Server starting on HTTP :8081...");

            while (true)
            {
                HttpListenerContext context = await listener.GetContextAsync();
                _ = Task.Run(() => HandleRequest(context));
            }
        }

        static async Task HandleRequest(HttpListenerContext context)
        {
            HttpListenerRequest request = context.Request;
            HttpListenerResponse response = context.Response;

            response.AddHeader("Access-Control-Allow-Origin", "*");
            response.AddHeader("Access-Control-Allow-Methods", "GET, POST, OPTIONS");
            response.AddHeader("Access-Control-Allow-Headers", "Content-Type, Authorization");

            if (request.HttpMethod == "OPTIONS")
            {
                response.StatusCode = (int)HttpStatusCode.OK;
                response.Close();
                return;
            }

            if (request.Url.AbsolutePath == "/")
            {
                byte[] buffer = System.Text.Encoding.UTF8.GetBytes("Warp C# Server Backend is operational.\n\nUse /status for health checks.");
                response.ContentType = "text/plain";
                response.ContentLength64 = buffer.Length;
                await response.OutputStream.WriteAsync(buffer, 0, buffer.Length);
            }
            else if (request.Url.AbsolutePath == "/status")
            {
                var statusObj = new
                {
                    status = "online",
                    version = "1.0.0",
                    message = "Warp Agent HTTP/WebSocket Orchestrator Active",
                    ready = true
                };
                string jsonString = JsonSerializer.Serialize(statusObj);
                byte[] buffer = System.Text.Encoding.UTF8.GetBytes(jsonString);

                response.ContentType = "application/json";
                response.ContentLength64 = buffer.Length;
                await response.OutputStream.WriteAsync(buffer, 0, buffer.Length);
            }
            else
            {
                response.StatusCode = (int)HttpStatusCode.NotFound;
            }

            response.Close();
        }
    }
}

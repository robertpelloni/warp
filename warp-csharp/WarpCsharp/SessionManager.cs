using System;
using System.Collections.Generic;

namespace WarpCsharp
{
    public enum SessionEntryType
    {
        Message,
        ModelChange,
        Compaction,
        BranchSummary
    }

    public class SessionEntryBase
    {
        public SessionEntryType Type { get; set; }
        public string Id { get; set; } = string.Empty;
        public string ParentId { get; set; } = string.Empty;
        public string Timestamp { get; set; } = string.Empty;
    }

    public class SessionMessageEntry : SessionEntryBase
    {
        public string Message { get; set; } = string.Empty;
    }

    public class ModelChangeEntry : SessionEntryBase
    {
        public string Provider { get; set; } = string.Empty;
        public string ModelId { get; set; } = string.Empty;
    }

    public class SessionManager
    {
        public string SessionId { get; private set; }
        public string Cwd { get; private set; }
        public string SessionDir { get; private set; }
        public List<object> FileEntries { get; private set; } = new List<object>();

        public SessionManager(string sessionId, string cwd, string sessionDir)
        {
            SessionId = sessionId;
            Cwd = cwd;
            SessionDir = sessionDir;
        }

        public void AppendMessageEntry(string id, string parentId, string message)
        {
            var entry = new SessionMessageEntry
            {
                Type = SessionEntryType.Message,
                Id = id,
                ParentId = parentId,
                Timestamp = "now",
                Message = message
            };

            FileEntries.Add(entry);
            Console.WriteLine($"[SessionManager] Appended Message: {message}");
        }
    }
}

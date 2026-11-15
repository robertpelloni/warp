using System;
using System.Collections.Generic;

namespace WarpCsharp
{
    public class StoredThread
    {
        public string ThreadId { get; set; } = string.Empty;
        public string ParentThreadId { get; set; } = string.Empty;
        public string Preview { get; set; } = string.Empty;
        public long CreatedAtMs { get; set; }
        public long UpdatedAtMs { get; set; }
        public string Cwd { get; set; } = string.Empty;
    }

    public class StoredThreadHistory
    {
        public string ThreadId { get; set; } = string.Empty;
        public List<string> Items { get; set; } = new List<string>();
    }

    public interface IThreadStore
    {
        void CreateThread(string threadId, string cwd);
        void AppendItem(string threadId, string item);
        StoredThread LoadThread(string threadId);
        StoredThreadHistory LoadHistory(string threadId);
    }

    public class InMemoryThreadStore : IThreadStore
    {
        private readonly Dictionary<string, StoredThread> _threads = new Dictionary<string, StoredThread>();
        private readonly Dictionary<string, List<string>> _histories = new Dictionary<string, List<string>>();
        private readonly object _lock = new object();

        public void CreateThread(string threadId, string cwd)
        {
            lock (_lock)
            {
                _threads[threadId] = new StoredThread
                {
                    ThreadId = threadId,
                    Preview = "New Session",
                    Cwd = cwd
                };
                _histories[threadId] = new List<string>();
            }
        }

        public void AppendItem(string threadId, string item)
        {
            lock (_lock)
            {
                if (_histories.ContainsKey(threadId))
                {
                    _histories[threadId].Add(item);
                }
                else
                {
                    throw new Exception($"Thread {threadId} not found");
                }
            }
        }

        public StoredThread LoadThread(string threadId)
        {
            lock (_lock)
            {
                if (_threads.TryGetValue(threadId, out var thread))
                {
                    return thread;
                }
                throw new Exception("Thread not found");
            }
        }

        public StoredThreadHistory LoadHistory(string threadId)
        {
            lock (_lock)
            {
                if (_histories.TryGetValue(threadId, out var items))
                {
                    return new StoredThreadHistory { ThreadId = threadId, Items = new List<string>(items) };
                }
                throw new Exception("History not found");
            }
        }
    }
}

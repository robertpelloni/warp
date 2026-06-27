using System;
using Xunit;
using WarpCsharp;

public class ThreadStoreTests
{
    [Fact]
    public void InMemoryThreadStore_CanCreateAndLoadThread()
    {
        var store = new InMemoryThreadStore();
        string threadId = "test_thread_1";
        string cwd = "/test/dir";

        store.CreateThread(threadId, cwd);

        var loadedThread = store.LoadThread(threadId);
        Assert.NotNull(loadedThread);
        Assert.Equal(threadId, loadedThread.ThreadId);
        Assert.Equal(cwd, loadedThread.Cwd);
        Assert.Equal("New Session", loadedThread.Preview);
    }

    [Fact]
    public void InMemoryThreadStore_CanAppendAndLoadHistory()
    {
        var store = new InMemoryThreadStore();
        string threadId = "test_thread_2";

        store.CreateThread(threadId, "/");

        store.AppendItem(threadId, "Item 1");
        store.AppendItem(threadId, "Item 2");

        var history = store.LoadHistory(threadId);
        Assert.NotNull(history);
        Assert.Equal(threadId, history.ThreadId);
        Assert.Equal(2, history.Items.Count);
        Assert.Equal("Item 1", history.Items[0]);
        Assert.Equal("Item 2", history.Items[1]);
    }

    [Fact]
    public void InMemoryThreadStore_ThrowsOnInvalidThread()
    {
        var store = new InMemoryThreadStore();
        Assert.Throws<Exception>(() => store.LoadThread("non_existent"));
        Assert.Throws<Exception>(() => store.LoadHistory("non_existent"));
        Assert.Throws<Exception>(() => store.AppendItem("non_existent", "item"));
    }
}

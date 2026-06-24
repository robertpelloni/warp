using System;
using System.Collections.Generic;

namespace WarpCsharp
{
    public class FileOperations
    {
        public HashSet<string> Read { get; set; } = new HashSet<string>();
        public HashSet<string> Written { get; set; } = new HashSet<string>();
        public HashSet<string> Edited { get; set; } = new HashSet<string>();
    }

    public class CompactionPreparation
    {
        public string FirstKeptEntryId { get; set; } = string.Empty;
        public List<string> MessagesToSummarize { get; set; } = new List<string>();
        public List<string> TurnPrefixMessages { get; set; } = new List<string>();
        public bool IsSplitTurn { get; set; }
        public int TokensBefore { get; set; }
        public string PreviousSummary { get; set; } = string.Empty;
        public FileOperations FileOps { get; set; } = new FileOperations();
    }

    public class CompactionResult
    {
        public string Summary { get; set; } = string.Empty;
        public string NewEntryId { get; set; } = string.Empty;
    }

    public static class CompactionService
    {
        public static CompactionPreparation PrepareCompaction(List<object> entries)
        {
            Console.WriteLine($"[Compaction] Analyzing {entries.Count} session entries for compaction...");

            if (entries.Count == 0)
            {
                throw new Exception("no entries to compact");
            }

            var prep = new CompactionPreparation
            {
                FirstKeptEntryId = "entry_xyz",
                MessagesToSummarize = new List<string> { "User asked to build feature", "Assistant planned feature" },
                TokensBefore = 4500,
                PreviousSummary = "Previous session context..."
            };

            prep.FileOps.Read.Add("/workspace/main.go");
            prep.FileOps.Edited.Add("/workspace/README.md");

            return prep;
        }

        public static CompactionResult Compact(CompactionPreparation prep)
        {
            Console.WriteLine($"[Compaction] Compacting {prep.TokensBefore} tokens down to summary...");

            string summary = $"Summarized {prep.MessagesToSummarize.Count} messages including: {prep.MessagesToSummarize[0]}";

            return new CompactionResult
            {
                Summary = summary,
                NewEntryId = "compact_abc123"
            };
        }
    }
}

package dev.warp;

import java.util.Scanner;

class ToolOrchestrator {
    private boolean sandboxEnabled = true;

    public String executeTask(String task) {
        System.out.println("[Orchestrator] Requesting approval for task: '" + task + "'");
        try {
            Thread.sleep(500);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }

        if (sandboxEnabled) {
            System.out.println("[Orchestrator] Running in sandbox mode...");
        }

        String[] steps = {"Plan", "Code", "Review"};
        for (String step : steps) {
            System.out.println("[Agent: " + step + "] Processing...");
            try {
                Thread.sleep(400);
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
            }
        }

        return "Task '" + task + "' completed successfully by Auto Drive.";
    }
}

public class Main {
    public static void main(String[] args) {
        System.out.println("Welcome to Warp CLI (Java Edition) - Inspired by just-every-code");
        System.out.println("Type '/help' for commands, or 'quit' to close.");

        ToolOrchestrator orchestrator = new ToolOrchestrator();
        Scanner scanner = new Scanner(System.in);

        while (true) {
            System.out.print("warp> ");
            if (!scanner.hasNextLine()) {
                break;
            }

            String input = scanner.nextLine().trim();
            if (input.isEmpty()) {
                continue;
            }

            if (input.equalsIgnoreCase("quit") || input.equalsIgnoreCase("/quit")) {
                break;
            }

            handleCommand(input, orchestrator);
        }
        scanner.close();
    }

    private static void handleCommand(String input, ToolOrchestrator orchestrator) {
        if (input.startsWith("/")) {
            String[] parts = input.substring(1).split(" ", 2);
            String cmd = parts[0].toLowerCase();
            String args = parts.length > 1 ? parts[1].trim() : "";

            switch (cmd) {
                case "help":
                    System.out.println("Available commands:");
                    System.out.println("  /help        - Show this help message");
                    System.out.println("  /auto <task> - Start Auto Drive orchestration for a task");
                    System.out.println("  quit         - Quit the application");
                    break;
                case "auto":
                    if (args.isEmpty()) {
                        System.out.println("[Error] /auto requires a task description.");
                    } else {
                        System.out.println("[Orchestrator: Auto Drive] Starting autonomous loop.");
                        String result = orchestrator.executeTask(args);
                        System.out.println("[Result] " + result);
                    }
                    break;
                default:
                    System.out.println("Unknown command: /" + cmd);
                    break;
            }
        } else {
            System.out.println("[Agent] Echoing input: " + input);
        }
    }
}

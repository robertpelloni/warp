package dev.warp;

import java.util.Scanner;

public class Main {
    public static void main(String[] args) {
        System.out.println("Welcome to Warp CLI (Java Edition) - Inspired by just-every-code");
        System.out.println("Type '/help' for commands, or 'quit' to close.");

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

            handleCommand(input);
        }
        scanner.close();
    }

    private static void handleCommand(String input) {
        if (input.startsWith("/")) {
            String cmd = input.substring(1).toLowerCase();
            switch (cmd) {
                case "help":
                    System.out.println("Available commands:");
                    System.out.println("  /help     - Show this help message");
                    System.out.println("  /plan     - (Stub) Coordinate planning agent");
                    System.out.println("  /code     - (Stub) Coordinate coding agent");
                    System.out.println("  /auto     - (Stub) Start Auto Drive orchestration");
                    System.out.println("  quit      - Quit the application");
                    break;
                case "plan":
                    System.out.println("[Agent: Planner] Acknowledged. Ready to plan task.");
                    break;
                case "code":
                    System.out.println("[Agent: Coder] Acknowledged. Ready to write code.");
                    break;
                case "auto":
                    System.out.println("[Orchestrator: Auto Drive] Starting autonomous loop (Mock).");
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

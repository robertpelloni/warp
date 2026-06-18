package dev.warp;

public class Main {
    public static void main(String[] args) {
        if (args.length == 0) {
            System.out.println("Warp Java Initialized");
            System.out.println("Usage: warp <command> [arguments]");
            return;
        }

        String command = args[0];

        switch (command) {
            case "run":
                String agent = "default";
                for (int i = 1; i < args.length; i++) {
                    if ("--agent".equals(args[i]) && i + 1 < args.length) {
                        agent = args[i + 1];
                    }
                }
                System.out.println("Warp Java running agent: " + agent);
                break;
            default:
                System.out.println("Unknown command: " + command);
                break;
        }
    }
}

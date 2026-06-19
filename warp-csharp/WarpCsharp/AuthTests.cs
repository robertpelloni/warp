using System;

namespace WarpCsharp.Tests
{
    // A simple manual test stub for now, as adding xUnit would require changing project structure significantly
    public class AuthTests
    {
        public static void RunTests()
        {
            Environment.SetEnvironmentVariable("WARP_ADMIN_USER", "admin");
            Environment.SetEnvironmentVariable("WARP_ADMIN_HASH", "8c6976e5b5410415bde908bd4dee15dfb167a9c873fc4bb8a81f6f2ab448a918");

            if (!Auth.AuthenticateUser("admin", "admin")) throw new Exception("Test Failed");
            if (Auth.AuthenticateUser("user", "pass")) throw new Exception("Test Failed");
        }
    }
}

using System;
using System.Security.Cryptography;
using System.Text;

namespace WarpCsharp
{
    public class Auth
    {
        public static bool AuthenticateUser(string username, string password)
        {
            Console.WriteLine($"Warp C#: Authenticating user {username}");

            string expectedUser = Environment.GetEnvironmentVariable("WARP_ADMIN_USER");
            string expectedHash = Environment.GetEnvironmentVariable("WARP_ADMIN_HASH");

            if (string.IsNullOrEmpty(expectedUser) || string.IsNullOrEmpty(expectedHash))
            {
                Console.WriteLine("Warning: Authentication not configured");
                return false;
            }

            using (SHA256 sha256Hash = SHA256.Create())
            {
                byte[] bytes = sha256Hash.ComputeHash(Encoding.UTF8.GetBytes(password));
                StringBuilder builder = new StringBuilder();
                for (int i = 0; i < bytes.Length; i++)
                {
                    builder.Append(bytes[i].ToString("x2"));
                }
                string passwordHash = builder.ToString();

                bool userMatch = CryptographicOperations.FixedTimeEquals(Encoding.UTF8.GetBytes(username), Encoding.UTF8.GetBytes(expectedUser));
                bool hashMatch = CryptographicOperations.FixedTimeEquals(Encoding.UTF8.GetBytes(passwordHash), Encoding.UTF8.GetBytes(expectedHash));

                return userMatch && hashMatch;
            }
        }
    }
}

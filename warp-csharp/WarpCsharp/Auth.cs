using System;
using System.Security.Cryptography;
using System.Text;

namespace WarpCsharp
{
    public class Auth
    {
        // In a real application, this would be retrieved from a database.
        private const string DummyHash = "8c6976e5b5410415bde908bd4dee15dfb167a9c873fc4bb8a81f6f2ab448a918";

        public static bool AuthenticateUser(string username, string password)
        {
            Console.WriteLine($"Warp C#: Authenticating user {username}");

            using (SHA256 sha256Hash = SHA256.Create())
            {
                byte[] bytes = sha256Hash.ComputeHash(Encoding.UTF8.GetBytes(password));
                StringBuilder builder = new StringBuilder();
                for (int i = 0; i < bytes.Length; i++)
                {
                    builder.Append(bytes[i].ToString("x2"));
                }
                string passwordHash = builder.ToString();

                return username == "admin" && passwordHash == DummyHash;
            }
        }
    }
}

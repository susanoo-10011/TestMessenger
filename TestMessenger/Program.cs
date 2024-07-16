using Npgsql;
using System.Net;
using System.Text;
using TestMessenger.Entity;
using Newtonsoft;
using Newtonsoft.Json;
using AuthenticationService;

namespace TestMessenger
{
    internal class Program
    {

        private static readonly string _connectionString = "Server=localhost;Port=5433;Database=MessengerDataBase;UserId=postgres;Password=12345;";
        static async Task Main(string[] args)
        {
            HttpListener listener = new HttpListener();
            listener.Prefixes.Add("http://localhost:5000/");
            listener.Start();

            Console.WriteLine("Listening...");

            while (true)
            {
                HttpListenerContext context = listener.GetContext();
                ProcessRequestAsync(context).Wait();
            }
        }

        private static async Task ProcessRequestAsync(HttpListenerContext context)
        {
            var request = context.Request;
            var response = context.Response;

            // Проверяем, что это POST-запрос и что путь совпадает с /auth/login
            if (request.HttpMethod == "POST" && request.RawUrl == "/auth/login")
            {
                using (var reader = new StreamReader(request.InputStream))
                {
                    string json = await reader.ReadToEndAsync();

                    User body = JsonConvert.DeserializeObject<User>(json);

                    UserRepository userRepository = new UserRepository(_connectionString);
                    bool checkUserLogin = userRepository.CheckUserInBD(body);

                    if (checkUserLogin)
                    {
                        userRepository.AddUser(body);

                        string jsonResponse = "{\"status\":\"success\",\"message\":\"Logged in successfully\"}";
                        byte[] buffer = Encoding.UTF8.GetBytes(jsonResponse);
                        response.ContentLength64 = buffer.Length;
                        Stream output = response.OutputStream;
                        await output.WriteAsync(buffer, 0, buffer.Length);
                        output.Close();
                    }
                    else
                    {
                        string jsonResponse = "{\"status\":\"error\",\"message\":\"Invalid login credentials\"}";
                        byte[] buffer = Encoding.UTF8.GetBytes(jsonResponse);
                        response.ContentLength64 = buffer.Length;
                        Stream output = response.OutputStream;
                        await output.WriteAsync(buffer, 0, buffer.Length);
                        output.Close();
                    }
                }
            }
            else
            {
                // Обработка других запросов или возврат стандартного сообщения
                string responseString = "<HTML><BODY> Invalid Request</BODY></HTML>";
                byte[] buffer = Encoding.UTF8.GetBytes(responseString);
                response.ContentLength64 = buffer.Length;
                Stream output = response.OutputStream;
                await output.WriteAsync(buffer, 0, buffer.Length);
                output.Close();
            }
        }

    }
}

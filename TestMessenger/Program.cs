using System.Net;
using System.Text;
using TestMessenger.Entity;
using Newtonsoft.Json;
using AuthenticationService;
using AuthenticationService.Entity;
using System.Security.Cryptography;

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
            HttpListenerRequest request = context.Request;
            HttpListenerResponse response = context.Response;

            // Проверяем, что это POST-запрос и что путь совпадает с /registration
            if (request.HttpMethod == "POST" && request.RawUrl == "/registration")
            {
                using (var reader = new StreamReader(request.InputStream))
                {
                    string json = await reader.ReadToEndAsync();

                    User body = JsonConvert.DeserializeObject<User>(json);

                    UserRepository userRepository = new UserRepository(_connectionString);
                    bool checkUserLogin = await userRepository.CheckUserLogin(body);

                    if (checkUserLogin)
                    {
                        User user = await userRepository.AddUser(body);

                        string token = "";
                        while (true)
                        {
                            token = GenerateApiKey();

                            if (await userRepository.CheckToken(token))
                            {
                                break;
                            }
                            else
                            {
                                continue;
                            }
                        }

                        TokenRecord tokenRecord = new TokenRecord();
                        tokenRecord.user_id = user.user_id;
                        tokenRecord.token = token;

                        long tokenId = await userRepository.AddToken(user, token);

                        if(tokenId != 0)
                            tokenRecord.id = tokenId;

                        var jsonData = new
                        {
                            status = "success",
                            message = "Logged in successfully",
                            tokenEntity = tokenRecord
                        };

                        string jsonResponse = JsonConvert.SerializeObject(jsonData);

                        await SendReply(response, jsonResponse, 200);
                    }
                    else
                    {
                        string jsonResponse = "{\"status\":\"error\",\"message\":\"This login already exists\"}";

                        await SendReply(response, jsonResponse, 401);
                    }
                }
            }
            else if (request.HttpMethod == "POST" && request.RawUrl == "/enterTheSystem")
            {
                using (StreamReader reader = new StreamReader(request.InputStream))
                {
                    string json = await reader.ReadToEndAsync();
                    User body = JsonConvert.DeserializeObject<User>(json);

                    UserRepository userRepository = new UserRepository(_connectionString);
                    long userId = await userRepository.CheckLoginPassword(body);

                    if (userId != 0)
                    {

                        TokenRecord tokenRecord = new TokenRecord();
                        tokenRecord.user_id = userId;
                        tokenRecord.token = "token";

                        var jsonData = new
                        {
                            status = "success",
                            message = "Successful login",
                            tokenRecord = tokenRecord
                        };

                        string jsonResponse = JsonConvert.SerializeObject(jsonData);

                        await SendReply(response, jsonResponse, 200);
                    }
                    else if (userId == 0)
                    {
                        string jsonResponse = "{\"status\":\"error\",\"message\":\"Invalid login or password\"}";

                        await SendReply(response, jsonResponse, 402);
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

        public static string GenerateApiKey()
        {
            using (var rng = new RNGCryptoServiceProvider())
            {
                byte[] randomBytes = new byte[32];
                rng.GetBytes(randomBytes);

                string apiKey = Convert.ToBase64String(randomBytes);

                return apiKey.Replace("/", "_").Replace("+", "-");
            }
        }

        private static async Task SendReply(HttpListenerResponse response, string jsonResponse, int statusCode)
        {
            byte[] buffer = Encoding.UTF8.GetBytes(jsonResponse);
            response.ContentLength64 = buffer.Length;
            response.StatusCode = statusCode;
            Stream output = response.OutputStream;
            await output.WriteAsync(buffer, 0, buffer.Length);
            output.Close();
        }
    }
}

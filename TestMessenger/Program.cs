using System.Net;
using System.Text;
using TestMessenger.Entity;
using Newtonsoft.Json;
using AuthenticationService;
using AuthenticationService.Entity;
using System.Security.Cryptography;
using System.Text.RegularExpressions;

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

                    ProcessingRegistrationRequests(response, body);

                }
            }
            else if (request.HttpMethod == "POST" && request.RawUrl == "/enterTheSystem")
            {
                using (StreamReader reader = new StreamReader(request.InputStream))
                {
                    string json = await reader.ReadToEndAsync();
                    User body = JsonConvert.DeserializeObject<User>(json);

                    ProcessingLoginRequests(response, body);
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

        #region Registration

        private static async Task ProcessingRegistrationRequests(HttpListenerResponse response, User body)
        {
            string pattern = @"^\w+[\w.-]*@\w+([\w-]\w+)*\.\w{2,3}(\.\w{2})?$";

            Regex regex = new Regex(pattern);
            bool isValidEmail = regex.IsMatch(body.email);
            if (!isValidEmail)
            {
                string jsonResponse = "{\"status\":\"error\",\"message\":\"Incorrect email\"}";
                response.StatusDescription = "Incorrect email";
                await SendReply(response, jsonResponse, 403);

                return;
            }

            if (!CheckPassword(body.password))
            {
                string jsonResponse = "{\"status\":\"error\",\"message\":\"Incorrect password\"}";
                response.StatusDescription = "Incorrect password";
                await SendReply(response, jsonResponse, 404);

                return;
            }

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

                if (tokenId != 0)
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
                response.StatusDescription = "This login already exists";
                await SendReply(response, jsonResponse, 401);
            }
        }

        private static bool CheckPassword(string password)
        {
            if (password == null || password.Length < 8)
            {
                return false;
            }

            Regex regex = new Regex(@"[a-zA-Z]");
            if (!regex.IsMatch(password))
            {
                return false;
            }

            return true;
        }

        #endregion

        private static async Task ProcessingLoginRequests(HttpListenerResponse response, User body)
        {
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
                response.StatusDescription = "Invalid login or password";
                await SendReply(response, jsonResponse, 402);
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

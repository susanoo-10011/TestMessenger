using Npgsql;
using TestMessenger.Entity;
using static Microsoft.EntityFrameworkCore.DbLoggerCategory.Database;

namespace AuthenticationService
{
    public class UserRepository
    {
        private readonly string _connectionString;

        public UserRepository(string connectionString)
        {
            _connectionString = connectionString;
        }

        public void Connect()
        {
            using (var connection = new NpgsqlConnection(_connectionString))
            {
                connection.Open();

                Console.WriteLine("Connected to the database.");
            }
        }

        public bool CheckUserInBD(User checkUser)
        {
            using (var connection = new NpgsqlConnection(_connectionString))
            {
                connection.Open();

                using (var command = new NpgsqlCommand())
                {
                    var users = new List<User>();

                    command.Connection = connection;
                    command.CommandText = "SELECT \"Id\", \"Login\", \"Password\", \"Email\" FROM auth_schema.users";
                    using (var reader = command.ExecuteReader())
                    {
                        while (reader.Read())
                        {
                            var user = new User
                            {
                                Id = reader.GetInt32(reader.GetOrdinal("Id")),
                                Login = reader.GetString(reader.GetOrdinal("Login")),
                                Password = reader.GetString(reader.GetOrdinal("Password")),
                                Email = reader.GetString(reader.GetOrdinal("Email")),
                            };

                            users.Add(user);
                        }
                    }

                    foreach (var user in users)
                    {
                        if(user.Login == checkUser.Login)
                        {
                            return false;
                        }
                    }

                    return true;
                }
            }
        }

        public void AddUser(User user)
        {
            using (var connection = new NpgsqlConnection(_connectionString))
            {
                connection.Open();

                using (var command = new NpgsqlCommand())
                {
                    command.Connection = connection;
                    command.CommandText = @"INSERT INTO auth_schema.users (""Login"", ""Password"", ""Email"") VALUES (@Login, @Password, @Email)";
                    command.Parameters.AddWithValue("@Login", $"{user.Login}");
                    command.Parameters.AddWithValue("@Password", $"{user.Password}");
                    command.Parameters.AddWithValue("@Email", $"{user.Email}");

                    command.ExecuteNonQuery();
                }

                Console.WriteLine("Пользователь зарегистрирован");
            }
        }

        public void DeleteUser(int id)
        {
            using (var connection = new NpgsqlConnection(_connectionString))
            {
                connection.Open();

                using (var command = new NpgsqlCommand())
                {
                    command.Connection = connection;
                    command.CommandText = "DELETE FROM auth_schema.users WHERE Id = @Id";
                    command.Parameters.AddWithValue("Id", id);

                    command.ExecuteNonQuery();
                }
            }
        }
    }
}

using Npgsql;
using System.Data.Common;
using TestMessenger.Entity;
using static Microsoft.EntityFrameworkCore.DbLoggerCategory.Database;

namespace AuthenticationService
{
    public class UserRepository
    {
        private readonly string _connectionString;
        private NpgsqlConnection _connection;

        public UserRepository(string connectionString)
        {
            _connectionString = connectionString;
        }

        private async Task<NpgsqlConnection> GetConnectionAsync()
        {
            if (_connection == null || _connection.State != System.Data.ConnectionState.Open)
            {
                _connection = new NpgsqlConnection(_connectionString);
                await _connection.OpenAsync();
            }

            return _connection;
        }

        public async Task AddToken(User user, string token)
        {
            using (var connection = await GetConnectionAsync())
            {
                using (var command = new NpgsqlCommand())
                {
                    command.Connection = connection;
                    command.CommandText = @"INSERT INTO auth_schema.token_records (""user_id"", ""token"") VALUES (@user_id, @token) RETURNING id";
                    command.Parameters.AddWithValue("@user_id", $"{user.user_id}");
                    command.Parameters.AddWithValue("@token", $"{token}");

                    long tokenId = (long)await command.ExecuteScalarAsync();

                    Console.WriteLine($"Токен пользователя создан {tokenId}");
                }

            }
        }

        public async Task<bool> CheckToken(string currentToken)
        {
            using (var connection = await GetConnectionAsync())
            {
                using (var command = new NpgsqlCommand())
                {
                    command.Connection = connection;
                    command.CommandText = "SELECT \"token\" FROM auth_schema.token_records";

                    using (var reader = command.ExecuteReader())
                    {
                        while (reader.Read())
                        {
                            string token = reader.GetString(reader.GetOrdinal("token"));

                            if (token == currentToken)
                            {
                                return false;
                            }
                        }
                    }

                    return true;
                }
            }
        }

        public async Task<bool> CheckUserLogin(User checkUser)
        {
            using (var connection = await GetConnectionAsync())
            {
                using (var command = new NpgsqlCommand())
                {
                    command.Connection = connection;
                    command.CommandText = "SELECT \"login\" FROM auth_schema.users"; 

                    using (var reader = command.ExecuteReader())
                    {
                        while (reader.Read())
                        {
                            var user = new User
                            {
                                login = reader.GetString(reader.GetOrdinal("login")),
                            };

                            if(user.login == checkUser.login)
                            {
                                 return false;
                            }
                        }
                    }

                    return true;
                }
            }
        }

        public async Task<long> CheckLoginPassword(User checkUser)
        {
            using (var connection = await GetConnectionAsync())
            {
                using(var command = new NpgsqlCommand())
                {
                    command.Connection = connection;
                    command.CommandText = "SELECT \"login\", \"password\", \"user_id\" FROM auth_schema.users";
                    using (var reader = command.ExecuteReader())
                    {
                        while (reader.Read())
                        {
                            var user = new User
                            {
                                login = reader.GetString(reader.GetOrdinal("login")),
                                password = reader.GetString(reader.GetOrdinal("password")),
                                user_id = reader.GetOrdinal("user_id"),
                            };

                            if(user.login == checkUser.login)
                            {
                                if (user.password == checkUser.password)
                                {
                                    return user.user_id;
                                }
                            }
                        }
                    }

                    return 0;
                }
            }
        }

        public async Task<User> AddUser(User user)
        {
            using (var connection = await GetConnectionAsync())
            {
                using (var command = new NpgsqlCommand())
                {
                    command.Connection = connection;
                    command.CommandText = @"INSERT INTO auth_schema.users (""login"", ""password"", ""email"") VALUES (@login, @password, @email) RETURNING user_id;";
                    command.Parameters.AddWithValue("@login", $"{user.login}");
                    command.Parameters.AddWithValue("@password", $"{user.password}");
                    command.Parameters.AddWithValue("@email", $"{user.email}");

                    long userId = (long)await command.ExecuteScalarAsync();
                    Console.WriteLine($"Пользователь зарегистрирован. ID: {userId}");

                    user.user_id = userId;
                    return user;
                }
            }
        }

        public async Task DeleteUser(int id)
        {
            using (var connection = await GetConnectionAsync())
            {
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

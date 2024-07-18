﻿using Npgsql;
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

        /// <summary>
        /// Устанавливаем соединение с БД
        /// </summary>
        private async Task<NpgsqlConnection> GetConnectionAsync()
        {
            if (_connection == null || _connection.State != System.Data.ConnectionState.Open)
            {
                _connection = new NpgsqlConnection(_connectionString);
                await _connection.OpenAsync();
            }

            return _connection;
        }

        /// <summary>
        /// Добавляем уникальный токен пользователя в БД
        /// </summary>
        public async Task<long> AddToken(User user, string token)
        {
            using (var connection = await GetConnectionAsync())
            {
                using (var command = new NpgsqlCommand())
                {
                    try
                    {
                        command.Connection = connection;
                        command.CommandText = @"INSERT INTO auth_schema.token_records (user_id, token) 
                                              SELECT u.user_id, :token 
                                              FROM auth_schema.users u
                                              WHERE u.user_id = :user_id
                                              RETURNING id;";

                        command.Parameters.AddWithValue("user_id", user.user_id);
                        command.Parameters.AddWithValue("token", token);

                        long tokenId = (long)await command.ExecuteScalarAsync();
                        Console.WriteLine($"Токен пользователя создан {tokenId}");

                        return tokenId;

                    }
                    catch (PostgresException pgex)
                    {
                        Console.WriteLine($"PostgreSQL Error: {pgex.Message}");
                        Console.WriteLine($"SQL State: {pgex.SqlState}");
                        Console.WriteLine($"Error Code: {pgex.ErrorCode}");

                        return 0;
                    }
                    catch (Exception ex)
                    {
                        Console.WriteLine($"Error: {ex}");
                        return 0;
                    }
                }

            }
        }

        /// <summary>
        /// Проверяем токен на совпадения в базе данных, чтобы не было одинаковых
        /// </summary>
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

        /// <summary>
        /// Проверяем логин на совпадения в базе данных, чтобы не было одинаковых
        /// </summary>
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

                            if (user.login == checkUser.login)
                            {
                                return false;
                            }
                        }
                    }

                    return true;
                }
            }
        }

        /// <summary>
        /// Проверяем логин и пароль пользователя для входа в систему
        /// </summary>
        public async Task<long> CheckLoginPassword(User checkUser)
        {
            using (var connection = await GetConnectionAsync())
            {
                using (var command = new NpgsqlCommand())
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

                            if (user.login == checkUser.login)
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

        /// <summary>
        /// Добавляем нового пользователя
        /// </summary>
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

        /// <summary>
        /// Удаляем пользователя из БД
        /// </summary>
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

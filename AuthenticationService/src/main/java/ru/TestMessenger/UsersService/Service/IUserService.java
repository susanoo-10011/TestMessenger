package ru.TestMessenger.UsersService.Service;

import jakarta.servlet.http.HttpServletRequest;
import ru.TestMessenger.UsersService.Model.TokenRecord;
import ru.TestMessenger.UsersService.Model.User;
import org.springframework.http.ResponseEntity;

import java.util.List;

public interface IUserService {

    ResponseEntity<List<User>> findAllUsers();

    ResponseEntity<?> saveUser(User user);

    ResponseEntity<User> findByEmail(String email);

    ResponseEntity<List<User>> findAllByName(String name);


    ResponseEntity<User> findUserByLogin(String login);

    ResponseEntity<?> getFollowers(String token);

    ResponseEntity<User> updateUser(User user, HttpServletRequest request);

    ResponseEntity<?> logInToSystem(User user);

    ResponseEntity<?> logOutFromSystem(String token);

    ResponseEntity<?> subscribeToUser(String token, String username);
}

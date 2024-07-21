package ru.TestMessenger.AuthenticationService.Service;

import ru.TestMessenger.AuthenticationService.Model.TokenRecord;
import ru.TestMessenger.AuthenticationService.Model.User;
import org.springframework.http.ResponseEntity;

import java.util.List;

public interface IUserService {

    ResponseEntity<List<User>> findAllUsers();

    ResponseEntity<?> saveUser(User user);

    ResponseEntity<User> findByEmail(String email);

    ResponseEntity<User> updateUser(User user);

    ResponseEntity<?> logInToSystem(User user);

}

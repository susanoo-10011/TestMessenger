package ru.TestMessenger.AuthenticationService.Controller;

import ru.TestMessenger.AuthenticationService.Exceptions.CustomException;
import ru.TestMessenger.AuthenticationService.Model.TokenRecord;
import ru.TestMessenger.AuthenticationService.Model.User;
import ru.TestMessenger.AuthenticationService.Service.UserService;
import lombok.AllArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.HashMap;
import java.util.List;
import java.util.Map;

@RestController
@RequestMapping("/api/v1/users")
@AllArgsConstructor
public class AuthorizationController {

    private final UserService _userService;

    @GetMapping()
    public ResponseEntity<List<User>> getUsers() {
        return _userService.findAllUsers();
    }

    @PostMapping("/save_user")
    public ResponseEntity<?> saveUser(@RequestBody User user) {
        return  _userService.saveUser(user);
    }

    @GetMapping("/{email}")
    public ResponseEntity<User> findUserByEmail(@PathVariable String email) {
        return _userService.findByEmail(email);
    }

    @PutMapping("/update_user/")
    public ResponseEntity<User> updateUser(@RequestBody User user) {
        return  _userService.updateUser(user);
    }

    @PutMapping("/logInSystem/")
    public ResponseEntity<?> logInSystem(@RequestBody User user) {
        return _userService.logInToSystem(user);
    }

    @PutMapping("/logOutFromSystem/")
    public ResponseEntity<?> logOutFromSystem(@RequestBody TokenRecord token) {
        return _userService.logOutFromSystem(token);
    }

    @DeleteMapping("/delete_user/{email}")
    public String deleteUser(@PathVariable String email) {
        _userService.deleteUser(email);
        return "Пользователь удален";
    }
}
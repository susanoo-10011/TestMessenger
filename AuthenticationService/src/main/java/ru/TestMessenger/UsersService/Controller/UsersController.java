package ru.TestMessenger.UsersService.Controller;

import jakarta.servlet.http.HttpServletRequest;
import org.springframework.http.HttpStatus;
import ru.TestMessenger.UsersService.Exceptions.CustomException;
import ru.TestMessenger.UsersService.Model.ErrorTDO;
import ru.TestMessenger.UsersService.Model.User;
import ru.TestMessenger.UsersService.Service.UserService;
import lombok.AllArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/api/v1/users")
@AllArgsConstructor
public class UsersController {

    private final UserService _userService;

    // region Get

    @GetMapping()
    public ResponseEntity<List<User>> getUsers() {

        return _userService.findAllUsers();
    }

    @GetMapping("/search/username")
    public ResponseEntity<?> getUsersByLogin(@RequestParam(required = true) String username) {
        if(username != null) {
            return _userService.findUserByLogin(username);
        }
        else {
            ErrorTDO errorTDO = new ErrorTDO();
            errorTDO.status = HttpStatus.BAD_REQUEST.value();
            errorTDO.message = "The parameters are null";
            errorTDO.error = HttpStatus.BAD_REQUEST;

            throw new CustomException(errorTDO);
        }
    }

    @GetMapping("/search/name")
    public ResponseEntity<List<User>> getUsersByName(@RequestParam(required = true) String name) {
        if(name != null) {
            return _userService.findAllByName(name);
        }
        else {
            ErrorTDO errorTDO = new ErrorTDO();
            errorTDO.status = HttpStatus.BAD_REQUEST.value();
            errorTDO.message = "The parameters are null";
            errorTDO.error = HttpStatus.BAD_REQUEST;

            throw new CustomException(errorTDO);
        }
    }

    @GetMapping("/subscriptions")
    public ResponseEntity<List<String>> getSubscriptions(@RequestHeader("X-User-Agent") String header) {
        return _userService.getFollowers(header);
    }

    //endregion


    @PostMapping()
    public ResponseEntity<?> saveUser(@RequestBody User user) {
        return  _userService.saveUser(user);
    }

    @PutMapping("/login")
    public ResponseEntity<?> logInSystem(@RequestBody User user) {
        return _userService.logInToSystem(user);
    }

    @PatchMapping("/subscriptions")
    public ResponseEntity<?> addSubscription(@RequestParam(required = true) String username,
                                             @RequestHeader("X-User-Agent") String header) {

        return  _userService.subscribeToUser(header, username);
    }

    @DeleteMapping("/logOut")
    public ResponseEntity<?> logOutFromSystem(@RequestHeader("X-User-Agent") String header) {
        return _userService.logOutFromSystem(header);
    }

    @DeleteMapping("/delete_user/{email}")
    public String deleteUser(@PathVariable String email) {
        _userService.deleteUser(email);
        return "Пользователь удален";
    }

    //region Profile

    @PutMapping("/profile/")
    public ResponseEntity<?> updateProfile(@RequestBody User user, HttpServletRequest request) {

        return _userService.updateUser(user, request);
    }

    //endregion
}
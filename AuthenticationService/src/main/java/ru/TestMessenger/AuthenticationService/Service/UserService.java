package ru.TestMessenger.AuthenticationService.Service;

import ru.TestMessenger.AuthenticationService.Exceptions.CustomException;
import ru.TestMessenger.AuthenticationService.Model.TokenRecord;
import ru.TestMessenger.AuthenticationService.Model.User;
import ru.TestMessenger.AuthenticationService.Repository.IUserRepository;
import ru.TestMessenger.AuthenticationService.Repository.IUserTokenRepository;
import ru.TestMessenger.AuthenticationService.Service.EntityResponseTDO.SaveUserResponseTDO;
import lombok.AllArgsConstructor;
import org.springframework.context.annotation.Primary;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.lang.Nullable;
import org.springframework.stereotype.Service;
import org.springframework.dao.DataIntegrityViolationException;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

@Service
@AllArgsConstructor
@Primary
@Transactional
public class UserService implements IUserService{

    private final IUserRepository _userRepository;
    private final IUserTokenRepository _tokenRepository;


    @Override
    public ResponseEntity<List<User>> findAllUsers() {
        List<User> users = _userRepository.findAll();

        if (users.isEmpty()) {
            throw new CustomException("Users not found", HttpStatus.NOT_FOUND);
        }
        return new ResponseEntity<>(users, HttpStatus.OK);
    }

    @Override
    public ResponseEntity<?> saveUser(User user) {

        CheckIncomingParametersToCreateUser(user);

        try {
            User savedUser = _userRepository.save(user);

            TokenRecord savedToken = _tokenRepository.save(createToken(user));

            return new ResponseEntity<>(CreateUserResponseTDO(savedToken), HttpStatus.OK);
        } catch (DataIntegrityViolationException e) {
            throw new CustomException("A user with this login already exists.", HttpStatus.CONFLICT);
        } catch (Exception e) {
            throw new CustomException("Failed to save user", HttpStatus.INTERNAL_SERVER_ERROR);
        }
    }

    @Override
    public ResponseEntity<User> findByEmail(String email) {
        User user = _userRepository.findUserByEmail(email);
        if (user == null) {
            throw new CustomException("User not found", HttpStatus.NOT_FOUND);
        }
        return new ResponseEntity<>(user, HttpStatus.OK);
    }

    @Override
    public ResponseEntity<User> updateUser(User user) {
        try{
            Optional<User> optionalUser = _userRepository.findById(user.getId());

            if (optionalUser.isEmpty()) {
                throw new CustomException("User not found", HttpStatus.NOT_FOUND);
            }

            _userRepository.save(user);
            return new ResponseEntity<>(user, HttpStatus.OK);
        }
        catch(Exception e){
            throw new CustomException("Failed to update user", HttpStatus.INTERNAL_SERVER_ERROR);
        }
    }

    @Override
    public ResponseEntity<?> logInToSystem(User user) {
        User foundUser = checkLoginPasswordLogInSystem(user);
        TokenRecord token = createToken(foundUser);

        try {
            _tokenRepository.save(token);
            return new ResponseEntity<>(CreateUserResponseTDO(token), HttpStatus.OK);
        }
        catch (DataIntegrityViolationException e) {
            throw new CustomException("Failed to save token", HttpStatus.CONFLICT);
        }
        catch(Exception e){
            throw new CustomException("Failed to save token", HttpStatus.INTERNAL_SERVER_ERROR);
        }
    }

    public ResponseEntity<?> logOutFromSystem(TokenRecord token) {
        try{
            _tokenRepository.deleteByTokenValue(token.getTokenValue());
            return new ResponseEntity<>(HttpStatus.OK);
        }
        catch (DataIntegrityViolationException e) {
            throw new CustomException("Failed to delete token", HttpStatus.CONFLICT); // todo проверить на ошибки
        }
        catch(Exception e){
            throw new CustomException("Failed to delete token", HttpStatus.INTERNAL_SERVER_ERROR);
        }
    }

    public void deleteUser(String email) {
        _userRepository.deleteByEmail(email);
    }

    //region Helpers

    private User checkLoginPasswordLogInSystem(User user) {

        User foundUser = _userRepository.findUserByEmail(user.getEmail());

        if (foundUser == null) {
            throw new CustomException("User not found", HttpStatus.NOT_FOUND);
        }

        if (!foundUser.getPassword().equals(user.getPassword())) {
            throw new CustomException("Wrong password", HttpStatus.UNAUTHORIZED);
        }

        return foundUser;
    }

    private SaveUserResponseTDO CreateUserResponseTDO(TokenRecord token) {
        SaveUserResponseTDO responseTDO = new SaveUserResponseTDO();
        responseTDO.status = 200;
        responseTDO.token = token.getTokenValue();
        responseTDO.createdDate = token.getCreatedDate();
        responseTDO.expirationDate = token.getExpirationDate();

        return responseTDO;
    }

    private TokenRecord createToken(User user){
        String token = UUID.randomUUID().toString();

        TokenRecord tokenRecord = new TokenRecord();
        tokenRecord.setUser(user);
        tokenRecord.setTokenValue(token);
        tokenRecord.setCreatedDate(LocalDateTime.now());
        tokenRecord.setExpirationDate(LocalDateTime.now().plusMonths(1));

        return tokenRecord;
    }

    @Nullable
    private String checkPassword(String password){
        boolean isValid = true;
        String message = null;

        if (password.length() < 8) {
            isValid = false;
            message = "Password must be at least 8 characters long";
        }

        if (!password.matches("^(?=.*[!@#$%^&*(),.?\":{}|<>]).*$")) {
            isValid = false;
            message = "Password must contain special characters";
        }

        return message;
    }

    @Nullable
    private String checkEmail(String email){
        if (!email.matches("^[A-Za-z0-9+_.-]+@[A-Za-z0-9.-]+$")) {
            return "Invalid email format";
        }
        else return null;
    }

    private void CheckIncomingParametersToCreateUser(User user) {
        String messageCheckPassword = checkPassword(user.getPassword());
        if (messageCheckPassword != null) {
            throw new CustomException(messageCheckPassword, HttpStatus.BAD_REQUEST);
        }

        String messageCheckEmail = checkEmail(user.getEmail());
        if(messageCheckEmail != null){
            throw new CustomException(messageCheckEmail, HttpStatus.BAD_REQUEST);
        }

        User existingUser = _userRepository.findUserByEmail(user.getEmail());
        if(existingUser != null){
            throw new CustomException("User with this email already exists", HttpStatus.CONFLICT);
        }
    }

    //endregion

}
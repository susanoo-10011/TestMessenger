package ru.TestMessenger.UsersService.Service;

import jakarta.servlet.http.HttpServletRequest;
import ru.TestMessenger.UsersService.Exceptions.CustomException;
import ru.TestMessenger.UsersService.Model.*;
import ru.TestMessenger.UsersService.Repository.IFollowerRepository;
import ru.TestMessenger.UsersService.Repository.IUserProfileRepository;
import ru.TestMessenger.UsersService.Repository.IUserRepository;
import ru.TestMessenger.UsersService.Repository.IUserTokenRepository;
import ru.TestMessenger.UsersService.Service.EntityResponseTDO.SaveUserResponseTDO;
import lombok.AllArgsConstructor;
import org.springframework.context.annotation.Primary;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.lang.Nullable;
import org.springframework.stereotype.Service;
import org.springframework.dao.DataIntegrityViolationException;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;
import java.util.*;

@Service
@AllArgsConstructor
@Primary
@Transactional
public class UserService implements IUserService{

    private final IUserRepository _userRepository;
    private final IUserTokenRepository _tokenRepository;
    private final IUserProfileRepository _userProfileRepository;
    private  final IFollowerRepository _followerRepository;


    @Override
    public ResponseEntity<List<User>> findAllUsers() {
        List<User> userEntities = _userRepository.findAll();

        if (userEntities.isEmpty()) {
           generateError("Users not found", HttpStatus.NOT_FOUND);
        }
        return new ResponseEntity<>(userEntities, HttpStatus.OK);
    }

    @Override
    public ResponseEntity<?> saveUser(User user) {

        CheckIncomingParametersToCreateUser(user);

        try {
            User savedUser = _userRepository.save(user);

            _userProfileRepository.save(createUserProfile(savedUser));

            return new ResponseEntity<>(savedUser, HttpStatus.OK);
        } catch (DataIntegrityViolationException e) {
            return generateError("Invalid request", HttpStatus.BAD_REQUEST);
        } catch (Exception e) {
            return generateError("Failed to save user", HttpStatus.INTERNAL_SERVER_ERROR);
        }
    }

    @Override
    public ResponseEntity<User> findByEmail(String email) {
        User user = _userRepository.findUserByEmail(email);
        if (user == null) {
            generateError("User not found", HttpStatus.NOT_FOUND);
        }
        return new ResponseEntity<>(user, HttpStatus.OK);
    }

    @Override
    public ResponseEntity<List<User>> findAllByName(String name){

        List<User> foundUsers = _userRepository.findAllByName(name);

        if (foundUsers == null) {
            generateError("User not found", HttpStatus.NOT_FOUND);
        }

        return new ResponseEntity<>(foundUsers, HttpStatus.OK);
    }

    @Override
    public ResponseEntity<User> findUserByLogin(String login){

        User foundUser = _userRepository.findUserByLogin(login);

        if (foundUser == null) {
            generateError("User not found", HttpStatus.NOT_FOUND);
        }

        return new ResponseEntity<>(foundUser, HttpStatus.OK);
    }

    @Override
    public ResponseEntity<List<String>> getFollowers(String token){

        TokenRecord foundToken = _tokenRepository.findByTokenValue(token);
        if (foundToken.getTokenValue() == null) {
            generateError("User not found", HttpStatus.NOT_FOUND);
        }

        User user = foundToken.getUser();

        List<Follower> followers = _followerRepository.getFollowersByUser(user);

        List<String> followerNames = new ArrayList<>();
        for (Follower follower : followers) {
            followerNames.add(follower.getFollowingUser().getLogin());
        }

        return new ResponseEntity<>(followerNames, HttpStatus.OK);
    }

    @Override
    public ResponseEntity<User> updateUser(User user, HttpServletRequest request) {
        try{

            String tokenValue = request.getHeaderNames().nextElement();
            TokenRecord foundToken = _tokenRepository.findByTokenValue(tokenValue);

            if(foundToken == null){
                generateError("The token is invalid", HttpStatus.BAD_REQUEST);
            }

            Optional<User> optionalUser = _userRepository.findById(user.getId());

            if (optionalUser.isEmpty()) {
                generateError("User not found", HttpStatus.NOT_FOUND);
            }

            _userRepository.save(user);
            return new ResponseEntity<>(user, HttpStatus.OK);
        }
        catch(Exception e){
            generateError("Failed to update user", HttpStatus.INTERNAL_SERVER_ERROR);
        }
        return null;
    }

    @Override
    public ResponseEntity<?> logInToSystem(User user) {
        User foundUser = checkLoginPasswordLogInSystem(user);

        TokenRecord token = createToken(user);

        try {
            _tokenRepository.save(token);

            return new ResponseEntity<>(token, HttpStatus.OK);
        } catch (DataIntegrityViolationException e) {
            return generateError("Invalid request", HttpStatus.BAD_REQUEST);
        } catch (Exception e) {
            return generateError("Failed to save user", HttpStatus.INTERNAL_SERVER_ERROR);
        }
    }

    @Override
    public ResponseEntity<?> logOutFromSystem(String token) {

        TokenRecord foundToken = _tokenRepository.findByTokenValue(token);
        if(foundToken == null){
            generateError("User not found", HttpStatus.NOT_FOUND);
        }

        try{
            _tokenRepository.deleteByTokenValue(token);
            return new ResponseEntity<>(HttpStatus.OK);
        }
        catch (DataIntegrityViolationException e) {
           return generateError("Failed to delete token", HttpStatus.CONFLICT);
        }
        catch(Exception e){
            return generateError("Failed to delete token", HttpStatus.INTERNAL_SERVER_ERROR);
        }
    }

    @Override
    public ResponseEntity<?> subscribeToUser(String token, String username) {

        TokenRecord foundToken = _tokenRepository.findByTokenValue(token);
        if (foundToken.getTokenValue() == null) {
            generateError("User not found", HttpStatus.NOT_FOUND);
        }

        try {
            User user = foundToken.getUser();
            User userSubscription = _userRepository.findUserByLogin(username);

            Follower follower = new Follower();
            follower.setUser(user);
            follower.setFollowingUser(userSubscription);

            _followerRepository.save(follower);
        } catch (Exception e){
            return generateError("Failed to subscribe to user", HttpStatus.INTERNAL_SERVER_ERROR);
        }

        return new ResponseEntity<>(foundToken, HttpStatus.OK);
    }

    public void deleteUser(String email) {
        _userRepository.deleteByEmail(email);
    }

    //region Helpers

    public Optional<User> findUserByTokenValue(String tokenValue) {
        TokenRecord tokenRecord = _tokenRepository.findByTokenValue(tokenValue);
        return tokenRecord != null ? Optional.of(tokenRecord.getUser()) : Optional.empty();
    }

    private ResponseEntity<?> generateError(String message, HttpStatus status){

        ErrorTDO errorTDO = new ErrorTDO();
        errorTDO.status = status.value();
        errorTDO.message = message;
        errorTDO.error = status;

        throw new CustomException(errorTDO);
    }


    private User checkLoginPasswordLogInSystem(User user) {

        User foundUser = _userRepository.findUserByLogin(user.getLogin());

        if (foundUser == null) {
            generateError("User not found", HttpStatus.NOT_FOUND);
        }

        if (!foundUser.getPassword().equals(user.getPassword())) {
            generateError("Wrong password", HttpStatus.UNAUTHORIZED);
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

        if(user.getId() == null){
            user = _userRepository.findUserByLogin(user.getLogin());
        }

        String token = UUID.randomUUID().toString();

        TokenRecord tokenRecord = new TokenRecord();
        tokenRecord.setUser(user);
        tokenRecord.setTokenValue(token);
        tokenRecord.setCreatedDate(LocalDateTime.now());
        tokenRecord.setExpirationDate(LocalDateTime.now().plusMonths(1));

        return tokenRecord;
    }

    private UserProfile createUserProfile(User user){

        UserProfile userProfile = new UserProfile();

        userProfile.setUser(user);

        return userProfile;
    }

    @Nullable
    private String checkPassword(String password){
        String message = null;

        if (password.length() < 8) {
            message = "Password must be at least 8 characters long";
        }

        if (!password.matches("^(?=.*[!@#$%^&*(),.?\":{}|<>]).*$")) {
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
            generateError(messageCheckPassword, HttpStatus.BAD_REQUEST);
        }

        String messageCheckEmail = checkEmail(user.getEmail());
        if(messageCheckEmail != null){
            generateError(messageCheckEmail, HttpStatus.BAD_REQUEST);
        }

        User existingUser = _userRepository.findUserByEmail(user.getEmail());
        if(existingUser != null){
            generateError("User with this email already exists", HttpStatus.CONFLICT);
        }
    }

    //endregion

}
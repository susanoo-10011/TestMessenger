package ru.TestMessenger.UsersService.Repository;

import ru.TestMessenger.UsersService.Model.User;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;

public interface IUserRepository extends JpaRepository<User, Long> {

    void deleteByEmail(String email);
    User findUserByEmail(String email);

    List<User> findAllByName(String name);

    User findUserByLogin(String login);

}


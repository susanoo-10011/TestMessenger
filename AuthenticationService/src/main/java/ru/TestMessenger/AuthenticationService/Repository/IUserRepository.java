package ru.TestMessenger.AuthenticationService.Repository;

import ru.TestMessenger.AuthenticationService.Model.User;
import org.springframework.data.jpa.repository.JpaRepository;

public interface IUserRepository extends JpaRepository<User, Long> {

    void deleteByEmail(String email);
    User findUserByEmail(String email);
}


package ru.TestMessenger.UsersService.Repository;

import ru.TestMessenger.UsersService.Model.TokenRecord;
import org.springframework.data.jpa.repository.JpaRepository;
import ru.TestMessenger.UsersService.Model.User;

import java.time.LocalDateTime;
import java.util.List;

public interface IUserTokenRepository extends JpaRepository<TokenRecord, Long> {

    void deleteByTokenValue(String token);
    List<TokenRecord> findTokensByExpirationDateLessThan(LocalDateTime expirationDate);

    TokenRecord findByTokenValue(String token);

    TokenRecord findByUser(User user);

    User findUserByTokenValue(String tokenValue);
}

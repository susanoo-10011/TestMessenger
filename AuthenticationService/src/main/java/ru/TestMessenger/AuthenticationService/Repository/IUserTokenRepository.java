package ru.TestMessenger.AuthenticationService.Repository;

import ru.TestMessenger.AuthenticationService.Model.TokenRecord;
import org.springframework.data.jpa.repository.JpaRepository;

import java.time.LocalDateTime;
import java.util.List;

public interface IUserTokenRepository extends JpaRepository<TokenRecord, Long> {

    void deleteByTokenValue(String token);
    List<TokenRecord> findTokensByExpirationDateLessThan(LocalDateTime expirationDate);
}

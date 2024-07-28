package ru.TestMessenger.UsersService.TaskService;

import ru.TestMessenger.UsersService.Model.TokenRecord;
import ru.TestMessenger.UsersService.Repository.IUserTokenRepository;
import lombok.AllArgsConstructor;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.util.List;

@Service
@AllArgsConstructor
public class TokenExpirationTaskService {

    private final IUserTokenRepository _tokenRepository;

    @Scheduled(cron = "30 15 * * * ?") // ызывается ежедневно в 15:30
    public void removeExpiredTokens() {
        List<TokenRecord> expiredTokens = _tokenRepository.findTokensByExpirationDateLessThan(LocalDateTime.now());

        expiredTokens.forEach(token -> {
            _tokenRepository.deleteAll();
        });
    }
}
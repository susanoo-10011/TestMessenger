package ru.TestMessenger.AuthenticationService;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.scheduling.annotation.EnableScheduling;

@SpringBootApplication
@EnableScheduling // активируем планировщик задач
public class TestMessengerApplication {

    public static void main(String[] args) {
        SpringApplication.run(TestMessengerApplication.class, args);
    }
}
package ru.TestMessenger.AuthenticationService.Service.EntityResponseTDO;

import org.springframework.http.HttpStatus;

import java.time.LocalDateTime;
import java.util.Date;


public class SaveUserResponseTDO {

    public int status;
    public String token;
    public LocalDateTime expirationDate;
    public LocalDateTime createdDate;
}

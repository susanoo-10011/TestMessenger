package ru.TestMessenger.UsersService.Service.EntityResponseTDO;

import java.time.LocalDateTime;


public class SaveUserResponseTDO {

    public int status;
    public String token;
    public LocalDateTime expirationDate;
    public LocalDateTime createdDate;
}

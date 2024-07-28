package ru.TestMessenger.UsersService.Exceptions;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Getter;
import org.springframework.http.HttpStatus;
import ru.TestMessenger.UsersService.Model.ErrorTDO;

@Getter
public class CustomException extends RuntimeException {

    @JsonProperty("status")
    private int statusCode;

    @JsonProperty("error")
    private ErrorTDO error;

    public CustomException(ErrorTDO errorTDO) {

        statusCode = errorTDO.status;
        error = errorTDO;
    }
}
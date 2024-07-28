package ru.TestMessenger.UsersService.Exceptions;

import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.ControllerAdvice;
import org.springframework.web.bind.annotation.ExceptionHandler;
import ru.TestMessenger.UsersService.Model.ErrorTDO;

@ControllerAdvice // перехватывает все исключения
public class GlobalExceptionHandler {

    @ExceptionHandler(CustomException.class)
    public ResponseEntity<ErrorTDO> handleCustomException(CustomException ex) {
        return ResponseEntity.status(ex.getStatusCode()).body(ex.getError());
    }
}

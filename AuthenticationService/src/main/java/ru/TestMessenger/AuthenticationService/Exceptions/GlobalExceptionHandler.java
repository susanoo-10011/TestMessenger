package ru.TestMessenger.AuthenticationService.Exceptions;

import ru.TestMessenger.AuthenticationService.Exceptions.ErrorEntity.ErrorResponseDTO;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.ControllerAdvice;
import org.springframework.web.bind.annotation.ExceptionHandler;

@ControllerAdvice // перехватывает все исключения
public class GlobalExceptionHandler {

    @ExceptionHandler(CustomException.class)
    public ResponseEntity<ErrorResponseDTO> handleCustomException(CustomException ex) {
        ErrorResponseDTO errorResponse = new ErrorResponseDTO(ex.getHttpStatus().value(), ex.getMessage());
        return ResponseEntity.status(ex.getHttpStatus().value()).body(errorResponse);
    }
}

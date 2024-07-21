package ru.TestMessenger.AuthenticationService.Exceptions.ErrorEntity;

import com.fasterxml.jackson.annotation.JsonProperty;

public class ErrorResponseDTO {
    @JsonProperty("status")
    private int statusCode;

    @JsonProperty("message")
    private String message;

    public ErrorResponseDTO(int statusCode, String message) {
        this.statusCode = statusCode;
        this.message = message;
    }

    // Getter и Setter
    public int getStatusCode() {
        return statusCode;
    }

    public void setStatusCode(int statusCode) {
        this.statusCode = statusCode;
    }

    public String getMessage() {
        return message;
    }

    public void setMessage(String message) {
        this.message = message;
    }
}

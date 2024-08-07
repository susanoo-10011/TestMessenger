package ru.TestMessenger.UsersService.Model;

import com.fasterxml.jackson.annotation.JsonIgnore;
import jakarta.persistence.*;
import lombok.Data;

import java.time.LocalDateTime;

@Data
@Entity
@Table(name = "token_records", schema = "auth_schema")
public class TokenRecord {

    @Id
    @JsonIgnore
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @OneToOne(fetch = FetchType.EAGER)
    @JoinColumn(name = "user_id", nullable = false)
    @JsonIgnore
    private User user;

    @Column(unique = true)
    private String tokenValue;

    @Column(name = "expires_at")
    private LocalDateTime expirationDate;

    @Column(name = "created_at")
    private LocalDateTime createdDate;
}
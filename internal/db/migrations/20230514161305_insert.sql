-- +goose Up
-- +goose StatementBegin
INSERT INTO users(email, password, role_name) VALUES ('jefferson@gmail.com','$2a$10$9EgGuxvFDbTfckfdc/yPsOFxfTXo3V5fmdVX73N6OBRyUbQ//2Yu6','admin');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM users WHERE id = 1;
-- +goose StatementEnd

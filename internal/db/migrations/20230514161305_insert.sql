-- +goose Up
-- +goose StatementBegin
INSERT INTO users(email, password, role_name) VALUES ('jefferson@gmail.com','$2a$10$9EgGuxvFDbTfckfdc/yPsOFxfTXo3V5fmdVX73N6OBRyUbQ//2Yu6','admin');
INSERT INTO categories VALUES ('Шутер');
INSERT INTO categories VALUES ('MMORPG');
INSERT INTO categories VALUES ('RPG');
INSERT INTO categories VALUES ('Стратегия');
INSERT INTO categories VALUES ('Симулятор');
INSERT INTO categories VALUES ('Survival');
INSERT INTO companies VALUES ('Valve');
INSERT INTO companies VALUES ('Epic Games');
INSERT INTO companies VALUES ('Bethesda Games');
INSERT INTO companies VALUES ('Blizzard');
INSERT INTO companies VALUES ('Sony');
INSERT INTO companies VALUES ('Electronic Arts');
INSERT INTO companies VALUES ('Ubisoft');
INSERT INTO products (name, description, price, stock, image, category_id, company_id) VALUES ('Fallout 3','Увлекательная RPG в мире постапокалипсиса', 550, 125, 'fallout3.jpeg', 3, 3);
INSERT INTO products (name, description, price, stock, image, category_id, company_id) VALUES ('Assassins Creed: Valhalla','Окунитесь в мир викингов в новой RPG от Ubisoft', 1499, 35, 'valhalla.jpeg', 3, 7);
INSERT INTO products (name, description, price, stock, image, category_id, company_id)
VALUES ('Counter-Strike: 2 Prime', 'Прайм статус для игры, покорившей миллионы', 850, 100, 'cs2.jpg', 1, 1);

INSERT INTO products (name, description, price, stock, image, category_id, company_id)
VALUES ('World of Warcraft: Dragonflight', 'Новый аддон для самой популярной MMORPG в мире', 4299, 25, 'dragonflight.jpeg', 2, 4);

INSERT INTO products (name, description, price, stock, image, category_id, company_id)
VALUES ('The Elder Scrolls V: Skyrim', 'Эпичная role-play фэнтези', 399, 30, 'skyrim.jpeg', 3, 3);

INSERT INTO products (name, description, price, stock, image, category_id, company_id)
VALUES ('Civilization VI', 'Одна из лучших пошаговых стратегий последних лет', 1999, 20, 'civ6.jpeg', 4, 2);

INSERT INTO products (name, description, price, stock, image, category_id, company_id)
VALUES ('FIFA 23', 'Отличный симулятор футбола', 2459, 10, 'fifa23.jpeg', 5, 6);

INSERT INTO products (name, description, price, stock, image, category_id, company_id)
VALUES ('Overwatch', 'Командный шутер, полюбившийся многим', 999, 50, 'overwatch.jpeg', 1, 4);

INSERT INTO products (name, description, price, stock, image, category_id, company_id)
VALUES ('God of War', 'Увлекательная RPG, повествующая о жизни бога войны', 3499, 20, 'god_of_war.jpeg', 3, 5);

INSERT INTO products (name, description, price, stock, image, category_id, company_id)
VALUES ('The Sims 4', 'Симулятор жизни в вашем компьютере', 799, 40, 'sims4.jpeg', 5, 6);

INSERT INTO products (name, description, price, stock, image, category_id, company_id)
VALUES ('1000 В-баксов Fortnite', 'Валюта для самого популярного battle-royale шутера', 1100, 100, 'vbucksfortnite.jpeg', 1, 2);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM users WHERE id = 1;
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "pages" (
    "id" INTEGER PRIMARY KEY AUTOINCREMENT,
    "url" TEXT NOT NULL,
    "user_name" TEXT NOT NULL,
    "status" INTEGER NOT NULL,
    "created_at" DATETIME NOT NULL,
    "updated_at" DATETIME NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "pages";
-- +goose StatementEnd

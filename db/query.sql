-- name: GetPetByID :one
SELECT * FROM pets
WHERE id = ? LIMIT 1;

-- name: ListPets :many
SELECT * FROM pets
ORDER BY id;

-- name: CreatePet :execresult
INSERT INTO pets (
  name, status
) VALUES (
  ?, ?
);

-- name: DeletePet :exec
DELETE FROM pets
WHERE id = ?;

-- name: UpdatePet :exec
UPDATE pets
set name = ?, status = ?
WHERE id = ?;

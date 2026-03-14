-- name: insert
INSERT INTO templates ("name", body, is_default, subject, type, team_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING
    id,
    created_at,
    updated_at,
    type,
    body,
    is_default,
    name,
    subject,
    is_builtin,
    team_id,
    NULL::TEXT AS team_name;

-- name: update
WITH u AS (
    UPDATE templates
    SET 
        name = $2,
        body = $3,
        is_default = $4,
        subject = $5,
        type = $6::template_type,
        team_id = $7,
        updated_at = NOW()
    WHERE id = $1
    RETURNING *
)
SELECT
    u.id,
    u.created_at,
    u.updated_at,
    u.type,
    u.body,
    u.is_default,
    u.name,
    u.subject,
    u.is_builtin,
    u.team_id,
    teams.name AS team_name
FROM u
LEFT JOIN teams ON teams.id = u.team_id
LIMIT 1;

-- name: get-default
SELECT
    id,
    created_at,
    updated_at,
    type,
    body,
    is_default,
    name,
    subject,
    is_builtin,
    team_id,
    NULL::TEXT AS team_name
FROM templates
WHERE is_default IS TRUE
  AND type = 'email_outgoing'
  AND team_id IS NULL;

-- name: get-all
SELECT
    templates.id,
    templates.created_at,
    templates.updated_at,
    templates.type,
    templates.body,
    templates.is_default,
    templates.name,
    templates.subject,
    templates.is_builtin,
    templates.team_id,
    teams.name AS team_name
FROM templates
LEFT JOIN teams ON teams.id = templates.team_id
WHERE templates.type = $1
ORDER BY templates.is_default DESC, templates.updated_at DESC;

-- name: get-all-by-team
SELECT
    templates.id,
    templates.created_at,
    templates.updated_at,
    templates.type,
    templates.body,
    templates.is_default,
    templates.name,
    templates.subject,
    templates.is_builtin,
    templates.team_id,
    teams.name AS team_name
FROM templates
LEFT JOIN teams ON teams.id = templates.team_id
WHERE templates.type = $1
  AND (
    templates.team_id = $2
    OR ($3::BOOLEAN IS TRUE AND templates.team_id IS NULL)
  )
ORDER BY
    CASE
        WHEN templates.team_id = $2 THEN 0
        WHEN templates.team_id IS NULL THEN 1
        ELSE 2
    END,
    templates.is_default DESC,
    templates.updated_at DESC;

-- name: get-template
SELECT
    templates.id,
    templates.created_at,
    templates.updated_at,
    templates.type,
    templates.body,
    templates.is_default,
    templates.name,
    templates.subject,
    templates.is_builtin,
    templates.team_id,
    teams.name AS team_name
FROM templates
LEFT JOIN teams ON teams.id = templates.team_id
WHERE templates.id = $1;

-- name: delete
DELETE FROM templates WHERE id = $1;

-- name: get-by-name
SELECT
    id,
    created_at,
    updated_at,
    type,
    body,
    is_default,
    name,
    subject,
    is_builtin,
    team_id,
    NULL::TEXT AS team_name
FROM templates
WHERE name = $1;

-- name: is-builtin
SELECT EXISTS(SELECT 1 FROM templates WHERE id = $1 AND is_builtin is TRUE);
